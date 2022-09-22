package status

import "fmt"

func ExampleLineBuffer() {
	// create a new line buffer
	buffer := LineBuffer{
		Line: func(line string) {
			fmt.Printf("Line(%q)\n", line)
		},
	}

	// write some text into it, calling Line() with each completed line
	buffer.WriteString("line 1\npartial")
	buffer.WriteString(" line 2\n\n")

	// Output: Line("line 1")
	// Line("partial line 2")
	// Line("")
}
