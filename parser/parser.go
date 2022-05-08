// Package parser implements compatiblity with the "github.com/jessevdk/go-flags" package
package parser

import (
	"errors"
	"os"
	"reflect"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/goprogram/lib/reflectx"
	"github.com/tkw1536/goprogram/meta"
)

// Parse represents a parser for arguments.
//
// It is internally backed by the "github.com/jessevdk/go-flags" package.
type Parser struct {
	// NOTE(twiesing): This entire struct is not directly tested
	// Instead the tests are performed using higher-level integration tests
	parser *flags.Parser
	tp     reflect.Type
}

const goFlagsCompletion = "GO_FLAGS_COMPLETION"

// Parse parses arguments for this parser.
// It guarantees that the completion handler is not called.
//
// The returned error may be nil, a help error, an unknown flag error or otherwise.
// See also IsHelp, IsUnknownFlag.
func (p Parser) Parse(args []string) ([]string, error) {
	// TESTME

	// if we don't have a parser, parsing is a no-op!
	if p.parser == nil {
		return args, nil
	}

	// store the completion variable
	defer os.Setenv(goFlagsCompletion, os.Getenv(goFlagsCompletion))

	// NOTE(twiesing): In a future version we probably want to wrap the error here
	// For now, the error is returned as-is with IsHelp() and IsUnknownFlag methods.
	os.Setenv(goFlagsCompletion, "")
	return p.parser.ParseArgs(args)
}

// Completion represents the completion of a single flag of the parser.
type Completion = flags.Completion

// Complete completes a set of partial arguments provided to this parser.
// It guarantees that argument parsing does not take place, and no state is modified.
//
// The returned error may be nil
func (p Parser) Complete(args []string) (items []Completion, err error) {
	// TESTME

	// if we don't have a parser, there is nothing to complete
	if p.parser == nil {
		return nil, nil
	}

	// store the old completion handler
	defer func(handler func(items []flags.Completion)) {
		p.parser.CompletionHandler = handler
	}(
		p.parser.CompletionHandler,
	)

	// store the completion env
	defer os.Setenv(goFlagsCompletion, os.Getenv(goFlagsCompletion))

	// setup a completion handler
	var ok bool
	p.parser.CompletionHandler = func(i []flags.Completion) {
		ok = true
		items = i
	}

	// do the completion!
	os.Setenv(goFlagsCompletion, "1")
	p.parser.ParseArgs(args)

	// if no one called the completion handler, then "go-flags" changed it's implementation
	// and we shouldn't assume anything anymore!
	if !ok {
		panic("Parser.Complete(): CompletionHandler was not called")
	}

	// and return them!
	return items, nil
}

// IsHelp checks if err represents the help flag being passed
func IsHelp(err error) bool {
	var flagError *flags.Error
	return errors.As(err, &flagError) && flagError.Type == flags.ErrHelp
}

// IsUnknownFlag checks if err indicates an unknown flag
func IsUnknownFlag(err error) bool {
	var flagError *flags.Error
	return errors.As(err, &flagError) && flagError.Type == flags.ErrUnknownFlag
}

// Positionals returns information about all the positionals found in this parser
func (p Parser) Positionals() []meta.Positional {
	// collect the args
	args := p.args()
	types := p.argTypes()
	if len(args) != len(types) {
		panic("Parser.Positionals(): len(args) != len(types)")
	}

	// turn them into proper positionals
	poss := make([]meta.Positional, len(args))
	for i, arg := range args {
		poss[i] = NewPositional(arg, types[i])
	}
	return poss
}

// args collects all the arguments found in this parser
func (p Parser) args() (args []*flags.Arg) {
	if p.parser == nil {
		return nil
	}

	return p.parser.Args()
}

// argTypes collects the struct fields corresponding to the arguments in this parser
func (p Parser) argTypes() (types []reflect.StructField) {
	if p.parser == nil {
		return nil
	}

	reflectx.IterateAllFields(p.tp, func(field reflect.StructField, index ...int) (cancel bool) {
		// check that we actually have a "positional-args" field
		if field.Tag.Get("positional-args") == "" || field.Type.Kind() != reflect.Struct {
			return
		}

		// iterate over all the fields in the nested struct
		nf := field.Type.NumField()
		for j := 0; j < nf; j++ {
			types = append(types, field.Type.Field(j))
		}

		return false
	})

	return
}

// Flags returns information about the flags belonging to this parser
func (p Parser) Flags() []meta.Flag {
	// collect the options
	options := p.options()

	// turn them into proper flags
	flags := make([]meta.Flag, len(options))
	for i, opt := range options {
		flags[i] = NewFlag(opt)
	}
	return flags
}

// options collects all options contained in p or inside a group of p
func (p Parser) options() (options []*flags.Option) {
	if p.parser == nil {
		return nil
	}

	groups := p.parser.Groups()
	for _, g := range groups {
		options = append(options, g.Options()...)
	}

	return
}
