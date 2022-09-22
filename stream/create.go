package stream

import "os"

// FromNil creates a new IOStream that silences all output and provides no input.
func FromNil() IOStream {
	return NewIOStream(nil, nil, nil, 0)
}

// FromEnv creates a new IOStream using the environment.
//
// The Stdin, Stdout and Stderr streams are used from the os package.
func FromEnv() IOStream {
	return IOStream{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		wrap:   ioDefaultWrap,
	}
}

// FromDebug returns a stream that behaves like [FromNil] when debugging is disabled.
// When debugging is enabled (by building with the "debug" tag), behaves like [FromEnv] instead.
func FromDebug() IOStream {
	if debugEnabled {
		return FromEnv()
	}
	return FromNil()
}
