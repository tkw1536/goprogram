// Package status provides status, line and group
package status

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/gosuri/uilive"
)

// Status represents an interactive status display that can write to multiple lines at once.
//
// A Status must be initialized using [New], then started (and stopped again) to write messages.
// Status may not be reused.
//
// A typical usage is as follows:
//
//   st := New(os.Stdout, 10)
//   st.Start()
//   defer st.Stop()
//
//   // ... whatever usage here ...
//   st.Set("line 0", 0)
//
//
// Using the status to Write messages outside of the Start / Stop process results in no-ops.
//
//
//
type Status struct {
	state uint64 // see state* message

	w *uilive.Writer // underlying uilive writer

	counter int32 // the first free message id, increased atomically

	ids      []int          // list of active message ids
	messages map[int]string // content of all the messages

	lastFlush time.Time // last time we flushed

	actions chan action // channel that status updates are sent to
	done    chan struct{}
}

// state of the state
const (
	stateInvalid uint64 = iota
	stateInit
	stateAlive
	stateDone
)

type actionTp uint8

const (
	writeLineAction actionTp = iota
	closeLineAction
	addLineAction
)

// action is sent to update a specific status line
type action struct {
	action actionTp
	index  int    // index is the index of the object to operate on
	data   string // data is the string to update
}

// New creates a new writer with the provided number of status lines.
//
// The ids of the status lines are guaranteed to be 0...(count-1).
// Count is assumed to be at least 0.
func New(writer io.Writer, count int) *Status {
	if count < 0 {
		count = 0
	}

	st := &Status{
		state: stateInit,

		w: uilive.New(),

		counter:  int32(count),
		ids:      make([]int, count),
		messages: make(map[int]string, count),

		actions: make(chan action, count),
		done:    make(chan struct{}),
	}
	for i := range st.ids {
		st.ids[i] = i
	}

	st.w.Out = writer
	return st
}

// Start instructs this Status to start writing output to the underlying writer.
//
// No other process should write to the underlying writer, while this process is running.
// Instead [Bypass] should be used.
// See also [Stop].
//
// Start may not be called more than once, extra calls may result in a panic.
func (st *Status) Start() {
	if atomic.LoadUint64(&st.state) == stateInvalid {
		panic("Status: Not created using New")
	}
	if !atomic.CompareAndSwapUint64(&st.state, stateInit, stateAlive) {
		panic("Status: Start() called multiple times")
	}

	go st.listen()
}

const minFlushDelay = 50 * time.Millisecond

// flush flushes the output of this Status to the underlying writer.
// flush respects [minFlushDelay], unless force is set to true.
func (st *Status) flush(force bool) {
	now := time.Now()
	if !force && now.Sub(st.lastFlush) < minFlushDelay {
		return
	}
	st.lastFlush = now

	// write out each of the lines
	var line io.Writer
	for i, key := range st.ids {
		if i == 0 {
			line = st.w
		} else {
			line = st.w.Newline()
		}

		fmt.Fprintln(line, st.messages[key])
	}

	// flush the output
	st.w.Flush()
}

// Stop blocks until all updates to finish processing.
// It then stops writing updates to the underlying writer.
//
// Stop must be called after [Start] has been called.
// Start may not be called more than once.
func (st *Status) Stop() {
	if !atomic.CompareAndSwapUint64(&st.state, stateAlive, stateDone) {
		panic("Status: Stop() called out-of-order")
	}

	close(st.actions)
	<-st.done
	st.flush(true) // force a flush!
}

// Set sets the status line with the given id to contain message.
// message should not contain newline characters.
// Set may block until the addition has been processed.
//
// Calling Set on a line which is not active results is a no-op.
//
// Set may safely be called concurrently with other methods.
//
// Set may only be called after [Start] has been called, but before [Stop].
// Other calls are silently ignored, and return an invalid line id.
func (st *Status) Set(message string, id int) {
	if atomic.LoadUint64(&st.state) != stateAlive {
		return
	}

	st.actions <- action{
		action: writeLineAction,
		index:  id,
		data:   message,
	}
}

// Line returns an [io.WriteCloser] linked to the status line with the provided id.
// Writing a complete newline-delimited line to it behaves just like [Set] with that line prefixed with prefix would.
// Calling [io.WriteCloser.Close] behaves just like [Done] would.
//
// Line may be called at any time.
func (st *Status) Line(prefix string, index int) io.WriteCloser {
	return &LineBuffer{
		Line:      func(line string) { st.Set(prefix+line, index) },
		CloseLine: func() { st.Done(index) },
	}
}

// Add adds a new status line and returns it's id.
// It may be further updated with calls to [Set], or removed with [Done].
// Add may block until the addition has been processed.
//
// Add may safely be called concurrently with other methods.
//
// Add may only be called after [Start] has been called, but before [Stop].
// Other calls are silently ignored, and return an invalid line id.
func (st *Status) Add(content string) (id int) {

	// even when not active, generate a new id
	// this guarantees that other calls are no-ops.
	id = int(atomic.AddInt32(&st.counter, 1))
	if atomic.LoadUint64(&st.state) != stateAlive {
		return
	}

	st.actions <- action{
		action: addLineAction,
		index:  id,
		data:   content,
	}
	return
}

// AddLine behaves like a call to [Add] followed by a call to [Line].
//
// AddLine may only be called after [Start] has been called, but before [Stop].
// Other calls are silently ignored, and return a no-op io.Writer.
func (st *Status) AddLine(prefix string) io.WriteCloser {
	return st.Line(prefix, st.Add(prefix))
}

// Done removes the status line with the provided id from this status.
// The last value of the status line is written to the top of the output.
// Done may block until the removal has been processed.
//
// Calling Done on a line which is not active results is a no-op.
//
// Done may safely be called concurrently with other methods.
//
// Done may only be called after [Start] has been called, but before [Stop].
// Other calls are silently ignored.
func (st *Status) Done(id int) {
	if atomic.LoadUint64(&st.state) != stateAlive {
		return
	}

	st.actions <- action{
		action: closeLineAction,
		index:  id,
	}
}

// listen listens for updates
func (st *Status) listen() {
	defer close(st.done)
	for msg := range st.actions {
		switch msg.action {
		case writeLineAction:
			for _, activeID := range st.ids {
				if activeID == msg.index {
					// update the message content
					st.messages[activeID] = msg.data
					break
				}
			}
			st.flush(false)
		case closeLineAction:
			for i, activeID := range st.ids {
				if activeID == msg.index {
					// write out the line to the bypass!
					fmt.Fprintln(st.w.Bypass(), st.messages[activeID])
					delete(st.messages, activeID)
					st.ids = append(st.ids[:i], st.ids[i+1:]...)
					break
				}
			}
			st.flush(true) // force a flush!
		case addLineAction:
			st.ids = append(st.ids, msg.index)
			st.messages[msg.index] = msg.data
			st.flush(true) // force a flush!
		}
	}
}

// Bypass returns a writer that completely bypasses this Status, and writes directly to the underlying writer.
// [Start] must have been called.
func (st *Status) Bypass() io.Writer {
	return st.w.Bypass()
}
