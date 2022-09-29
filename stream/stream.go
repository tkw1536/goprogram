// Packate stream defines input and output streams.
package stream

import (
	"fmt"
	"io"

	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/lib/docfmt"
	"github.com/tkw1536/goprogram/lib/wrap"
	"golang.org/x/term"
)

// IOStream represents a set of input and output streams commonly associated to a process.
type IOStream struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer

	// Number of columns to wrap input and output in
	wrap int
}

// StdinIsATerminal checks if standard input is a terminal
func (io IOStream) StdinIsATerminal() bool {
	return streamIsTerminal(io.Stdin)
}

func (io IOStream) StdoutIsATerminal() bool {
	return streamIsTerminal(io.Stdout)
}

func (io IOStream) StderrIsATerminal() bool {
	return streamIsTerminal(io.Stderr)
}

// streamIsTerminal checks if stream is a terminal
func streamIsTerminal(stream any) bool {
	file, ok := stream.(interface{ Fd() uintptr })
	return ok && term.IsTerminal(int(file.Fd()))
}

// Printf is like [fmt.Printf] but prints to io.Stdout.
func (io IOStream) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(io.Stdout, format, args...)
}

// EPrintf is like [fmt.Printf] but prints to io.Stderr.
func (io IOStream) EPrintf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(io.Stderr, format, args...)
}

// Print is like [fmt.Print] but prints to io.Stdout.
func (io IOStream) Print(args ...interface{}) (n int, err error) {
	return fmt.Fprint(io.Stdout, args...)
}

// EPrint is like [fmt.Print] but prints to io.Stderr.
func (io IOStream) EPrint(args ...interface{}) (n int, err error) {
	return fmt.Fprint(io.Stderr, args...)
}

// Println is like [fmt.Println] but prints to io.Stdout.
func (io IOStream) Println(args ...interface{}) (n int, err error) {
	return fmt.Fprintln(io.Stdout, args...)
}

// EPrintln is like [fmt.Println] but prints to io.Stderr.
func (io IOStream) EPrintln(args ...interface{}) (n int, err error) {
	return fmt.Fprintln(io.Stderr, args...)
}

// ioDefaultWrap is the default value for Wrap of an IOStream.
const ioDefaultWrap = 80

// NewIOStream creates a new IOStream with the provided readers and writers.
// If any of them are set to an empty stream, they are set to util.NullStream.
// When wrap is set to 0, it is set to a reasonable default.
//
// It furthermore wraps output as set by wrap.
func NewIOStream(Stdout, Stderr io.Writer, Stdin io.Reader, wrap int) IOStream {
	if Stdout == nil {
		Stdout = Null
	}
	if Stderr == nil {
		Stderr = Null
	}
	if Stdin == nil {
		Stdin = Null
	}
	if wrap == 0 {
		wrap = ioDefaultWrap
	}
	return IOStream{
		Stdin:  Stdin,
		Stdout: Stdout,
		Stderr: Stderr,
		wrap:   wrap,
	}
}

// Streams creates a new IOStream with the provided streams and wrap.
// If any parameter is the zero value, copies the values from io.
func (io IOStream) Streams(Stdout, Stderr io.Writer, Stdin io.Reader, wrap int) IOStream {
	if Stdout == nil {
		Stdout = io.Stdout
	}
	if Stderr == nil {
		Stderr = io.Stderr
	}
	if Stdin == nil {
		Stdin = io.Stdin
	}
	if wrap == 0 {
		wrap = io.wrap
	}
	return NewIOStream(Stdout, Stderr, Stdin, wrap)
}

// NonInteractive creates a new IOStream with [Null] as standard input.
func (io IOStream) NonInteractive() IOStream {
	return io.Streams(nil, nil, Null, 0)
}

var newLine = []byte("\n")

// StdoutWriteWrap is like
//
//	io.Stdout.Write([]byte(s + "\n"))
//
// but wrapped at a reasonable length
func (io IOStream) StdoutWriteWrap(s string) (int, error) {
	n, err := wrap.Write(io.Stdout, io.wrap, s)
	if err != nil {
		return n, err
	}
	m, err := io.Stdout.Write(newLine)
	n += m
	return n, err
}

// StderrWriteWrap is like
//
//	io.Stdout.Write([]byte(s + "\n"))
//
// but wrapped at length Wrap.
func (io IOStream) StderrWriteWrap(s string) (int, error) {
	n, err := wrap.Write(io.Stderr, io.wrap, s)
	if err != nil {
		return n, err
	}
	m, err := io.Stderr.Write(newLine)
	n += m
	return n, err
}

var errDieUnknown = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unknown Error",
}

// Die prints a non-nil err to io.Stderr and returns an error of type Error or nil.
//
// When error is non-nil, this function first turns err into type Error.
// Then if err.Message is not the empty string, it prints it to io.Stderr with wrapping.
//
// If err is nil, it does nothing and returns nil.
func (io IOStream) Die(err error) error {
	var e exit.Error
	switch ee := err.(type) {
	case nil:
		return nil
	case exit.Error:
		e = ee
	default:
		e = errDieUnknown.Wrap(ee)
	}

	// print the error message to standard error in a wrapped way
	if message := e.Error(); message != "" {
		io.StderrWriteWrap(docfmt.Format(message))
	}

	return e
}
