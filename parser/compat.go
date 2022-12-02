package parser

import (
	"reflect"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/goprogram/meta"
)

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

	dflt := option.Default
	if len(dflt) != 0 {
		flag.Default = strings.Join(dflt, ", ")
	}
	if option.DefaultMask != "" {
		flag.Default = option.DefaultMask
	}

	flag.Choices = option.Choices

	return
}

// AllFlags is a convenience method to get all flags for the provided argument type
func AllFlags[T any]() []meta.Flag {
	data := new(T)
	return Parser{
		parser: flags.NewParser(data, flags.None),
		tp:     reflect.TypeOf(data).Elem(),
	}.Flags()
}

// NewPositional creates a new Positional from a flag argument
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
