//spellchecker:words exit
package exit

//spellchecker:words github pkglib docfmt
import (
	"errors"
)

//spellchecker:words nolint errorlint

// NewErrorWithCode creates a new error that additionally holds the given exit code.
func NewErrorWithCode(message string, code ExitCode) error {
	return &codeError{message: message, code: code}
}

type codeError struct {
	code    ExitCode
	message string
}

func (err *codeError) Error() string {
	return err.message
}

// CodeFromError returns the ExitCode contained in error, if any.
// The exit code is found by [errors.As] unwrapping into an error created by this package.
//
// When err is nil, returns [ExitZero] and zero.
// When err does not hold any [Error]s, returns [ExitGeneric] and false.
func CodeFromError(err error) (code ExitCode, ok bool) {
	if err == nil {
		return ExitZero, true
	}
	var codeErr *codeError
	if !errors.As(err, &codeErr) {
		return ExitGeneric, false
	}
	return codeErr.code, true
}
