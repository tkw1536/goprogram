//spellchecker:words goprogram
package goprogram

//spellchecker:words reflect slices github goprogram exit meta parser pkglib reflectx
import (
	"reflect"
	"slices"

	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/meta"
	"github.com/tkw1536/goprogram/parser"
	"github.com/tkw1536/pkglib/reflectx"
)

//spellchecker:words jessevdk

// Command represents a command associated with a program.
// It takes the same type parameters as a program.
//
// Each command is first initialized using any means by the user.
// Next, it is registered with a program using the Program.Register Method.
// Once a program is called with this command, the arguments for it are parsed by making use of the Description.
// Eventually the Run method of this command is invoked, using an appropriate new Context.
// See Context for details on the latter.
//
// A command should be implemented as a struct or pointer to a struct.
// Typically command contains state that represents the parsed options.
// Flag parsing is implemented using the "github.com/jessevdk/go-flags" package.
//
// To ensure consistency between runs, a command is copied before being executed.
// When command is a pointer, the underlying data (not the pointer to it) will be copied.
type Command[E any, P any, F any, R Requirement[F]] interface {
	// Run runs this command in the given context.
	//
	// It is called only once and must return either nil or an error of type Error.
	Run(context Context[E, P, F, R]) error

	// Description returns a description of this command.
	// It may be called multiple times.
	Description() Description[F, R]
}

// AfterParseCommand represents a command with an AfterParse function
type AfterParseCommand[E any, P any, F any, R Requirement[F]] interface {
	Command[E, P, F, R]

	// AfterParse is called after arguments have been parsed, but before the command is being run.
	// It may perform additional argument checking and should return an error if needed.
	//
	// It is called only once and must return either nil or an error of type Error.
	AfterParse() error
}

// Description describes a command, and specifies any potential requirements.
type Description[F any, R Requirement[F]] struct {
	// Command and Description the name and human-readable description of this command.
	// Command must not be taken by any other command registered with the corresponding program.
	Command     string
	Description string

	ParserConfig parser.Config // information about how to configure a parser for this command

	// Requirements on the environment to be able to run the command
	Requirements R
}

// Requirement describes a requirement on a type of Flags F.
//
// A caller may use EmptyRequirement if no such abstraction is desired.
type Requirement[F any] interface {
	// AllowsFlag checks if the provided flag may be passed to fullfil this requirement
	// By default it is used only for help page generation, and may be inaccurate.
	AllowsFlag(flag meta.Flag) bool

	// Validate validates if this requirement is fulfilled for the provided global flags.
	// It should return either nil, or an error of type exit.Error.
	//
	// Validate does not take into account AllowsOption, see ValidateAllowedOptions.
	Validate(arguments Arguments[F]) error
}

// EmptyRequirement represents a requirement that allows any flag and validates all arguments
type EmptyRequirement[F any] struct{}

func (EmptyRequirement[F]) AllowsFlag(meta.Flag) bool   { return true }
func (EmptyRequirement[F]) Validate(Arguments[F]) error { return nil }

// Register registers a command c with this program.
//
// It expects that the command does not have a name that is already taken.
func (p *Program[E, P, F, R]) Register(c Command[E, P, F, R]) {
	if p.commands == nil {
		p.commands = make(map[string]Command[E, P, F, R])
	}

	Name := c.Description().Command

	if _, ok := p.commands[Name]; ok {
		panic("Register(): Command already registered")
	}

	p.commands[Name] = c
}

// Commands returns a list of known commands
func (p Program[E, P, F, R]) Commands() []string {
	commands := make([]string, 0, len(p.commands))
	for cmd := range p.commands {
		commands = append(commands, cmd)
	}
	slices.Sort(commands)
	return commands
}

// Command returns the command with the provided name and if it exists
func (p Program[E, P, F, R]) Command(name string) (Command[E, P, F, R], bool) {
	// NOTE: This function is not tested because it is so trivial
	cmd, ok := p.commands[name]
	if ok {
		cmd, _ = reflectx.CopyInterface(cmd)
	}
	return cmd, ok
}

// FmtCommands returns a human readable string describing the commands.
// See also Commands.
func (p Program[E, P, F, R]) FmtCommands() string {
	return meta.JoinCommands(p.Commands())
}

var errTakesNoArgument = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "wrong number of arguments: %q takes no %q argument",
}

// Validate validates that every flag f in args.flags either passes the AllowsOption method of the given requirement, or has the zero value.
// If this is not the case returns an error of type ValidateAllowedFlags.
//
// This function is intended to be used to implement the validate method of a Requirement.
func ValidateAllowedFlags[F any](r Requirement[F], args Arguments[F]) error {
	fVal := reflect.ValueOf(args.Flags)

	for _, flag := range parser.AllFlags[F]() {
		if r.AllowsFlag(flag) {
			continue
		}

		v := fVal.FieldByName(flag.FieldName)
		if !v.IsZero() { // flag was set!
			name := flag.Long
			if len(name) == 0 {
				name = []string{""}
			}
			return errTakesNoArgument.WithMessageF(args.Command, "--"+name[0])
		}
	}

	return nil

}

var universalOpts = parser.AllFlags[Universals]()

// globalOptions returns a list of global options for a command with the provided flag type
func globalOptions[F any]() (flags []meta.Flag) {
	flags = append(flags, universalOpts...)
	flags = append(flags, parser.AllFlags[F]()...)
	return
}

// globalFlagsFor returns a list of global options for a command with the provided flag type
func globalFlagsFor[F any](r Requirement[F]) (flags []meta.Flag) {
	// filter options to be those that are allowed
	gFlags := parser.AllFlags[F]()
	n := 0
	for _, flag := range gFlags {
		if !r.AllowsFlag(flag) {
			continue
		}
		gFlags[n] = flag
		n++
	}
	gFlags = gFlags[:n]

	// concat universal flags and normal flags
	flags = append(flags, universalOpts...)
	flags = append(flags, gFlags...)
	return
}
