package status

import (
	"bytes"
	"io"
	"sync"
)

// LineBuffer is an [io.Writer] that calls a function Line for every newline-delimited line written to it.
// Do not copy a non-zero LineBuffer.
type LineBuffer struct {
	m      sync.Mutex
	buffer bytes.Buffer

	// Line is called once a complete newline-terminated line has been written to this LineBuffer.
	// It is called for each line, in the correct order.
	// The parameter has a trailing '\r\n' or '\n' trimmed.
	//
	// Write methods block until the Line function has returned.
	// Therefor Line must not trigger another Write into the LineBuffer.
	Line func(line string)
}

// Write writes b into the internal buffer.
// When this completes one or more lines, calls Line appropriatly.
func (lb *LineBuffer) Write(b []byte) (int, error) {
	lb.m.Lock()
	defer lb.m.Unlock()

	defer lb.flush()
	return lb.buffer.Write(b)
}

// WriteByte is like [Write], but takes a single byte.
func (lb *LineBuffer) WriteByte(b byte) error {
	lb.m.Lock()
	defer lb.m.Unlock()

	defer lb.flush()
	return lb.buffer.WriteByte(b)
}

// WriteRune is like [Write], but takes a single rune
func (lb *LineBuffer) WriteRune(r rune) (int, error) {
	lb.m.Lock()
	defer lb.m.Unlock()

	defer lb.flush()
	return lb.buffer.WriteRune(r)
}

// WriteString is like [Write], but takes a string
func (lb *LineBuffer) WriteString(s string) (int, error) {
	lb.m.Lock()
	defer lb.m.Unlock()

	defer lb.flush()
	return lb.buffer.WriteString(s)
}

func (lb *LineBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	lb.m.Lock()
	defer lb.m.Unlock()

	defer lb.flush()
	return lb.buffer.ReadFrom(r)

}

// runeR and runeN represent the bytes corresponding to '\r' and '\n' respecitively.
const runeR byte = '\r'
const runeN byte = '\n'

// flush takes any completed lines in the internal buffer and calls the Line function
func (lb *LineBuffer) flush() {
	for {
		// find the index of any '\n'
		index := bytes.IndexByte(lb.buffer.Bytes(), runeN)
		if index == -1 {
			return
		}

		// take the line, and trim any trailing '\r'
		line := lb.buffer.Next(index + 1)[:index]
		if index > 0 && line[index-1] == runeR {
			line = line[:index-1]
		}

		// call the line function!
		lb.Line(string(line))
	}
}
