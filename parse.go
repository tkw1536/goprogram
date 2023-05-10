package goprogram

import (
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
	"golang.org/x/exp/slices"
)

var errParseArgsNeedOneArgument = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "unable to parse arguments: need at least one argument",
}

var errParseArgsUnknownError = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "unable to parse arguments: %s",
}

// parseProgramFlags parses program-wide arguments.
//
// In particular, it *does not* parse command specific arguments.
// Any flags are just returned as unparsed positionals.
//
// When parsing fails, returns an error of type Error.
func (args *Arguments[F]) parseProgramFlags(argv []string) error {
	var err error

	argsParser := parser.NewArgumentsParser(args)
	args.pos, err = argsParser.ParseArgs(argv)

	// intercept unknown flags
	if parser.IsUnknownFlag(err) {
		err = errParseArgsUnknownError.WithMessageF(err.Error())
	}

	// store the arguments we got and complain if there are none.
	// If we had a 'for' argument though, we should raise an error.
	if len(args.pos) == 0 {
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
	args.Command = args.pos[0]
	args.pos = args.pos[1:]

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
	// this prevents any kind of side effect from occurring.
	if slices.Contains(context.Args.pos, "--help") || slices.Contains(context.Args.pos, "-h") {
		context.Args.Universals.Help = true
		return nil
	}

	// check that the requirements for the command are fulfilled
	if err := context.Description.Requirements.Validate(context.Args); err != nil {
		return err
	}

	// parse the command flags
	if err := context.parseCommandFlags(); err != nil {
		return err
	}

	// check that no positional arguments are left over
	if len(context.Args.pos) > 0 {
		return errParseArgCount.WithMessageF(context.Args.Command, len(context.Args.pos))
	}

	return nil
}

var errWrongArguments = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "wrong arguments for %s: %s",
}

// parseCommandFlags uses the parser to parse flags passed directly to the command.
//
// When an error occurs, returns an error of type Error.
func (context *Context[E, P, F, R]) parseCommandFlags() (err error) {
	context.Args.pos, err = context.parser.ParseArgs(context.Args.pos)

	// catch the help error
	if parser.IsHelp(err) {
		context.Args.Universals.Help = true
		err = nil
	}

	// if an error occurred, return it!
	if err != nil {
		err = errWrongArguments.WithMessageF(context.Args.Command, err.Error())
	}

	return err
}

var errParseArgCount = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "wrong number of positional arguments for %s: %d additional arguments were provided",
}
