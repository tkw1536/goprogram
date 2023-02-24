package exit

import (
	"fmt"

	"github.com/tkw1536/pkglib/docfmt"
)

// Error represents any error state by a program.
// It implements the builtin error interface.
//
// The zero value represents that no error occured and is ready to use.
type Error struct {
	// Exit code of the program (if applicable)
	ExitCode
	// Message for this error
	// Messages should pass docfmt.Validate
	Message string

	// underlying wrapped error, if any
	err error
}

// Unwrap unwraps this error, if any
func (err Error) Unwrap() error {
	return err.err
}

// Error returns the error message belonging to this error.
func (err Error) Error() string {
	docfmt.AssertValid(err.Message)
	return err.Message
}

// AsError asserts that err is either nil or of type Error and returns it.
// When err is nil, the zero value of type Error is returned.
//
// If err is not nil and not of type Error, calls panic().
func AsError(err error) Error {
	switch e := err.(type) {
	case nil:
		return Error{}
	case Error:
		return e
	}
	panic("AsError: err must be nil or Error")
}

// WithMessage returns a copy of this error with the same Code but different Message.
//
// The new message is the message passed as an argument.
func (err Error) WithMessage(message string) Error {
	docfmt.AssertValid(err.Message)
	return Error{
		ExitCode: err.ExitCode,
		Message:  message,
	}
}

// WithMessageF returns a copy of this error with the same Code but different Message.
// The new message is the current message, formatted using a call to SPrintf and the arguments.
func (err Error) WithMessageF(args ...any) Error {
	docfmt.AssertValid(err.Message)
	docfmt.AssertValidArgs(args...)
	return err.WithMessage(fmt.Sprintf(err.Message, args...))
}

// Wrap creates a new error with same exit code, wrapping the inner error.
// When inner is nil, returns an empty error.
//
// The message of the new error will contain the Error() result of the inner error.
func (err Error) Wrap(inner error) Error {
	if inner == nil {
		return Error{}
	}
	err.Message = fmt.Sprintf("%s: %s", err.Message, inner)
	err.err = inner
	return err
}
