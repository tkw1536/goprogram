//spellchecker:words exit
package exit

//spellchecker:words github pkglib docfmt stream
import (
	"fmt"

	"github.com/tkw1536/pkglib/docfmt"
	"github.com/tkw1536/pkglib/stream"
)

// Die prints a non-nil err to io.Stderr and returns an error of type Error or nil.
//
// When error is non-nil, this function first turns err into type Error.
// If error is of type Error (or wraps an Error) uses the same procedure as [AsError].
// If not, creates a new Error wrapping in the error and a generic error message.
//
// Then if err.Message is not the empty string, it prints it to io.Stderr.
//
// If err is nil, it does nothing and returns nil.
func Die(str stream.IOStream, err error) error {
	// fast case: not an error
	if err == nil {
		return nil
	}

	ourError, ok := asError(err)
	if !ok {
		ourError = Error{
			ExitCode: ExitGeneric,
			Message:  fmt.Sprintf("unknown error: %s", err),
			err:      err,
		}
	}

	// print the error message to standard error in a wrapped way
	if message := fmt.Sprint(ourError); message != "" {
		if stream.IsNullWriter(str.Stderr) {
			docfmt.Format(message)
			return ourError
		}
		_, _ = str.EPrintln(docfmt.Format(message)) // no way to report the failure
	}

	return ourError
}
