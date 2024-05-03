//spellchecker:words parser
package parser

//spellchecker:words reflect github jessevdk flags
import (
	"reflect"

	"github.com/jessevdk/go-flags"
)

// Config represents configuration for the parser for a command
type Config struct {
	// IncludeUnknown causes unknown flags to be parsed as positional arguments.
	// When IncludeUnknown in false, unknown flags produce an error instead.
	IncludeUnknown bool
}

// NewCommandParser checks if command represents a valid command and, when this is the case, creates a new parser for it
// with the config provided in ParserConfig.
func (cfg Config) NewCommandParser(command any) (p Parser) {

	// the command must be backed by a pointed-to struct
	// when this is not the case, we don't need to create a parser
	if ptrVal := reflect.TypeOf(command); command == nil || ptrVal.Kind() != reflect.Ptr || ptrVal.Elem().Kind() != reflect.Struct {
		return
	}

	// create options for the parser
	var options flags.Options = flags.PassDoubleDash | flags.HelpFlag
	if cfg.IncludeUnknown {
		options |= flags.IgnoreUnknown
	}

	// and make it
	p.parser = flags.NewParser(command, options)
	p.tp = reflect.TypeOf(command).Elem()
	return
}

// NewArgumentsParser creates a new parser fo parse a set of arguments
func NewArgumentsParser(args any) Parser {
	return Parser{
		parser: flags.NewParser(args, flags.PassAfterNonOption|flags.PassDoubleDash),
	}
}
