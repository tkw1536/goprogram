// Package goprogram provides a program abstraction that can be used to create programs.
//
//spellchecker:words goprogram
package goprogram

//spellchecker:words context github goprogram exit meta pkglib stream
import (
	"context"
	"fmt"

	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/goprogram/meta"
	"go.tkw01536.de/pkglib/stream"
)

//spellchecker:words nolint wrapcheck

// Program represents an executable program.
// A program is intended to be invoked on the command line.
// Each invocation of a program executes a command.
//
// Programs have 4 type parameters:
// An environment of type E, a type of parameters P, a type of flags F and a type requirements R.
//
// The Environment type E defines a runtime environment for commands to execute in.
// An Environment is created using the NewEnvironment function, taking parameters P.
//
// The type of (global) command line flags F is backed by a struct type.
// It is jointed by a type of Requirements R which impose restrictions on flags for commands.
//
// Internally a program also contains a list of commands, keywords and aliases.
//
// See the Main method for a description of how program execution takes place.
type Program[E any, P any, F any, R Requirement[F]] struct {
	// Meta-information about the current program
	// Used to generate help and version pages
	Info meta.Info

	// The NewContext function is called to create a new context for a command.
	// It may optionally return a ContextCleanupFunc.
	//
	// It is passed the parameters as well as a parent context.
	// The parameters are only passed when Main() is invoked; otherwise they are the zero value.
	//
	// The parent context is either a new context, or the context from the parent command this command was invoked from.
	//
	// To access the full context object within the NewContext function, use GetContext.
	//
	// When NewContext is nil, the parent context is used.
	//
	// If the context is closed before a command would be invoked, then the command is not invoked.
	//
	// See also [env.NewOSContext].
	NewContext func(params *P, parent context.Context) (context.Context, ContextCleanupFunc[E, P, F, R], error)

	// The NewEnvironment function is used to create a new environment.
	// The returned error must be nil or of type exit.Error.
	//
	// NewEnvironment may be nil, in which case a new environment is assumed to be
	// the zero value of type E.
	NewEnvironment func(params P, context Context[E, P, F, R]) (E, error)

	// BeforeKeyword, BeforeAlias and BeforeCommand (if non-nil) are invoked right before their respective datum is executed.
	// They are intended to act as a guard before executing a particular datum.
	//
	// The returned error must be nil or of type exit.Error.
	// When non-nil, the error is returned to the caller of Main().
	BeforeKeyword func(context Context[E, P, F, R], keyword Keyword[F]) error
	BeforeAlias   func(context Context[E, P, F, R], alias Alias) error
	BeforeCommand func(context Context[E, P, F, R], command Command[E, P, F, R]) error

	// Commands, Keywords, and Aliases associated with this program.
	// They are expanded in order; see Main for details.
	keywords map[string]Keyword[F]
	aliases  map[string]Alias
	commands map[string]Command[E, P, F, R]
}

// initContext initializes the context of the context.
func (p Program[E, P, F, R]) initContextContext(params *P, context *Context[E, P, F, R]) error {
	context.Context = context.withContext(context.Context)

	if p.NewContext == nil {
		return nil
	}

	// make a new context
	ctx, cleanup, err := p.NewContext(params, context.withContext(context.Context))
	if err != nil {
		return err
	}

	// store it and add the cleanup function
	context.Context = ctx
	context.AddCleanupFunction(cleanup)

	return nil
}

// Main invokes this program and returns an error of type exit.Error or nil.
//
// Main takes input / output streams, parameters for the environment and a set of command-line arguments.
//
// It first parses these into arguments for a specific command to be executed.
// Next, it executes any keywords and expands any aliases.
// Finally, it executes the requested command or displays a help or version page.
//
// For keyword actions, see Keyword.
// For alias expansion, see Alias.
// For command execution, see Command.
//
// For help pages, see MainUsage, CommandUsage, AliasUsage.
// For version pages, see FmtVersion.
func (p Program[E, P, F, R]) Main(str stream.IOStream, params P, argv []string) (err error) {
	// whenever an error occurs, we want it printed
	defer func() {
		err = exit.Die(str, err)
	}()

	// create a new context
	context := Context[E, P, F, R]{
		Context:  context.Background(),
		IOStream: str,
		Program:  p,
	}
	defer context.handleCleanup()()

	// parse flags!
	if err := context.Args.parseProgramFlags(argv); err != nil {
		return err
	}

	// initialize the underlying context
	if err := p.initContextContext(&params, &context); err != nil {
		return err
	}

	// and run!
	return p.run(context, func(context Context[E, P, F, R]) (E, error) {
		e, err := p.makeEnvironment(params, context)
		return e, err
	})
}

