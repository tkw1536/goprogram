//spellchecker:words meta
package meta

//spellchecker:words slices github pkglib docfmt text
import (
	"io"
	"slices"

	"github.com/tkw1536/pkglib/docfmt"
	"github.com/tkw1536/pkglib/text"
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
func (f Flag) WriteSpecTo(w io.Writer) error {
	return f.spec(w, "|", true, true)
}

// WriteLongSpecTo writes a long specification of f into w.
// It is of the form
//
//	-f, --flag value
//
// WriteLongSpecTo does not add any brackets around the argument.
func (opt Flag) WriteLongSpecTo(w io.Writer) error {
	return opt.spec(w, ", ", false, false)
}

// spec implements SpecShort and SpecLong.
//
// sep indicates how to separate arguments.
// longFirst indicates that long argument names should be listed before short arguments.
// optionalBraces indicates if braces should be placed around the argument if it is optional.
func (opt Flag) spec(w io.Writer, sep string, longFirst bool, optionalBraces bool) (err error) {
	// if the argument is optional put braces around it!
	if optionalBraces && !opt.Required {
		if _, err := io.WriteString(w, "["); err != nil {
			return err
		}
		defer func() {
			if err != nil {
				return
			}
			_, err = io.WriteString(w, "]")
		}()
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
	if _, err := text.Join(w, args, sep); err != nil {
		return err
	}

	// write the value (if any)
	if value := opt.Value; value != "" {
		if _, err := io.WriteString(w, " "); err != nil {
			return err
		}
		if _, err := io.WriteString(w, value); err != nil {
			return err
		}
	}

	return nil
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
func (opt Flag) WriteMessageTo(w io.Writer) error {

	if _, err := io.WriteString(w, usageMsg1); err != nil {
		return err
	}
	if err := opt.WriteLongSpecTo(w); err != nil {
		return err
	}
	if _, err := io.WriteString(w, usageMsg2); err != nil {
		return err
	}

	if _, err := io.WriteString(w, docfmt.Format(opt.Usage)); err != nil {
		return err
	}

	{
		Default := opt.Default
		hasDefault := Default != ""
		choices := opt.Choices
		hasChoices := len(choices) > 0

		if hasDefault || hasChoices {
			if _, err := io.WriteString(w, " ("); err != nil {
				return err
			}
			if hasChoices {
				if _, err := io.WriteString(w, "choices: "); err != nil {
					return err
				}
				if _, err := text.Join(w, opt.Choices, ", "); err != nil {
					return err
				}
				if hasDefault {
					if _, err := io.WriteString(w, "; "); err != nil {
						return err
					}
				}
			}

			if hasDefault {
				if _, err := io.WriteString(w, "default "); err != nil {
					return err
				}
				if _, err := io.WriteString(w, Default); err != nil {
					return err
				}
			}

			if _, err := io.WriteString(w, ")"); err != nil {
				return err
			}
		}
	}

	if _, err := io.WriteString(w, usageMsg3); err != nil {
		return err
	}

	return nil
}
