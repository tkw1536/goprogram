// Package parser implements compatibility with the "github.com/jessevdk/go-flags" package
package parser

import (
	"errors"
	"reflect"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/goprogram/meta"
	"github.com/tkw1536/pkglib/reflectx"
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

// ParseArgs parses arguments for this parser.
//
// The returned error may be nil, a help error, an unknown flag error or otherwise.
// See also IsHelp, IsUnknownFlag.
func (p Parser) ParseArgs(args []string) ([]string, error) {
	// if we don't have a parser, parsing is a no-op!
	if p.parser == nil {
		return args, nil
	}

	// NOTE(twiesing): In a future version we probably want to wrap the error here
	// For now, the error is returned as-is with IsHelp() and IsUnknownFlag methods.
	return p.parser.ParseArgs(args)
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

// options collects all options contained in p or inside a group of p
func (p Parser) args() (options []*flags.Arg) {
	if p.parser == nil {
		return nil
	}

	return p.parser.Args()
}

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