// Exec executes this program from within a given context.
//
// It does not create a new environment.
// It does not re-parse arguments preceding the keyword, alias or command.
//
// This function is intended to safely run a command from within another command.
func (p Program[E, P, F, R]) Exec(context Context[E, P, F, R], command string, pos ...string) (err error) {
	// NOTE: This function is untested, because it is nearly identical to Main

	// create a new context
	execContext := Context[E, P, F, R]{
		IOStream: context.IOStream,
		Context:  context.Context,
		Program:  p,

		Args: Arguments[F]{
			Universals: context.Args.Universals,
			Flags:      context.Args.Flags,

			Command: command,
			pos:     pos,
		},

		inExec: true,
	}
	defer context.handleCleanup()()

	// initialize the underlying context
	if err := p.initContextContext(nil, &context); err != nil {
		return err
	}

	// reset the arguments to the context
	return p.run(execContext, func(Context[E, P, F, R]) (E, error) { return context.Environment, nil })
}

var (
	errProgramUnknownCommand = exit.NewErrorWithCode("unknown command", exit.ExitUnknownCommand) //  must be one of %s
	errProgramContext        = exit.NewErrorWithCode("context was closed before main could run", exit.ExitContext)
	errProgramIO             = exit.NewErrorWithCode("failed to write to context", exit.ExitContext)
)

// run implements Main and Exec.
//
//nolint:wrapcheck
func (p Program[E, P, F, R]) run(context Context[E, P, F, R], setupEnvironment func(context Context[E, P, F, R]) (E, error)) (err error) {
	// expand keyword
	keyword, hasKeyword := p.keywords[context.Args.Command]
	if hasKeyword {
		// invoke BeforeKeyword (if any)
		if p.BeforeKeyword != nil {
			err := p.BeforeKeyword(context, keyword)
			if err != nil {
				return err
			}
		}
		if err := keyword(&context.Args, &context.Args.pos); err != nil {
			return err
		}
	}

	// handle universals
	switch {
	case context.Args.Universals.Help:
		_, err = context.Println(p.MainUsage().String())
		if err != nil {
			return fmt.Errorf("%w: %w", errProgramIO, err)
		}
		return nil
	case context.Args.Universals.Version:
		_, err = context.Println(p.Info.FmtVersion())
		if err != nil {
			return fmt.Errorf("%w: %w", errProgramIO, err)
		}
		return nil
	}

	// expand the alias (if any)
	alias, hasAlias := p.aliases[context.Args.Command]
	if hasAlias {
		// invoke BeforeAlias (if any)
		if p.BeforeAlias != nil {
			err := p.BeforeAlias(context, alias)
			if err != nil {
				return err
			}
		}
		context.Args.Command, context.Args.pos = alias.Invoke(context.Args.pos)
	}

	// load the command if we have it
	command, hasCommand := p.Command(context.Args.Command)
	if !hasCommand {
		return fmt.Errorf("%w: must be one of %s", errProgramUnknownCommand, p.FmtCommands())
	}

	// make the context use the given command
	if err := context.use(command); err != nil {
		return err
	}

	// write out help information (if given)
	if context.Args.Universals.Help {
		if hasAlias {
			_, err = context.Println(p.AliasUsage(context, alias).String())
			if err != nil {
				return fmt.Errorf("%w: %w", errProgramIO, err)
			}
			return nil
		}
		_, err = context.Println(p.CommandUsage(context).String())
		if err != nil {
			return fmt.Errorf("%w: %w", errProgramIO, err)
		}
		return nil
	}

	// call the AfterParse hook
	if ap, isAP := command.(AfterParseCommand[E, P, F, R]); isAP {
		if err := ap.AfterParse(); err != nil {
			return err
		}
	}

	// create the environment
	if context.Environment, err = setupEnvironment(context); err != nil {
		return err
	}

	// invoke BeforeCommand (if any)
	if p.BeforeCommand != nil {
		err := p.BeforeCommand(context, command)
		if err != nil {
			return err
		}
	}

	// check that the context isn't closed!
	if err := context.Context.Err(); err != nil {
		return fmt.Errorf("%w: %w", errProgramContext, err)
	}

	// do the command!
	return command.Run(context)
}

// makeEnvironment creates a new environment for the given command.
func (p Program[E, P, F, R]) makeEnvironment(params P, context Context[E, P, F, R]) (E, error) {
	if p.NewEnvironment == nil {
		var zeroE E
		return zeroE, nil
	}

	return p.NewEnvironment(params, context)
}
