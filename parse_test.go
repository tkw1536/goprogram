//spellchecker:words goprogram
package goprogram

//spellchecker:words reflect testing
import (
	"reflect"
	"testing"
)

func TestArguments_parseProgramFlags(t *testing.T) {
	type args struct {
		argv []string
	}
	tests := []struct {
		name       string
		args       args
		wantParsed iArguments
		wantErr    error
	}{
		{"no arguments", args{[]string{}}, iArguments{}, errParseArgsNeedOneArgument},
		{"command without arguments", args{[]string{"cmd"}}, iArguments{Command: "cmd", pos: []string{}}, nil},

		{"help with command (1)", args{[]string{"--help", "cmd"}}, iArguments{Universals: Universals{Help: true}, pos: []string{"cmd"}}, nil},
		{"help with command (2)", args{[]string{"-h", "cmd"}}, iArguments{Universals: Universals{Help: true}, pos: []string{"cmd"}}, nil},

		{"help without command (1)", args{[]string{"--help"}}, iArguments{Universals: Universals{Help: true}, pos: []string{}}, nil},
		{"help without command (2)", args{[]string{"-h"}}, iArguments{Universals: Universals{Help: true}, pos: []string{}}, nil},

		{"version with command (1)", args{[]string{"--version", "cmd"}}, iArguments{Universals: Universals{Version: true}, pos: []string{"cmd"}}, nil},
		{"version with command (2)", args{[]string{"-v", "cmd"}}, iArguments{Universals: Universals{Version: true}, pos: []string{"cmd"}}, nil},

		{"version without command (2)", args{[]string{"--version"}}, iArguments{Universals: Universals{Version: true}, pos: []string{}}, nil},
		{"version without command (3)", args{[]string{"-v"}}, iArguments{Universals: Universals{Version: true}, pos: []string{}}, nil},

		{"command with arguments", args{[]string{"cmd", "a1", "a2"}}, iArguments{Command: "cmd", pos: []string{"a1", "a2"}}, nil},

		{"command with help (1)", args{[]string{"cmd", "help", "a1"}}, iArguments{Command: "cmd", pos: []string{"help", "a1"}}, nil},
		{"command with help (2)", args{[]string{"cmd", "--help", "a1"}}, iArguments{Command: "cmd", pos: []string{"--help", "a1"}}, nil},
		{"command with help (3)", args{[]string{"cmd", "-h", "a1"}}, iArguments{Command: "cmd", pos: []string{"-h", "a1"}}, nil},

		{"command with version (1)", args{[]string{"cmd", "version", "a1"}}, iArguments{Command: "cmd", pos: []string{"version", "a1"}}, nil},
		{"command with version (2)", args{[]string{"cmd", "--version", "a1"}}, iArguments{Command: "cmd", pos: []string{"--version", "a1"}}, nil},
		{"command with version (3)", args{[]string{"cmd", "-v", "a1"}}, iArguments{Command: "cmd", pos: []string{"-v", "a1"}}, nil},

		{"global flag without command (1)", args{[]string{"-a", "stuff"}}, iArguments{}, errParseArgsNeedOneArgument},
		{"global flag without command (2)", args{[]string{"--global-one", "stuff"}}, iArguments{}, errParseArgsNeedOneArgument},

		{"global flag with command (1)", args{[]string{"-a", "stuff", "cmd"}}, iArguments{Command: "cmd", Flags: tFlags{GlobalOne: "stuff"}, pos: []string{}}, nil},
		{"global flag with command (2)", args{[]string{"--global-one", "stuff", "cmd"}}, iArguments{Command: "cmd", Flags: tFlags{GlobalOne: "stuff"}, pos: []string{}}, nil},

		{"global flag with command and arguments (1)", args{[]string{"--global-two", "stuff", "cmd", "a1", "a2"}}, iArguments{Command: "cmd", Flags: tFlags{GlobalTwo: "stuff"}, pos: []string{"a1", "a2"}}, nil},
		{"global flag with command and arguments (2)", args{[]string{"-b", "stuff", "cmd", "a1", "a2"}}, iArguments{Command: "cmd", Flags: tFlags{GlobalTwo: "stuff"}, pos: []string{"a1", "a2"}}, nil},

		{"global looking flag", args{[]string{"--not-a-global-flag", "stuff", "command"}}, iArguments{}, errParseArgsUnknownError.WithMessageF("unknown flag `not-a-global-flag'")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args iArguments
			err := args.parseProgramFlags(tt.args.argv)

			// turn wantErr into a string
			var wantErr string
			if tt.wantErr != nil {
				wantErr = tt.wantErr.Error()
			}

			// turn gotErr into a string
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}

			// compare error messages
			if wantErr != gotErr {
				t.Errorf("Arguments.parseP() error = %#v, wantErr %#v", err, tt.wantErr)
			}

			if tt.wantErr != nil { // ignore checks when an error is returned; we don't care
				return
			}

			if !reflect.DeepEqual(args, tt.wantParsed) {
				t.Errorf("Arguments.parseProgramFlags() args = %#v, wantArgs %#v", args, &tt.wantParsed)
			}
		})
	}
}
