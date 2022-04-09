package goprogram

import (
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/meta"
	"golang.org/x/exp/slices"
)

var errParseArgsNeedOneArgument = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to parse arguments: Need at least one argument. ",
}

var errParseArgsUnknownError = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to parse arguments: %s",
}

// parseP parses program-wide arguments.
//
// In particular, it *does not* parse command specific arguments.
// Any flags are just returned as unparsed positionals.
//
// When parsing fails, returns an error of type Error.
func (args *Arguments[F]) parseP(argv []string) error {
	var err error

	parser := meta.NewArgumentsParser(args)
	args.Pos, err = parser.ParseArgs(argv)

	// intercept unknonw flags
	if meta.IsUnknownFlag(err) {
		err = errParseArgsUnknownError.WithMessageF(err.Error())
	}

	// store the arguments we got and complain if there are none.
	// If we had a 'for' argument though, we should raise an error.
	if len(args.Pos) == 0 {
		switch {
		case args.Universals.Help || args.Universals.Version:
			return nil
		default:
			return errParseArgsNeedOneArgument
		}
	}

	// if we had help or version arguments we don't need to do
	// any more parsing and can bail out.
	if args.Universals.Help || args.Universals.Version {
		return nil
	}

	// setup command and arguments
	args.Command = args.Pos[0]
	args.Pos = args.Pos[1:]

	return err
}

// use prepares this context for using the provided command.
// It expects the context.Arguments object to exist, see the parseP method of Arguments.
//
// It expects that neither the Help nor Version flag of Arguments are true.
//
// When parsing fails, returns an error of type Error.
func (context *Context[E, P, F, R]) use(command Command[E, P, F, R]) error {
	context.Description = command.Description()
	context.parser = context.Description.ParserConfig.NewCommandParser(command)

	// specifically intercept the "--help" and "-h" arguments.
	// this prevents any kind of side effect from occuring.
	if slices.Contains(context.Args.Pos, "--help") || slices.Contains(context.Args.Pos, "-h") {
		context.Args.Universals.Help = true
		return nil
	}

	// check that the requirements for the command have been fullfilled
	if err := context.checkRequirements(); err != nil {
		return err
	}

	// do the actual parsing of the flags and validate that the right number of arguments has been given.
	if err := context.parseFlags(); err != nil {
		return err
	}

	if err := context.checkPositionalCount(); err != nil {
		return err
	}

	return nil
}

var errParseFlagSet = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Error parsing flags: %s",
}

// parseFlagset calls Parse() on the flagset.
// If the flagset has no defined flags (or is nil), immediatly returns nil
//
// When an error occurs, returns an error of type Error.
func (context *Context[E, P, F, R]) parseFlags() (err error) {
	context.Args.Pos, err = context.parser.ParseArgs(context.Args.Pos)

	// catch the help error
	if meta.IsHelp(err) {
		context.Args.Universals.Help = true
		err = nil
	}

	// if an error occured, return it!
	if err != nil {
		err = errParseFlagSet.WithMessageF(err.Error())
	}

	return err
}

func (context Context[E, P, F, R]) checkRequirements() error {
	return context.Description.Requirements.Validate(context.Args)
}

var errParseArgCount = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Wrong number of positional arguments for %s: %s",
}

func (context Context[E, P, F, R]) checkPositionalCount() error {
	err := context.Description.Positional.Validate(len(context.Args.Pos))
	if err != nil {
		return errParseArgCount.WithMessageF(context.Args.Command, err)
	}
	return nil
}
