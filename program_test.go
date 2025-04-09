//spellchecker:words goprogram
package goprogram //nolint:testpackage

//spellchecker:words positionals nolint testpackage

//spellchecker:words bytes path filepath reflect runtime testing time github goprogram exit meta parser pkglib stream testlib
import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/meta"
	"github.com/tkw1536/goprogram/parser"
	"github.com/tkw1536/pkglib/stream"
	"github.com/tkw1536/pkglib/testlib"
)

// This file contains dummy implementations of everything required to assemble a program.
// It is reused across the test suite, however there is no versioning guarantee.
// It may change in a future revision of the test suite.

// Environment of each command is a single string value.
// Parameters to initialize each command is also a string value.
type tEnvironment string
type tParameters string

type echoStruct = struct {
	Arguments []string `description:"arguments"`
}

func makeEchoCommand(name string) iCommand {
	cmd := &tCommand[echoStruct]{
		MDesc: iDescription{
			Command:      name,
			Requirements: func(flag meta.Flag) bool { return true },
		},

		MAfterParse: func() error { return nil },
	}
	cmd.MRun = func(command tCommand[echoStruct], context iContext) error {
		defer func() { command.Positionals.Arguments = nil }() // for the next time

		_, err := context.Printf("%v\n", command.Positionals.Arguments)
		if err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		return nil
	}
	return cmd
}

// makeProgram creates a new program and registers an echo command with it.
func makeProgram() iProgram {
	return iProgram{
		Info: meta.Info{
			BuildVersion: "42.0.0",
			BuildTime:    time.Unix(0, 0).UTC(),

			Executable:  "exe",
			Description: "something something dark side",
		},
	}
}

// tFlags holds a set of dummy global flags.
type tFlags struct {
	GlobalOne string `long:"global-one" short:"a"`
	GlobalTwo string `long:"global-two" short:"b"`
}

// tRequirements is the implementation of the AllowsFlag function.
type tRequirements func(flag meta.Flag) bool

func (t tRequirements) AllowsFlag(flag meta.Flag) bool { return t(flag) }
func (t tRequirements) Validate(args Arguments[tFlags]) error {
	return ValidateAllowedFlags[tFlags](t, args)
}

// instantiated types for the test suite.
type iProgram = Program[tEnvironment, tParameters, tFlags, tRequirements]
type iCommand = Command[tEnvironment, tParameters, tFlags, tRequirements]
type iContext = Context[tEnvironment, tParameters, tFlags, tRequirements]
type iArguments = Arguments[tFlags]
type iDescription = Description[tFlags, tRequirements]

// tCommand represents a sample test suite command.
// It runs the associated private functions, or prints an info message to stdout.
type tCommand[Pos any] struct {
	StdoutMsg string `default:"write to stdout" long:"stdout" short:"o" value-name:"message"`
	StderrMsg string `default:"write to stderr" long:"stderr" short:"e" value-name:"message"`

	Positionals Pos `positional-args:"true"`

	MDesc iDescription

	MAfterParse func() error
	MRun        func(command tCommand[Pos], context iContext) error
}

func (t tCommand[Pos]) Description() iDescription {
	return t.MDesc
}

func (t tCommand[Pos]) AfterParse() error {
	if t.MAfterParse == nil {
		fmt.Println("AfterParse()")
		return nil
	}
	return t.MAfterParse()
}
func (t tCommand[Pos]) Run(ctx iContext) error {
	if t.MRun == nil {
		fmt.Println("Run()")
		return nil
	}
	return t.MRun(t, ctx)
}

// makeTPM_Positionals makes a new command with the provided positional arguments.
func makeTPM_Positionals[Pos any]() iCommand {
	return &tCommand[Pos]{
		MAfterParse: func() error { return nil },
	}
}

