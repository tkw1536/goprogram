package status

import (
	"io"

	"github.com/tkw1536/goprogram/stream"
)

// StreamGroup intelligently runs handler over items concurrently.
//
// Count determines the number of concurrent invocations to run.
// count <= 0 indicates no limit.
// count = 1 indicates running handler in order.
//
// handler is additionally passed an IOStream.
// When there is only one concurrent invocation, the original stream as a parameter.
// When there is more than one concurrent invocation, each invocation is passed a single line of a new [Status].
// The [Status] will send output to the standard output of str.
//
// StreamGroup returns the first non-nil error returned by each call to handler; or nil otherwise.
func StreamGroup[T any](str stream.IOStream, count int, handler func(value T, stream stream.IOStream) error, items []T, opts ...StreamGroupOption[T]) error {

	// create a group
	var group Group[T, error]
	group.HandlerLimit = count

	// apply all the options
	isParallel := count != 1
	for _, opt := range opts {
		group = opt(isParallel, group)
	}

	// setup the default prefix string
	if group.PrefixString == nil {
		group.PrefixString = DefaultPrefixString[T]
	}

	// then just iterate over the items
	if !isParallel {
		for index, item := range items {
			str.Println(group.PrefixString(item, index))
			err := handler(item, str)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// if we are running in parallel, setup a handler
	group.Handler = func(item T, index int, writer io.Writer) error {
		ios := stream.NewIOStream(writer, writer, nil, 0)
		return handler(item, ios)
	}

	// create a new status display
	st := NewWithCompat(str.Stdout, 0)
	st.Start()
	defer st.Stop()

	// and use it!
	return UseErrorGroup(st, group, items)
}

// StreamGroupOption represents an option for [StreamGroup].
// The boolean indicates if the option is being applied to a status line or not.
type StreamGroupOption[T any] func(bool, Group[T, error]) Group[T, error]

// SmartMessage sets the message to display as a prefix before invoking a handler.
func SmartMessage[T any](handler func(value T) string) StreamGroupOption[T] {
	return func(p bool, s Group[T, error]) Group[T, error] {
		s.PrefixString = func(item T, index int) string {
			message := handler(item)
			if p {
				return "[" + message + "]: "
			}
			return message
		}
		s.PrefixAlign = true
		return s
	}
}
