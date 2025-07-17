//spellchecker:words parser
package parser

//spellchecker:words reflect github jessevdk flags
import (
	"reflect"

	"github.com/jessevdk/go-flags"
)

// NewCommandParser checks if command represents a valid command and, when this is the case, creates a new parser for it
// with the config provided in ParserConfig.
func NewCommandParser(command any) (p Parser) {
	// the command must be backed by a pointed-to struct
	// when this is not the case, we don't need to create a parser
	if ptrVal := reflect.TypeOf(command); command == nil || ptrVal.Kind() != reflect.Ptr || ptrVal.Elem().Kind() != reflect.Struct {
		return
	}

	// and make the parser
	p.parser = flags.NewParser(command, flags.PassDoubleDash|flags.HelpFlag)
	p.tp = reflect.TypeOf(command).Elem()
	return
}

// NewArgumentsParser creates a new parser fo parse a set of arguments.
func NewArgumentsParser(args any) Parser {
	return Parser{
		parser: flags.NewParser(args, flags.PassAfterNonOption|flags.PassDoubleDash),
	}
}
