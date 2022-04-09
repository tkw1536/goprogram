package meta

import (
	"errors"
	"reflect"

	"github.com/jessevdk/go-flags"
)

// Parse represents a parser for arguments.
//
// It is internally backed by the "github.com/jessevdk/go-flags" package.
type Parser struct {
	parser *flags.Parser
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

// Flags returns information about the flags belonging to this parser
func (p Parser) Flags() []Flag {
	// no parser => no options
	if p.parser == nil {
		return nil
	}

	// collect all the options
	var options []*flags.Option
	groups := p.parser.Groups()
	for _, g := range groups {
		options = append(options, g.Options()...)
	}

	// turn them into proper flags
	flags := make([]Flag, len(options))
	for i, opt := range options {
		flags[i] = NewFlag(opt)
	}
	return flags
}

// ParserConfig represents configuration for the parser for a command
type ParserConfig struct {
	// IncludeUnknown causes unknown flags to be parsed as positional arguments.
	// When IncludeUnknown in false, unknown flags produce an error instead.
	IncludeUnknown bool
}

// NewCommandParser checks if command represents a valid command and, when this is the case, creates a new parser for it
// with the config provided in ParserConfig.
func (cfg ParserConfig) NewCommandParser(command any) (p Parser) {

	// the command must be backed by a pointed-to struct
	// when this is not the case, we don't need to create a parser
	if ptrval := reflect.TypeOf(command); command == nil || ptrval.Kind() != reflect.Ptr || ptrval.Elem().Kind() != reflect.Struct {
		return
	}

	// create options for the parser
	var options flags.Options = flags.PassDoubleDash | flags.HelpFlag
	if cfg.IncludeUnknown {
		options |= flags.IgnoreUnknown
	}

	// and make it
	p.parser = flags.NewParser(command, options)
	return
}

// NewArgumentsParser creates a new parser fo parse a set of arguments
func NewArgumentsParser(args any) Parser {
	return Parser{
		parser: flags.NewParser(args, flags.PassAfterNonOption|flags.PassDoubleDash),
	}
}
