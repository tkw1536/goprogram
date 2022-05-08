package goprogram

import (
	"github.com/tkw1536/goprogram/parser"
	"github.com/tkw1536/goprogram/stream"
)

// Context represents an execution environment for a command.
// it takes the same type parameters as a command and program.
type Context[E any, P any, F any, R Requirement[F]] struct {
	// IOStream describes the input and output the command reads from and writes to.
	stream.IOStream

	// Args holds arguments passed to this command.
	Args Arguments[F]

	// Description is the description of the command being invoked
	Description Description[F, R]
	// Environment holds the environment for this command.
	Environment E

	// Program is the current program this command is executing under
	Program Program[E, P, F, R]

	// parser holds a parser for command-specific arguments
	// this refers to the command itself
	parser parser.Parser
}

// expandKeywords expands any keywords in the context (if any)
func (context *Context[E, P, F, R]) expandKeywords() (hasKeyword bool, err error) {
	var keyword Keyword[F]
	keyword, hasKeyword = context.Program.keywords[context.Args.Command]
	if hasKeyword {
		err = keyword(&context.Args, &context.Args.pos)
	}
	return
}

// expandAliases expands any alias in the context (if any)
func (context *Context[E, P, F, R]) expandAliases() (alias Alias, hasAlias bool) {
	alias, hasAlias = context.Program.aliases[context.Args.Command]
	if hasAlias {
		context.Args.Command, context.Args.pos = alias.Invoke(context.Args.pos)
	}
	return
}

// Arguments represent a set of command-independent arguments passed to a command.
// These should be further parsed into CommandArguments using the appropriate Parse() method.
//
// Command line argument are annotated using syntax provided by "github.com/jessevdk/go-flags".
type Arguments[F any] struct {
	Universals Universals `group:"universals"`
	Flags      F          `group:"flags"`

	Command string   // command to run
	pos     []string // positional arguments
}

// Universals holds flags added to every executable.
//
// Command line arguments are annotated using syntax provided by "github.com/jessevdk/go-flags".
type Universals struct {
	Help    bool `short:"h" long:"help" description:"print a help message and exit"`
	Version bool `short:"v" long:"version" description:"print a version message and exit"`
}
