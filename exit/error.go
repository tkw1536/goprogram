package exit

import (
	"fmt"

	"github.com/tkw1536/goprogram/lib/doccheck"
)

// Error represents any error state by a program.
// It implements the builtin error interface.
//
// The zero value represents that no error occured and is ready to use.
type Error struct {
	// Exit code of the program (if applicable)
	ExitCode
	// Message for this error
	Message string
}

// Error returns the error message belonging to this error.
func (err Error) Error() string {
	doccheck.Check(err.Message)
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
	doccheck.Check(err.Message)
	return Error{
		ExitCode: err.ExitCode,
		Message:  message,
	}
}

// WithMessageF returns a copy of this error with the same Code but different Message.
// The new message is the current message, formatted using a call to SPrintf and the arguments.
func (err Error) WithMessageF(args ...interface{}) Error {
	doccheck.Check(err.Message)
	return err.WithMessage(fmt.Sprintf(err.Message, args...))
}
