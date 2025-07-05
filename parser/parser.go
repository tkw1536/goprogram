// Package parser implements compatibility with the "github.com/jessevdk/go-flags" package
//
//spellchecker:words parser
package parser

//spellchecker:words errors reflect github jessevdk flags goprogram meta pkglib reflectx
import (
	"errors"
	"reflect"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/goprogram/meta"
	"go.tkw01536.de/pkglib/reflectx"
)

//spellchecker:words positionals nolint wrapcheck

// Parse represents a parser for arguments.
//
// It is internally backed by the "github.com/jessevdk/go-flags" package.
type Parser struct {
	// NOTE: This entire struct is not directly tested
	// Instead the tests are performed using higher-level integration tests
	parser *flags.Parser
	tp     reflect.Type
}

// ParseArgs parses arguments for this parser.
//
// The returned error may be nil, a help error, an unknown flag error or otherwise.
// See also IsHelp, IsUnknownFlag.
//
//nolint:wrapcheck
func (p Parser) ParseArgs(args []string) ([]string, error) {
	// if we don't have a parser, parsing is a no-op!
	if p.parser == nil {
		return args, nil
	}

	// NOTE: In a future version we probably want to wrap the error here
	// For now, the error is returned as-is with IsHelp() and IsUnknownFlag methods.
	return p.parser.ParseArgs(args)
}

// IsHelp checks if err represents the help flag being passed.
func IsHelp(err error) bool {
	var flagError *flags.Error
	return errors.As(err, &flagError) && flagError.Type == flags.ErrHelp
}

// IsUnknownFlag checks if err indicates an unknown flag.
func IsUnknownFlag(err error) bool {
	var flagError *flags.Error
	return errors.As(err, &flagError) && flagError.Type == flags.ErrUnknownFlag
}

// options collects all options contained in p or inside a group of p.
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

	for field := range reflectx.IterAllFields(p.tp) {
		// check that we actually have a "positional-args" field
		if field.Tag.Get("positional-args") == "" || field.Type.Kind() != reflect.Struct {
			continue
		}

		// iterate over all the fields in the nested struct
		nf := field.Type.NumField()
		for j := range nf {
			types = append(types, field.Type.Field(j))
		}

		break
	}

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

// options collects all options contained in p or inside a group of p.
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

// Flags returns information about the flags belonging to this parser.
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
