//spellchecker:words meta
package meta

//spellchecker:words errors github pkglib text
import (
	"errors"
	"fmt"
	"io"

	"github.com/tkw1536/pkglib/text"
)

// Positional holds meta-information about a positional argument.
type Positional struct {
	// Name and Description of the positional in help texts
	Value string // defaults to "ARGUMENT"
	Usage string

	// Min and Max indicate how many positional arguments are expected for this command.
	// Min must be >= 0. Max must be either Min, or -1.
	// Max == -1 indicates an unlimited number of repeats.
	Min, Max int
}

// ValidRange checks if positional has valid min and max values.
func (pos Positional) ValidRange() bool {
	extra := pos.Max - pos.Min
	return pos.Min >= 0 && (pos.Max <= 0 || extra >= 0)
}

// defaultPositionalValue is the default name used for a positional argument.
const defaultPositionalValue = "ARGUMENT"

var errPositionalInvalidRange = errors.New("positional: invalid range")

// WriteSpecTo writes a specification of this argument into w.
// A specification looks like "arg [arg...]".
func (pos Positional) WriteSpecTo(w io.Writer) error {
	extra := pos.Max - pos.Min

	if !pos.ValidRange() {
		return errPositionalInvalidRange
	}

	if pos.Value == "" {
		pos.Value = defaultPositionalValue
	}

	// nothing to generate!
	if pos.Max == 0 && extra == 0 {
		return nil
	}

	// arg arg arg
	if _, err := text.RepeatJoin(w, pos.Value, " ", pos.Min); err != nil {
		return fmt.Errorf("unable to repeat and join value: %w", err)
	}
	if pos.Min > 0 && extra != 0 {
		if _, err := io.WriteString(w, " "); err != nil {
			return fmt.Errorf("unable to write space: %w", err)
		}
	}

	if pos.Max < 0 {
		// [arg ...]
		if _, err := io.WriteString(w, "["); err != nil {
			return fmt.Errorf("unable to write '[': %w", err)
		}
		if _, err := io.WriteString(w, pos.Value); err != nil {
			return fmt.Errorf("unable to write value: %w", err)
		}
		if _, err := io.WriteString(w, " ...]"); err != nil {
			return fmt.Errorf("unable to write '...]': %w", err)
		}
		return nil
	}

	// [arg [arg]]
	if _, err := text.RepeatJoin(w, "["+pos.Value, " ", extra); err != nil {
		return fmt.Errorf("unable to repeat and join args: %w", err)
	}
	if _, err := text.Repeat(w, "]", extra); err != nil {
		return fmt.Errorf("unable to repeat args: %w", err)
	}
	return nil
}

var (
	errWrongArgumentCount    = errors.New("wrong argument count")
	errParseTakesNoArguments = fmt.Errorf("%w: no arguments permitted", errWrongArgumentCount)
)

// Validate checks if the correct number of positional arguments have been passed.
func (pos Positional) Validate(count int) error {
	// If we are outside the range for the arguments, we reset the counter to 0
	// and return the appropriate error message.
	//
	// - we always need to be more than the minimum
	// - we need to be below the max if the maximum is not unlimited
	if count < pos.Min || ((pos.Max != -1) && (count > pos.Max)) {
		switch {
		case pos.Min == 0 && pos.Max == 0: // 0 arguments, but some given
			return (errParseTakesNoArguments)
		case pos.Min == pos.Max: // exact number of arguments is wrong
			return fmt.Errorf("%w: exactly %d argument(s) required", errWrongArgumentCount, pos.Min)
		case pos.Max == -1: // less than min arguments
			return fmt.Errorf("%w: at least %d argument(s) required", errWrongArgumentCount, pos.Min)
		default: // between set number of arguments
			return fmt.Errorf("%w: between %d and %d argument(s) required", errWrongArgumentCount, pos.Min, pos.Max)
		}
	}

	return nil
}
