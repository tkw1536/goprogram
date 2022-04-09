package goprogram

import (
	"reflect"
	"testing"

	"github.com/tkw1536/goprogram/meta"
	"github.com/tkw1536/goprogram/parser"
)

func TestProgram_MainUsage(t *testing.T) {
	program := makeProgram()

	program.Register(makeEchoCommand("a"))
	program.Register(makeEchoCommand("c"))
	program.Register(makeEchoCommand("b"))

	got := program.MainUsage()
	want := meta.Meta{Executable: "exe", Command: "", Description: "something something dark side", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}, {FieldName: "GlobalOne", Short: []string{"a"}, Long: []string{"global-one"}, Required: false, Value: "", Usage: "", Default: ""}, {FieldName: "GlobalTwo", Short: []string{"b"}, Long: []string{"global-two"}, Required: false, Value: "", Usage: "", Default: ""}}, CommandFlags: []meta.Flag(nil), Positionals: []meta.Positional(nil), Commands: []string{"a", "b", "c"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Program.MainUsage() = %#v, want %#v", got, want)
	}
}

// makeTPM_Positionals makes a new parser with the provided positional arguments
func makeTPCU_Positionals[Pos any]() parser.Parser {
	return parser.Config{}.NewCommandParser(&struct {
		Boolean     bool `short:"b" value-name:"random" long:"bool" description:"a random boolean argument with short"`
		Int         int  `long:"int" value-name:"dummy" description:"a dummy integer flag" default:"12"`
		Positionals Pos  `positional-args:"true"`
	}{})
}

func TestProgram_CommandUsage(t *testing.T) {

	program := makeProgram()

	// define requirements to allow only the Global1 (or any) arguments
	reqOne := tRequirements(func(flag meta.Flag) bool {
		return flag.FieldName == "Global1"
	})

	// define requirements to allow anything
	reqAny := tRequirements(func(flag meta.Flag) bool { return true })

	type args struct {
		Command     string
		Description string
		Requirement tRequirements
		Positionals parser.Parser // should use newParser[/* positionals struct */]()
	}
	tests := []struct {
		name string
		args args
		want meta.Meta
	}{
		{
			"command without args and allowing all globals",
			args{Command: "cmd", Requirement: reqAny, Positionals: makeTPCU_Positionals[struct{}]()},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}, {FieldName: "GlobalOne", Short: []string{"a"}, Long: []string{"global-one"}, Required: false, Value: "", Usage: "", Default: ""}, {FieldName: "GlobalTwo", Short: []string{"b"}, Long: []string{"global-two"}, Required: false, Value: "", Usage: "", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positionals: []meta.Positional{}, Commands: []string(nil)},
		},

		{
			"command without args and allowing only global1",
			args{Command: "cmd", Requirement: reqOne, Positionals: makeTPCU_Positionals[struct {
				Meta string `description:"usage" positional-arg-name:"META"`
			}]()},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positionals: []meta.Positional{{Value: "META", Usage: "usage", Min: 0, Max: 1}}, Commands: []string(nil)},
		},

		{
			"command with max finite args",
			args{Command: "cmd", Requirement: reqOne, Positionals: makeTPCU_Positionals[struct {
				Meta []string `description:"usage" positional-arg-name:"META" required:"0-4"`
			}]()},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positionals: []meta.Positional{{Value: "META", Usage: "usage", Min: 0, Max: 4}}, Commands: []string(nil)},
		},

		{
			"command with finite args",
			args{Command: "cmd", Requirement: reqOne, Positionals: makeTPCU_Positionals[struct {
				Meta []string `description:"usage" positional-arg-name:"META" required:"1-2"`
			}]()},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positionals: []meta.Positional{{Value: "META", Usage: "usage", Min: 1, Max: 2}}, Commands: []string(nil)},
		},

		{
			"command with infinite args",
			args{Command: "cmd", Requirement: reqOne, Positionals: makeTPCU_Positionals[struct {
				Meta []string `description:"usage" positional-arg-name:"META" required:"1"`
			}]()},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positionals: []meta.Positional{{Value: "META", Usage: "usage", Min: 1, Max: -1}}, Commands: []string(nil)},
		},

		{
			"command with description",
			args{Command: "cmd", Description: "A fake command", Requirement: reqOne, Positionals: makeTPCU_Positionals[struct {
				Meta []string `description:"usage" positional-arg-name:"META" required:"1"`
			}]()},
			meta.Meta{Executable: "exe", Command: "cmd", Description: "A fake command", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}}, CommandFlags: []meta.Flag{{FieldName: "Boolean", Short: []string{"b"}, Long: []string{"bool"}, Required: false, Value: "random", Usage: "a random boolean argument with short", Default: ""}, {FieldName: "Int", Short: []string(nil), Long: []string{"int"}, Required: false, Value: "dummy", Usage: "a dummy integer flag", Default: "12"}}, Positionals: []meta.Positional{{Value: "META", Usage: "usage", Min: 1, Max: -1}}, Commands: []string(nil)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := tt.args.Positionals

			context := iContext{
				Args: iArguments{
					Command: tt.args.Command,
				},

				parser: parser,

				Description: iDescription{
					Command:      tt.args.Command,
					Description:  tt.args.Description,
					Requirements: tt.args.Requirement,
				},
			}
			got := program.CommandUsage(context)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Program.CommandUsage() = %#v\n\n\n, want %#v", got, tt.want)
			}
		})
	}
}

func TestProgram_AliasUsage(t *testing.T) {
	program := makeProgram()
	program.Register(makeEchoCommand("a"))

	alias := Alias{Name: "nice", Command: "a", Args: []string{"nice", "command"}, Description: "Do one nice thing"}
	program.RegisterAlias(alias)

	context := iContext{
		Description: iDescription{
			Requirements: func(flag meta.Flag) bool { return true },
		},
	}

	got := program.AliasUsage(context, alias)
	want := meta.Meta{Executable: "exe", Command: "nice", Description: "Do one nice thing\n\nAlias for `exe a nice command`. See `exe a --help` for detailed help page about a. ", GlobalFlags: []meta.Flag{{FieldName: "Help", Short: []string{"h"}, Long: []string{"help"}, Required: false, Value: "", Usage: "Print a help message and exit", Default: ""}, {FieldName: "Version", Short: []string{"v"}, Long: []string{"version"}, Required: false, Value: "", Usage: "Print a version message and exit", Default: ""}, {FieldName: "GlobalOne", Short: []string{"a"}, Long: []string{"global-one"}, Required: false, Value: "", Usage: "", Default: ""}, {FieldName: "GlobalTwo", Short: []string{"b"}, Long: []string{"global-two"}, Required: false, Value: "", Usage: "", Default: ""}}, CommandFlags: []meta.Flag(nil), Positionals: []meta.Positional{{Value: "ARG", Usage: "Arguments to pass after `exe a nice command`.", Min: 0, Max: -1}}, Commands: []string(nil)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Program.AliasUsage() = %#v, want %#v", got, want)
	}
}
