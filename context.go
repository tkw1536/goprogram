package goprogram

import (
	"context"

	"github.com/tkw1536/goprogram/parser"
	"github.com/tkw1536/pkglib/stream"
)

// Context represents an execution environment for a command.
// it takes the same type parameters as a command and program.
type Context[E any, P any, F any, R Requirement[F]] struct {
	stream.IOStream // IOStream describes the input and output the command reads from and writes to.
	Context         context.Context

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

	// inExec indicates if the current command is being called from within a program.Exec call.
	inExec bool

	// cleanup receives cleanup handlers
	cleanup chan ContextCleanupFunc[E, P, F, R]
}

type goProgramKey struct{}

var contextKey goProgramKey

// GetContext returns the goprogram context stored inside of ctx.
// When ctx has no context associated, returns nil.
func GetContext[E any, P any, F any, R Requirement[F]](ctx context.Context) *Context[E, P, F, R] {
	value := ctx.Value(contextKey)
	if value == nil {
		return nil
	}
	if vctx, ok := value.(*Context[E, P, F, R]); ok {
		return vctx
	}
	return nil
}

// withContext creates a new context that stores contextKey in parent.
func (ctx *Context[E, P, F, R]) withContext(parent context.Context) context.Context {
	return context.WithValue(parent, contextKey, ctx)
}

// ContextCleanupFunc represents a function that is called to cleanup a context.
// It is called with the context to be cleaned up.
//
// ContextCleanupFunc is guaranteed to be called even if the underlying operation associated with the context calls panic().
// This also holds if other ContextCleanupFuncs panic.
// There is no guarantee on the order in which functions are called.
//
// A nil ContextCleanupFunc is permitted.
type ContextCleanupFunc[E any, P any, F any, R Requirement[F]] func(context *Context[E, P, F, R])

// AddCleanupHandler adds f to be called when this context is no longer needed.
// f may be nil, in which case the call is ignored.
//
// Multiple handlers may be called in any order.
// f may not invoke AddCleanupHandler.
func (context Context[E, P, F, R]) AddCleanupFunction(f ContextCleanupFunc[E, P, F, R]) {
	context.cleanup <- f
}

// handleCleanup initializes a new cleanup channel and begins handling cleanup functions.
func (context *Context[E, P, F, R]) handleCleanup() func() {
	context.cleanup = make(chan ContextCleanupFunc[E, P, F, R])
	done := make(chan struct{})
	go func() {
		defer close(done)

		// collect all the cleanup functions
		var cleanups []ContextCleanupFunc[E, P, F, R]
		for handler := range context.cleanup {
			if handler == nil {
				continue
			}
			cleanups = append(cleanups, handler)
		}

		// defer them in the correct order!
		for i := len(cleanups) - 1; i >= 0; i-- {
			defer cleanups[i](context)
		}

	}()
	return func() {
		close(context.cleanup)
		<-done
	}
}

// InExec returns true if this execution environment was started from within a program.Exec command.
func (context Context[E, P, F, R]) InExec() bool {
	return context.inExec
}

// Exec is like context.Program.Exec.
func (context Context[E, P, F, R]) Exec(command string, args ...string) error {
	return context.Program.Exec(context, command, args...)
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
