package meta

import (
	"fmt"
	"io"

	"github.com/tkw1536/goprogram/lib/text"
)

// Positional holds meta-information about a positional argument.
type Positional struct {
	// Name and Description of the positional in help texts
	Value string // defaults to "ARGUMENT"
	Usage string

	// Min and Max indicate how many positional arguments are expected for this command.
	// Min must be >= 0. Max must be either Min, or -1.
	// Max == -1 inidicates an unlimited number of repeats.
	Min, Max int
}

// ValidRange checks if posititional has valid min and max values.
func (pos Positional) ValidRange() bool {
	extra := pos.Max - pos.Min
	return pos.Min >= 0 && (pos.Max <= 0 || extra >= 0)
}

// defaultPositionalValue is the default name used for a positional argument.
const defaultPositionalValue = "ARGUMENT"

// WriteSpecTo writes a specification of this argument into w.
// A specification looks like "arg [arg...]".
func (pos Positional) WriteSpecTo(w io.Writer) {
	extra := pos.Max - pos.Min

	if !pos.ValidRange() {
		panic("Positional: invalid range")
	}

	if pos.Value == "" {
		pos.Value = defaultPositionalValue
	}

	// nothing to generate!
	if pos.Max == 0 && extra == 0 {
		return
	}

	// arg arg arg
	text.RepeatJoin(w, pos.Value, " ", pos.Min)
	if pos.Min > 0 && extra != 0 {
		io.WriteString(w, " ")
	}

	if pos.Max < 0 {
		// [arg ...]
		io.WriteString(w, "[")
		io.WriteString(w, pos.Value)
		io.WriteString(w, " ...]")
		return
	}

	// [arg [arg]]
	text.RepeatJoin(w, "["+pos.Value, " ", extra)
	text.Repeat(w, "]", extra)
}

const (
	errParseTakesNoArguments      = "no arguments permitted"
	errParseTakesExactlyArguments = "exactly %d argument(s) required"
	errParseTakesMinArguments     = "at least %d argument(s) required"
	errParseTakesBetweenArguments = "between %d and %d argument(s) required"
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
		case pos.Min == pos.Max && pos.Min == 0: // 0 arguments, but some given
			return fmt.Errorf(errParseTakesNoArguments)
		case pos.Min == pos.Max: // exact number of arguments is wrong
			return fmt.Errorf(errParseTakesExactlyArguments, pos.Min)
		case pos.Max == -1: // less than min arguments
			return fmt.Errorf(errParseTakesMinArguments, pos.Min)
		default: // between set number of arguments
			return fmt.Errorf(errParseTakesBetweenArguments, pos.Min, pos.Max)
		}
	}

	return nil
}
