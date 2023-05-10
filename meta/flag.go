package meta

import (
	"io"

	"github.com/tkw1536/pkglib/docfmt"
	"github.com/tkw1536/pkglib/text"
	"golang.org/x/exp/slices"
)

// Flag holds meta-information about a single flag of a command.
//
// To create a new flag, see parser.NewFlag.
type Flag struct {

	// For the purposes of documentation we use the following argument as an example.
	//   -n, --number digit  A digit used within something (default: 42)

	// The name of the underlying struct field this flag comes from.
	FieldName string // "Number"

	// Short and Long Names of the flag
	// each potentially more than one
	Short []string // ["n"]
	Long  []string // ["number"]

	// Indicates if the flag is required
	Required bool // false

	// Name and Description of the flag in help texts
	Value string // "digit"
	Usage string // "A digit used within something"

	// Default value of the flag, as shown to the user.
	// When multiple default values are set, they are joined as a string.
	Default string // "42"

	// Valid choices for the option
	Choices []string
}

// WriteSpecTo writes a short specification of f into w.
// It is of the form
//
//	--flag|-f value
//
// WriteSpecTo adds braces around the argument if it is optional.
func (f Flag) WriteSpecTo(w io.Writer) {
	f.spec(w, "|", true, true)
}

// WriteLongSpecTo writes a long specification of f into w.
// It is of the form
//
//	-f, --flag value
//
// WriteLongSpecTo does not add any brackets around the argument.
func (opt Flag) WriteLongSpecTo(w io.Writer) {
	opt.spec(w, ", ", false, false)
}

// spec implements SpecShort and SpecLong.
//
// sep indicates how to separate arguments.
// longFirst indicates that long argument names should be listed before short arguments.
// optionalBraces indicates if braces should be placed around the argument if it is optional.
func (opt Flag) spec(w io.Writer, sep string, longFirst bool, optionalBraces bool) {
	// if the argument is optional put braces around it!
	if optionalBraces && !opt.Required {
		io.WriteString(w, "[")
		defer io.WriteString(w, "]")
	}

	// collect long and short arguments and combine them
	la := slices.Clone(opt.Long)
	for k, v := range la {
		la[k] = "--" + v
	}

	sa := slices.Clone(opt.Short)
	for k, v := range sa {
		sa[k] = "-" + v
	}

	// write the joined versions of the arguments into the specification
	var args []string
	if longFirst {
		args = append(la, sa...)
	} else {
		args = append(sa, la...)
	}
	text.Join(w, args, sep)

	// write the value (if any)
	if value := opt.Value; value != "" {
		io.WriteString(w, " ")
		io.WriteString(w, value)
	}
}

// usageMsgTpl is the template for long usage messages
// it is split into three parts, that are joined by the arguments.
//
//	const usageMsgTpl = usageMsg1 + "%s" + usageMsg2 + "%s" + usageMsg3
const (
	usageMsg1 = "\n\n   "
	usageMsg2 = "\n      "
	usageMsg3 = ""
)

// WriteMessageTo writes a long message of f to w.
// It is of the form
//
//	-f, --flag ARG
//
// and
//
//	DESCRIPTION (choices CHOICE1, CHOICE2. default DEFAULT)
//
// .
//
// This function is implicitly tested via other tests.
func (opt Flag) WriteMessageTo(w io.Writer) {

	io.WriteString(w, usageMsg1)
	opt.WriteLongSpecTo(w)
	io.WriteString(w, usageMsg2)

	io.WriteString(w, docfmt.Format(opt.Usage))

	{
		dflt := opt.Default
		hasDefault := dflt != ""
		choices := opt.Choices
		hasChoices := len(choices) > 0

		if hasDefault || hasChoices {
			io.WriteString(w, " (")
			if hasChoices {
				io.WriteString(w, "choices: ")
				text.Join(w, opt.Choices, ", ")
				if hasDefault {
					io.WriteString(w, "; ")
				}
			}

			if hasDefault {
				io.WriteString(w, "default ")
				io.WriteString(w, dflt)
			}

			io.WriteString(w, ")")
		}
	}

	io.WriteString(w, usageMsg3)
}
