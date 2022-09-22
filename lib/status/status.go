// Package status provides status, line and group
package status

import (
	"fmt"
	"io"
	"time"

	"github.com/gosuri/uilive"
)

// Status represents a writer that can write to multiple lines at once.
// The zero value is not ready to use, new instances should be created using [New].
type Status struct {
	w *uilive.Writer // underlying uilive writer

	updates chan updatemsg // channel that status updates are sent to
	done    chan struct{}
	status  []string // current status messages

	lastFlush time.Time // last time the writer was flushed
}

// updatemsg is sent to update a specific status line
type updatemsg struct {
	index   int
	message string
}

// New creates a new writer with the provided number of lines.
// Count is the initial number of lines to support, it must be >= 0, or New panics.
func New(writer io.Writer, count int) *Status {
	if count <= 0 {
		panic("New: count <= 0")
	}

	st := &Status{
		w: uilive.New(),

		updates: make(chan updatemsg, count),
		done:    make(chan struct{}),

		status: make([]string, count),
	}
	st.w.Out = writer
	return st
}

// Start instructs this Status to start writing output to the underlying writer.
// Calling start multiple times is an error, and may cause undefined behavior.
//
// No other process should write to the underlying writer, while this process is running.
// Instead [Bypass] should be used.
// See also [Stop].
func (st *Status) Start() {
	go st.listen()
}

func (st *Status) listen() {
	defer close(st.done)
	for msg := range st.updates {
		// update the right index
		if msg.index < 0 || msg.index > len(st.status) {
			continue
		}
		st.status[msg.index] = msg.message

		// and flush
		st.flush(false)
	}
}

const minFlushDelay = 50 * time.Millisecond

// flush flushes the output of this writer, unless less than minFlushDelay has elapsed
func (st *Status) flush(force bool) {
	now := time.Now()
	if !force && now.Sub(st.lastFlush) < minFlushDelay {
		return
	}
	st.lastFlush = now

	// write out each of the lines
	var line io.Writer
	for i, msg := range st.status {
		if i == 0 {
			line = st.w
		} else {
			line = st.w.Newline()
		}

		fmt.Fprintln(line, msg)
	}

	// flush the output
	st.w.Flush()
}

// Stop stops listening for updates, and blocks until all updates have completed.
func (st *Status) Stop() {
	close(st.updates)
	<-st.done
	st.flush(true) // force a flush!
}

// Set sets the provided status line to the given message.
// [Start] must have been called.
//
// Message must not end with a newline character.
//
// Set may safely be called called concurrently.
func (st *Status) Set(message string, index int) {
	st.updates <- updatemsg{
		message: message,
		index:   index,
	}
}

// Line returns an io.Writer that writes to the specified line of this status, always prepending the line with the given output.
func (st *Status) Line(prefix string, index int) io.Writer {
	return &LineBuffer{
		Line: func(line string) { st.Set(prefix+line, index) },
	}
}

// Bypass returns a writer that completely bypasses this Status, and writes directly to the underlying writer.
// [Start] must have been called.
func (st *Status) Bypass() io.Writer {
	return st.w.Bypass()
}
