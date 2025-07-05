//spellchecker:words goprogram
package goprogram //nolint:testpackage

//spellchecker:words reflect testing github pkglib stream
import (
	"reflect"
	"testing"

	"go.tkw01536.de/pkglib/stream"
)

//spellchecker:words nolint testpackage

// Register a command for a program.
// See the test suite for instantiated types.
func ExampleCommand() {
	// create a new program that only has an echo command
	// this code is reused across the test suite, hence not shown here.
	p := makeProgram()
	p.Register(makeEchoCommand("echo"))

	// Execute the command with some arguments
	_ = p.Main(stream.FromEnv(), "", []string{"echo", "hello", "world"})

	// Output: [hello world]
}

func TestProgram_Commands(t *testing.T) {
	t.Parallel()

	p := makeProgram()
	p.Register(makeEchoCommand("a"))
	p.Register(makeEchoCommand("c"))
	p.Register(makeEchoCommand("b"))

	got := p.Commands()
	want := []string{"a", "b", "c"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Program.Commands() = %v, want = %v", got, want)
	}
}

func TestProgram_FmtCommands(t *testing.T) {
	t.Parallel()

	p := makeProgram()
	p.Register(makeEchoCommand("a"))
	p.Register(makeEchoCommand("c"))
	p.Register(makeEchoCommand("b"))

	got := p.FmtCommands()
	want := `"a", "b", "c"`

	if got != want {
		t.Errorf("Program.FmtCommands() = %v, want = %v", got, want)
	}
}