func TestProgram_Main(t *testing.T) {
	t.Parallel()

	root := testlib.TempDirAbs(t)
	if err := os.Mkdir(filepath.Join(root, "real"), os.ModeDir&os.ModePerm); err != nil {
		panic(err)
	}

	// define requirements to allow only the Global1 (or any) arguments
	reqOne := tRequirements(func(flag meta.Flag) bool {
		return flag.FieldName == "Global1"
	})

	// define requirements to allow anything
	reqAny := tRequirements(func(flag meta.Flag) bool { return true })

	tests := []struct {
		name string
		args []string

		// should use makeTPM_Positionals[/* positionals type */]()
		// for failure detection, set Args inside the positionals type
		positionals iCommand

		desc       iDescription
		parameters tParameters

		// alias to register (if any)
		alias Alias

		wantStdout string
		wantStderr string
		wantCode   uint8
	}{
		{
			name:        "no arguments",
			args:        []string{},
			positionals: makeTPM_Positionals[struct{}](),

			wantStderr: "Unable to parse arguments: Need at least one argument\n",
			wantCode:   3,
		},

		{
			name:        "unknown general args",
			args:        []string{"--this-flag-does-not-exist", "--", "fake"},
			positionals: makeTPM_Positionals[struct{}](),

			wantStderr: "Unable to parse arguments: Unknown flag `this-flag-does-not-exist'\n",
			wantCode:   3,
		},

		{
			name:        "display help",
			args:        []string{"--help"},
			positionals: makeTPM_Positionals[struct{}](),

			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--] COMMAND [ARGS...]\n\nSomething something dark side\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\n   COMMAND [ARGS...]\n      Command to call. One of \"fake\". See individual commands for more help.\n",
			wantCode:   0,
		},

		{
			name:        "display help, don't run command",
			args:        []string{"--help", "fake", "whatever"},
			positionals: makeTPM_Positionals[struct{}](),

			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--] COMMAND [ARGS...]\n\nSomething something dark side\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\n   COMMAND [ARGS...]\n      Command to call. One of \"fake\". See individual commands for more help.\n",
			wantCode:   0,
		},

		{
			name:        "display version",
			args:        []string{"--version"},
			positionals: makeTPM_Positionals[struct{}](),

			wantStdout: "exe version 42.0.0, built 1970-01-01 00:00:00 +0000 UTC, using " + runtime.Version() + "\n",
			wantCode:   0,
		},

		{
			name:        "command help (1)",
			args:        []string{"fake", "--help"},
			desc:        iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct{}](),

			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--] fake [--stdout|-o message] [--stderr|-e message]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   -o, --stdout message\n       (default write to stdout)\n\n   -e, --stderr message\n       (default write to stderr)\n",
			wantCode:   0,
		},

		{
			name:        "command help (2)",
			args:        []string{"--", "fake", "--help"},
			desc:        iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct{}](),

			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--] fake [--stdout|-o message] [--stderr|-e message]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   -o, --stdout message\n       (default write to stdout)\n\n   -e, --stderr message\n       (default write to stderr)\n",
			wantCode:   0,
		},

		{
			name: "alias help (1)",

			alias: Alias{
				Name:    "alias",
				Command: "fake",
			},

			args:        []string{"alias", "--help"},
			desc:        iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct{}](),

			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--] alias [--] [ARG ...]\n\nAlias for `exe fake`. See `exe fake --help` for detailed help page about fake\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   [ARG ...]\n      Arguments to pass after `exe fake`\n",
			wantCode:   0,
		},

		{
			name: "alias help (2)",

			alias: Alias{
				Name:    "alias",
				Command: "fake",
			},

			args:        []string{"--", "alias", "--help"},
			desc:        iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct{}](),

			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--] alias [--] [ARG ...]\n\nAlias for `exe fake`. See `exe fake --help` for detailed help page about fake\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   [ARG ...]\n      Arguments to pass after `exe fake`\n",
			wantCode:   0,
		},

		{
			name: "alias help (3)",

			alias: Alias{
				Name:        "alias",
				Command:     "fake",
				Args:        []string{"something", "else"},
				Description: "some useful alias",
			},

			args:        []string{"alias", "--help"},
			desc:        iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct{}](),

			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--] alias [--] [ARG ...]\n\nSome useful alias\n\nalias for `exe fake something else`. See `exe fake --help` for detailed help page about fake\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   [ARG ...]\n      Arguments to pass after `exe fake something else`\n",
			wantCode:   0,
		},

		{
			name: "not enough arguments for fake",
			args: []string{"fake"},
			desc: iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct {
				Args []string `required:"1-2"`
			}](),

			wantStderr: "Wrong arguments for fake: The required argument `Args (at least 1 argument)` was not provided\n",
			wantCode:   4,
		},

		{
			name: "'fake' with unknown argument (not allowed)",
			args: []string{"fake", "--argument-not-declared"},
			desc: iDescription{Requirements: reqAny, ParserConfig: parser.Config{IncludeUnknown: false}},
			positionals: makeTPM_Positionals[struct {
				Args []string
			}](),

			wantStdout: "",
			wantStderr: "Wrong arguments for fake: Unknown flag `argument-not-declared'\n",
			wantCode:   4,
		},

		{
			name: "'fake' with unknown argument (allowed)",
			args: []string{"fake", "--argument-not-declared"},
			desc: iDescription{Requirements: reqAny, ParserConfig: parser.Config{IncludeUnknown: true}},
			positionals: makeTPM_Positionals[struct {
				Args []string
			}](),

			wantStdout: "Got Flags: { }\nGot Pos: {[--argument-not-declared]}\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name: "'fake' without global",
			args: []string{"fake", "hello", "world"},
			desc: iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct {
				Args []string `required:"1-2"`
			}](),

			wantStdout: "Got Flags: { }\nGot Pos: {[hello world]}\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name: "'fake' with global (1)",
			args: []string{"-a", "real", "fake", "hello", "world"},
			desc: iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct {
				Args []string `required:"1-2"`
			}](),

			wantStdout: "Got Flags: {real }\nGot Pos: {[hello world]}\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name: "'fake' with global (2)",
			args: []string{"--global-two", "real", "fake", "hello", "world"},
			desc: iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct {
				Args []string `required:"1-2"`
			}](),

			wantStdout: "Got Flags: { real}\nGot Pos: {[hello world]}\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name: "'fake' with disallowed global",
			args: []string{"--global-one", "not-allowed", "fake", "hello", "world"},
			desc: iDescription{Requirements: reqOne},
			positionals: makeTPM_Positionals[struct {
				Args []string `required:"1-2"`
			}](),

			wantStderr: "Wrong number of arguments: \"fake\" takes no \"--global-one\" argument\n",
			wantCode:   4,
		},

		{
			name: "'fake' with allowed and disallowed global",
			args: []string{"--global-one", "one", "--global-two", "two", "fake", "hello", "world"},
			desc: iDescription{Requirements: reqOne},
			positionals: makeTPM_Positionals[struct {
				Args []string `required:"1-2"`
			}](),

			wantStderr: "Wrong number of arguments: \"fake\" takes no \"--global-one\" argument\n",
			wantCode:   4,
		},

		{
			name:        "'fake' with non-global argument with identical name",
			args:        []string{"--", "fake", "--global-one"},
			desc:        iDescription{Requirements: reqAny, ParserConfig: parser.Config{IncludeUnknown: true}},
			positionals: makeTPM_Positionals[struct{ Args []string }](),

			wantStdout: "Got Flags: { }\nGot Pos: {[--global-one]}\nwrite to stdout\n",
			wantStderr: "write to stderr\n", //
			wantCode:   0,
		},

		{
			name:        "'fake' with parsed short argument",
			args:        []string{"fake", "-o", "message"},
			desc:        iDescription{Requirements: reqAny, ParserConfig: parser.Config{IncludeUnknown: true}},
			positionals: makeTPM_Positionals[struct{ Args []string }](),

			wantStdout: "Got Flags: { }\nGot Pos: {[]}\nmessage\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:        "'fake' with non-parsed short argument",
			args:        []string{"fake", "--", "--s", "message"},
			desc:        iDescription{Requirements: reqAny, ParserConfig: parser.Config{IncludeUnknown: true}},
			positionals: makeTPM_Positionals[struct{ Args []string }](),

			wantStdout: "Got Flags: { }\nGot Pos: {[--s message]}\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:        "'fake' with parsed long argument",
			args:        []string{"fake", "--stdout", "message"},
			desc:        iDescription{Requirements: reqAny, ParserConfig: parser.Config{IncludeUnknown: true}},
			positionals: makeTPM_Positionals[struct{ Args []string }](),

			wantStdout: "Got Flags: { }\nGot Pos: {[]}\nmessage\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:        "'fake' with non-parsed long argument",
			args:        []string{"fake", "--", "--stdout", "message"},
			desc:        iDescription{Requirements: reqAny, ParserConfig: parser.Config{IncludeUnknown: true}},
			positionals: makeTPM_Positionals[struct{ Args []string }](),

			wantStdout: "Got Flags: { }\nGot Pos: {[--stdout message]}\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name: "'fake' with failure ",
			args: []string{"fake", "fail"},
			desc: iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct {
				Args []string `required:"1-2"`
			}](),

			wantStdout: "Got Flags: { }\nGot Pos: {[fail]}\nwrite to stdout\n",
			wantStderr: "write to stderr\nTest failure\n",
			wantCode:   1,
		},

		{
			name:        "'notExistent' command",
			args:        []string{"notExistent"},
			positionals: makeTPM_Positionals[struct{}](),

			wantStderr: "Unknown command: Must be one of \"fake\"\n",
			wantCode:   2,
		},

		{
			name: "'notExistent' command (with alias)",

			alias: Alias{
				Name:    "alias",
				Command: "fake",
			},

			args:        []string{"notExistent"},
			positionals: makeTPM_Positionals[struct{}](),

			wantStderr: "Unknown command: Must be one of \"fake\"\n",
			wantCode:   2,
		},

		{
			name: "'alias' without args",
			args: []string{"alias", "hello", "world"},

			alias: Alias{
				Name:    "alias",
				Command: "fake",
			},

			desc:        iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct{ Args []string }](),
			wantStdout:  "Got Flags: { }\nGot Pos: {[hello world]}\nwrite to stdout\n",
			wantStderr:  "write to stderr\n",
			wantCode:    0,
		},

		{
			name: "'alias' with args",
			args: []string{"alias", "world"},

			alias: Alias{
				Name:    "alias",
				Command: "fake",
				Args:    []string{"hello"},
			},

			desc:        iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct{ Args []string }](),

			wantStdout: "Got Flags: { }\nGot Pos: {[hello world]}\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name: "recursive 'fake' alias ",
			args: []string{"fake", "world"},

			alias: Alias{
				Name:    "fake",
				Command: "fake",
				Args:    []string{"hello"},
			},

			desc:        iDescription{Requirements: reqAny},
			positionals: makeTPM_Positionals[struct{ Args []string }](),

			wantStdout: "Got Flags: { }\nGot Pos: {[hello world]}\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// create buffers for input and output
			var stdoutBuffer bytes.Buffer
			var stderrBuffer bytes.Buffer
			stream := stream.NewIOStream(&stdoutBuffer, &stderrBuffer, nil)

			// var fakeCommand *tCommand[T]
			fakeCommand := reflect.ValueOf(tt.positionals).Elem()

			// tt.MDesc = ...
			fakeCommand.FieldByName("MDesc").Set(reflect.ValueOf(tt.desc))

			// tt.MDesc.Command = ...
			fakeCommand.FieldByName("MDesc").FieldByName("Command").Set(reflect.ValueOf("fake"))

			// tt.MRun = ...
			MRun := fakeCommand.FieldByName("MRun")
			MRun.Set(reflect.MakeFunc(MRun.Type(), func(args []reflect.Value) (results []reflect.Value) {
				err := (func(command reflect.Value, context iContext) error {
					pos := command.FieldByName("Positionals")

					_, _ = context.Printf("Got Flags: %s\n", context.Args.Flags)
					_, _ = context.Printf("Got Pos: %v\n", pos.Interface())

					_, _ = context.Println(command.FieldByName("StdoutMsg").Interface())
					_, _ = context.EPrintln(command.FieldByName("StderrMsg").Interface())

					command.FieldByName("Positionals").FieldByName("Args")

					// fail when requested to fail
					if argField := pos.FieldByName("Args"); argField.IsValid() {
						pos, ok := argField.Interface().([]string)
						if ok && len(pos) > 0 && pos[0] == "fail" {
							return exit.Error{ExitCode: exit.ExitGeneric, Message: "test failure"}
						}
					}

					return nil
				})(args[0], args[1].Interface().(iContext))

				var rErr = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())
				if err != nil {
					rErr = reflect.ValueOf(err)
				}

				return []reflect.Value{rErr}
			}))

			program := makeProgram()
			program.Register(tt.positionals)

			if tt.alias.Name != "" {
				program.RegisterAlias(tt.alias)
			}

			// run the program
			ret := exit.AsError(program.Main(stream, tt.parameters, tt.args))

			// check all the error values
			gotCode := uint8(ret.ExitCode)
			gotStdout := stdoutBuffer.String()
			gotStderr := stderrBuffer.String()

			if gotCode != tt.wantCode {
				t.Errorf("Program.Main() code = %v, wantCode %v", gotCode, tt.wantCode)
			}

			if gotStdout != tt.wantStdout {
				t.Errorf("Program.Main() stdout = %q, wantStdout %q", gotStdout, tt.wantStdout)
			}

			if gotStderr != tt.wantStderr {
				t.Errorf("Program.Main() stderr = %q, wantStderr %q", gotStderr, tt.wantStderr)
			}
		})
	}
}
