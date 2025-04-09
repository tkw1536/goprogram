//spellchecker:words exit
package exit

//spellchecker:words github pkglib docfmt
import (
	"errors"
	"fmt"

	"github.com/tkw1536/pkglib/docfmt"
)

//spellchecker:words nolint errorlint

// Error represents any error state by a program.
// It implements the builtin error interface.
//
// The zero value represents that no error occurred and is ready to use.
type Error struct {
	// Exit code of the program (if applicable)
	ExitCode
	// Message for this error
	// Messages should pass docfmt.Validate
	Message string

	// underlying wrapped error, if any
	err error
}

// Unwrap unwraps this error, if any.
func (err Error) Unwrap() error {
	return err.err
}

// Error returns the error message belonging to this error.
func (err Error) Error() string {
	return err.Message
}

// AsError asserts that err is either nil, wraps an error of type Error, or is of type Error itself.
// When failing the precondition, panic()s.
//
// If nil or of type Error, returns err unchanged.
// When wrapping an Error, returns a new Error object with the appropriate exit code that wraps the original.
//
// If err is not nil and not of type Error, calls panic().
func AsError(err error) Error {
	ourError, ok := asError(err)
	if !ok && err != nil {
		panic("AsError: err must be nil or wrap type Error")
	}
	return ourError
}

// asError tries to turn error into an Error maintaining exit code.
// If error is nil, return as is
//
// - if an Error, return that Error and true.
// - if wrapping an Error, lift the exit code and wrapping.
// - in all other cases: return Error{}, false.
func asError(err error) (Error, bool) {
	if err == nil {
		return Error{}, false
	}

	// when nil, or an error, return as is!
	if ourError, ok := err.(Error); ok { //nolint:errorlint
		return ourError, true
	}

	var wrapped Error
	if errors.As(err, &wrapped) {
		return Error{
			ExitCode: wrapped.ExitCode,
			Message:  err.Error(),
			err:      err,
		}, true
	}

	return Error{}, false
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

// WrapError creates a new Error with same exit code, wrapping the inner error.
// When inner is nil, returns nil.
// This function will return either nil, or an error of type Error.
//
// The message of the new error will contain the Error() result of the inner error.
func (err Error) WrapError(inner error) error {
	if inner == nil {
		return nil
	}
	err.Message = fmt.Sprintf("%s: %s", err.Message, inner)
	err.err = inner
	return err
}
