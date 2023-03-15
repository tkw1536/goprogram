package exit

import (
	"github.com/tkw1536/pkglib/docfmt"
	"github.com/tkw1536/pkglib/stream"
)

var errDieUnknown = Error{
	ExitCode: ExitGeneric,
	Message:  "unknown error",
}

// Die prints a non-nil err to io.Stderr and returns an error of type Error or nil.
//
// When error is non-nil, this function first turns err into type Error.
// Then if err.Message is not the empty string, it prints it to io.Stderr with wrapping.
//
// If err is nil, it does nothing and returns nil.
func Die(str stream.IOStream, err error) error {
	var e Error
	switch ee := err.(type) {
	case nil:
		return nil
	case Error:
		e = ee
	default:
		e = errDieUnknown.Wrap(ee)
	}

	// print the error message to standard error in a wrapped way
	if message := e.Error(); message != "" {
		if stream.IsNullWriter(str.Stderr) {
			docfmt.Format(message)
			return e
		}
		str.StderrWriteWrap(docfmt.Format(message))
	}

	return e
}
