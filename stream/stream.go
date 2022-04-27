// Packate stream defines input and output streams.
package stream

import (
	"fmt"
	"io"

	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/lib/docfmt"
	"github.com/tkw1536/goprogram/lib/wrap"
)

// IOStream represents a set of input and output streams commonly associated to a process.
type IOStream struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer

	// Number of columns to wrap input and output in
	wrap int
}

// Printf is like "fmt.Printf" but prints to io.Stdout.
func (io IOStream) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(io.Stdout, format, args...)
}

// EPrintf is like "fmt.EPrintf" but prints to io.Stderr.
func (io IOStream) EPrintf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(io.Stderr, format, args...)
}

// Println is like "fmt.Println" but prints to io.Stdout.
func (io IOStream) Println(args ...interface{}) (n int, err error) {
	return fmt.Fprintln(io.Stdout, args...)
}

// EPrintln is like "fmt.Println" but prints to io.Stderr.
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

var newLine = []byte("\n")

// StdoutWriteWrap is like
//  io.Stdout.Write([]byte(s + "\n"))
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
//  io.Stdout.Write([]byte(s + "\n"))
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
	Message:  "Unknown Error: %s",
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
		e = errDieUnknown.WithMessageF(ee)
	}

	// print the error message to standard error in a wrapped way
	if message := e.Error(); message != "" {
		docfmt.Check(message)
		io.StderrWriteWrap(message)
	}

	return e
}
