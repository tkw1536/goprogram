//spellchecker:words parser
package parser

//spellchecker:words reflect strings github jessevdk flags goprogram meta
import (
	"reflect"
	"strings"

	"github.com/jessevdk/go-flags"
	"go.tkw01536.de/goprogram/meta"
)

//spellchecker:words positionals

// NewFlag creates a new flag based on an option from the flags package.
func NewFlag(option *flags.Option) (flag meta.Flag) {
	flag.Required = option.Required

	short := option.ShortName
	if short != rune(0) {
		flag.Short = []string{string(short)}
	}

	long := option.LongName
	if long != "" {
		flag.Long = []string{long}
	}

	flag.FieldName = option.Field().Name

	flag.Value = option.ValueName

	flag.Usage = option.Description

	Default := option.Default
	if len(Default) != 0 {
		flag.Default = strings.Join(Default, ", ")
	}
	if option.DefaultMask != "" {
		flag.Default = option.DefaultMask
	}

	flag.Choices = option.Choices

	return
}

// AllFlags is a convenience method to get all flags for the provided argument type.
func AllFlags[T any]() []meta.Flag {
	return AllFlagsOf(new(T))
}

// AllFlagsOf is a convenience method to get all flags of the provided argument.
func AllFlagsOf(data any) []meta.Flag {
	return Parser{
		parser: flags.NewParser(data, flags.None),
		tp:     reflect.TypeOf(data).Elem(),
	}.Flags()
}

// NewPositional creates a new Positional from a flag argument.
func NewPositional(arg *flags.Arg, field reflect.StructField) (pos meta.Positional) {
	pos.Value = arg.Name
	pos.Usage = arg.Description

	pos.Min, pos.Max = arg.Required, arg.RequiredMaximum
	if pos.Min == -1 {
		pos.Min = 0
	}
	if pos.Max == -1 && field.Type.Kind() != reflect.Slice {
		pos.Max = 1
	}
	return
}

func AllPositionals[T any]() []meta.Positional {
	data := new(T)
	return Parser{
		parser: flags.NewParser(data, flags.None),
		tp:     reflect.TypeOf(data).Elem(),
	}.Positionals()
}
