package stream

import (
	"errors"
	"io"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// ReadLine reads the current line from the provided reader.
//
// A line is considered to end when one of the following is encountered: '\r\n', '\n' or EOF or '\r' followed by EOF.
// Note that only a '\r' is not considered an end-of-line.
//
// The returned line never contains the end-of-line markers, such as '\n' or '\r\n'.
// A line may be empty, however when only EOF is read, returns "", EOF.
func ReadLine(reader io.Reader) (value string, err error) {
	var builder strings.Builder // buffer for the string to construct
	var lastR bool              // delay writing a '\r', in case it is followed by an '\n'
	var readSomething bool
	for {
		// read the next valid rune
		r, err := readRune(reader)
		if err == io.EOF { // at EOF, we are done!
			break
		}
		readSomething = true
		if err != nil { // unknown reading error => bail out
			return "", err
		}
		if r == '\n' { // \n or \r\n
			break
		}

		if lastR {
			// flag is set, but we didn't encounter a '\n' or EOF.
			// so we need to write it back to the buffer
			if _, err := builder.WriteRune('\r'); err != nil {
				return "", err
			}
			lastR = false
		}
		if r == '\r' {
			lastR = true
			continue
		}

		// store it to the builder
		if _, err := builder.WriteRune(r); err != nil {
			return "", err
		}
	}

	// if we didn't read anything, return EOF!
	if !readSomething {
		return "", io.EOF
	}

	// make it a string
	return builder.String(), nil
}

// readRune reads a single valid rune from reader.
func readRune(reader io.Reader) (r rune, err error) {
	var runeBuffer []byte

	var count int
	for !utf8.FullRune(runeBuffer) {
		// expand the rune buffer
		runeBuffer = append(runeBuffer, 0)

		// read the next byte into it into or bail out!
		if _, err = reader.Read(runeBuffer[count:]); err != nil {
			return
		}
		count++
	}

	// decode the rune!
	r, _ = utf8.DecodeRune(runeBuffer)
	return r, nil
}

// ReadPassword is like ReadLine, except that it turns off terminal echo.
// When standard input is not a terminal, behaves like ReadLine()
func ReadPassword(reader io.Reader) (value string, err error) {
	value, err = ReadPasswordStrict(reader)
	if err == ErrNoTerminal {
		return ReadLine(reader)
	}
	return
}

// ErrNoTerminal is returned by ReadPasswordStrict() when stdin is not a terminal
var ErrNoTerminal = errors.New("ReadPasswordStrict: Stdin is not a terminal")

// ReadPasswordSrict is like ReadPassword, except that when stdin is not a terminal, returns ErrNoTerminal.
func ReadPasswordStrict(reader io.Reader) (value string, err error) {
	// check if stdin is a terminal
	file, ok := reader.(interface{ Fd() uintptr })
	if !ok || !term.IsTerminal(int(file.Fd())) {
		return "", ErrNoTerminal
	}

	// read the bytes
	bytes, err := term.ReadPassword(int(file.Fd()))
	return string(bytes), err
}

// ReadLine is like ReadLine(io.Stdin)
func (io IOStream) ReadLine() (string, error) {
	return ReadLine(io.Stdin)
}

// ReadPassword is like ReadPassword(io.Stdin)
func (io IOStream) ReadPassword() (string, error) {
	return ReadPassword(io.Stdin)
}

// ReadPasswordStrict is like ReadPasswordStrict(io.Stdin)
func (io IOStream) ReadPasswordStrict() (string, error) {
	return ReadPasswordStrict(io.Stdin)
}
