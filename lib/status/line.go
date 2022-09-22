package status

import (
	"strings"
	"sync"
)

// LineBuffer is an [io.Writer] that calls Line for every newline-delimited line written to it.
// Do not copy a non-zero LineBuffer.
type LineBuffer struct {
	l       sync.Mutex // l protects builder
	builder strings.Builder

	// Line is called once a complete newline-terminated line has been written to this LineBuffer.
	// It is called for each line, in the correct order.
	//
	// The parameter will be the line written, with terminating '\r\n' or '\n' stripped.
	Line func(line string)
}

// Write writes b into the internal buffer.
// When this completes one or more lines, calls Line appropriatly.
func (lb *LineBuffer) Write(b []byte) (int, error) {
	lb.l.Lock()
	defer lb.l.Unlock()

	defer lb.flush() // flush when done!

	return lb.builder.Write(b)
}

// WriteByte is like [Write], but takes a single byte.
func (lb *LineBuffer) WriteByte(b byte) error {
	lb.l.Lock()
	defer lb.l.Unlock()

	defer lb.flush()

	return lb.builder.WriteByte(b)
}

// WriteRune is like [Write], but takes a single rune
func (lb *LineBuffer) WriteRune(r rune) (int, error) {
	lb.l.Lock()
	defer lb.l.Unlock()

	defer lb.flush()

	return lb.builder.WriteRune(r)
}

// WriteString is like [Write], but takes a string
func (lb *LineBuffer) WriteString(s string) (int, error) {
	lb.l.Lock()
	defer lb.l.Unlock()

	defer lb.flush()

	return lb.builder.WriteString(s)
}

// flush takes any completed lines in the internal buffer and flushes them by calling [Line].
// Returns the number of calls to line performed.
func (lb *LineBuffer) flush() (count int) {
	// grab the text from the buffer, and while there are newlines in it
	// trim off the first line.
	text := lb.builder.String()
	var line string
	for strings.ContainsRune(text, '\n') {
		line, text, _ = strings.Cut(text, "\n")
		lb.Line(strings.TrimSuffix(line, "\r"))
		count++
	}

	// reset the builder to contain only the remaining text.
	// except if there wasn't any change to the text.
	if count != 0 {
		lb.builder.Reset()
		lb.builder.WriteString(text)
	}

	return
}
