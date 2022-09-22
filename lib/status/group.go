package status

import (
	"io"
	"strings"
	"sync"
)

// Group represents a group of operations that each write to a separate line of a [Status].
type Group[Item any] struct {
	// Writer is the underyling writer of this Group.
	Writer io.Writer

	// PrefixString is called on each line of the [Status] to add a prefix.
	// When nil, it is assumed to return the empty string instead.
	PrefixString func(item Item, index int) string

	// When PrefixAlign is set, automatically ensure that all prefixes are of the same length,
	// by adding appropriate spaces.
	PrefixAlign bool

	// ErrString is called to generate a message for when the given item has finished processing.
	// It is called with the returned error, and MUST NOT be nil.
	ErrString func(item Item, index int, err error) string

	// Handler is a handler called for each item to run.
	// It is passed an io.Writer that writes directly to the specified line of the status.
	// Handler must not be nil.
	Handler func(item Item, index int, writer io.Writer) error
}

// Run calls Handler for all passed Items concurrently, each passing output to a dedicated line of a new [Status].
//
// It returns the first non-nil error returned from the Handler invocations.
// Run always waits for all handlers to return, regardless which one returns an error.
func (group Group[Item]) Run(items []Item) error {
	// setup the status!
	status := New(group.Writer, len(items))
	status.Start()
	defer status.Stop()

	// create a waitgroup
	var wg sync.WaitGroup
	wg.Add(len(items))

	// generate all the prefixes and compute the maximum prefix length
	prefixes := make([]string, len(items))
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

	// generate an errors array
	errors := make([]error, len(items))
	for index, item := range items {
		// and call the handler functions
		go func(index int, item Item, prefix string) {
			defer wg.Done()

			// and a writer
			writer := status.Line(prefix, index)

			// write into the error array
			errors[index] = group.Handler(item, index, writer)

			// and write out the result
			io.WriteString(writer, "\n"+group.ErrString(item, index, errors[index])+"\n")
		}(index, item, prefixes[index])
	}

	// wait for everything to finish
	wg.Wait()

	// return the first non-nil error
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
