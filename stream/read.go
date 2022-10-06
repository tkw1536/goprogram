package stream

import "github.com/tkw1536/goprogram/lib/nobufio"

// ReadLine is like [nobufio.ReadLine] on the standard input
func (io IOStream) ReadLine() (string, error) {
	return nobufio.ReadLine(io.Stdin)
}

// ReadPassword is like [nobufio.ReadPassword] on the standard input
func (io IOStream) ReadPassword() (string, error) {
	return nobufio.ReadPassword(io.Stdin)
}

// ReadPasswordStrict is like [nobufio.ReadPasswordStrict] on the standard input
func (io IOStream) ReadPasswordStrict() (string, error) {
	return nobufio.ReadPasswordStrict(io.Stdin)
}
