package status

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

// Group represents a concurrent set of operations.
// Each operation takes an Item as a parameter, as well as an [io.Writer].
// Each operation returns a Result.
// Each writer writes to a dedicated line of a [Status].
type Group[Item any, Result any] struct {
	// PrefixString is called once on each line of the [Status] to add a prefix.
	// When nil, [DefaultPrefixString] is used.
	PrefixString func(item Item, index int) string

	// When PrefixAlign is set, automatically ensure that all prefixes are of the same length,
	// by adding appropriate spaces.
	PrefixAlign bool

	// ResultString is called to generate a message for when the given item has finished processing.
	// It is called with the returned error.
	// When nil, [DefaultErrString] is used.
	ResultString func(res Result, item Item, index int) string

	// Handler is a handler called for each item to run.
	// It is passed an io.Writer that writes directly to the specified line of the status.
	// Handler must not be nil.
	Handler func(item Item, index int, writer io.Writer) Result

	// HandlerLimit is the maximum number of handlers to run concurrently.
	// A HandlerLimit of <= 0 indicates no limit.
	//
	// Handlers are in principle called in order, however for HandlerLimit > 1
	// this cannot be strictly enforced.
	//
	// The Limit is only enforced within a single call to [Use] or [].
	HandlerLimit int

	// WaitString is called when the status line for a specific handler is initialized, but the Handler has not yet been called.
	//
	// When WaitString is nil, lines are only initialized once they are needed.
	// Setting WaitString != nil causes output to appear in order.
	WaitString func(item Item, index int) string
}

// DefaultPrefixString is the default implementation of [Group.PrefixString].
// It uses the default 'v' verb of the 'fmt' package to format the item.
func DefaultPrefixString[Item any](item Item, index int) string {
	return fmt.Sprintf("%v: ", item)
}

// DefaultResultString is the default implementation of [Group.ResultString].
// It uses fmt.Sprint on the result type.
func DefaultResultString[Item, Result any](result Result, item Item, index int) string {
	return fmt.Sprint(result)
}

// DefaultWaitString returns the string "waiting" for any item.
func DefaultWaitString[Item any](item Item, index int) string {
	return "waiting"
}

// Use calls Handler for all passed items.
//
// It sends output to the provided status, while respecting HandlerLimit.
// Each output is displayed on a separate line.
//
// If group.WaitString is nil, lines are closed as soon as they are no longer needed.
// Otherwise they are closed right before returning.
//
// Use returns once all Handlers have returned.
func (group Group[Item, Result]) Use(status *Status, items []Item) []Result {
	// setup defaults
	if group.PrefixString == nil {
		group.PrefixString = DefaultPrefixString[Item]
	}
	if group.ResultString == nil {
		group.ResultString = DefaultResultString[Item, Result]
	}

	// create data arrays
	prefixes := make([]string, len(items))        // prefixes per-line
	writers := make([]io.WriteCloser, len(items)) // writers per-line
	results := make([]Result, len(items))         // results per item

	// generate all the prefixes and compute the maximum prefix length
	var maxPrefixLength int
	if group.PrefixString != nil {
		for index, item := range items {
			prefixes[index] = group.PrefixString(item, index)

			if len(prefixes[index]) > maxPrefixLength {
				maxPrefixLength = len(prefixes[index])
			}
		}
	}

	// if requested, align the prefixes
	if group.PrefixAlign {
		for index, prefix := range prefixes {
			prefixes[index] += strings.Repeat(" ", maxPrefixLength-len(prefix))
		}
	}

	// if we have a limit, create a channel for tokens
	hasLimit := group.HandlerLimit > 0

	var tokens chan struct{}
	if hasLimit {
		tokens = make(chan struct{}, group.HandlerLimit)
	}

	// initialize all the lines (if needed)
	if group.WaitString != nil {
		for index, item := range items {
			writers[index] = status.OpenLine(prefixes[index], group.WaitString(item, index))
		}
	}

	// prepare a waitgroup for all the tasks.
	var wg sync.WaitGroup
	wg.Add(len(items))

	// we want to run tasks as much in order as is feasible.
	// so we spawn as many workers as possible, and send them tasks in order.
	indexes := make(chan int)

	// start all the workers first
	for range items {
		go func() {
			defer wg.Done()

			// if we have a limit, wait for it!
			if hasLimit {
				tokens <- struct{}{}
				defer func() {
					<-tokens
				}()
			}

			// grab the next index
			index := <-indexes
			item := items[index]

			// if the line hasn't yet been setup, create it!
			if group.WaitString == nil {
				writers[index] = status.OpenLine(prefixes[index], "")
				defer writers[index].Close()
			}

			// write into the result array
			results[index] = group.Handler(item, index, writers[index])

			// and write out the result
			io.WriteString(writers[index], "\n"+group.ResultString(results[index], item, index)+"\n")
		}()
	}

	// run tasks for all the indexes
	for index := range items {
		indexes <- index
	}

	// and wait for them to complete
	wg.Wait()

	// if we didn't initialize the waiters beforehand
	// then we still need to close them all!
	if group.WaitString != nil {
		for _, w := range writers {
			w.Close()
		}
	}

	return results
}

// Run creates a new Status, and then directs output to it using [Use].
//
// See also [New], [Use].
func (group Group[Item, Result]) Run(writer io.Writer, items []Item) []Result {
	// setup the status!
	status := New(writer, 0)
	status.Start()
	defer status.Stop()

	// and use it!
	return group.Use(status, items)
}

// DefaultErrorString implements the default result handler for [UseErrorGroup] and [RunErrorGroup].
// When error is nil, returns the string "done", else returns the string "failed" with an error description.
func DefaultErrorString[Item any](err error, item Item, index int) string {
	if err == nil {
		return "done"
	}
	return fmt.Sprintf("failed (%v)", err)
}

// UseErrorGroup calls group.Use(status, items) and returns the first non-nil error, or nil otherwise.
// When group.ResultString is nil, uses [DefaultErrorString] instead.
func UseErrorGroup[Item any](status *Status, group Group[Item, error], items []Item) error {
	if group.ResultString == nil {
		group.ResultString = DefaultErrorString[Item]
	}

	errors := group.Use(status, items)
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

// RunErrorGroup calls group.Run(status, items) and returns the first non-nil error, or nil otherwise.
// When group.ResultString is nil, uses [DefaultErrorString] instead.
func RunErrorGroup[Item any](writer io.Writer, group Group[Item, error], items []Item) error {
	// setup the status!
	status := New(writer, 0)
	status.Start()
	defer status.Stop()

	// Use it!
	return UseErrorGroup(status, group, items)
}
