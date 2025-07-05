//spellchecker:words meta
package meta_test

//spellchecker:words strings testing
import (
	"strings"
	"testing"

	"go.tkw01536.de/goprogram/meta"
)

func TestPositional_WriteSpecTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		pos  meta.Positional
		want string
	}{
		{"arg 0, 0", meta.Positional{Value: "arg", Min: 0, Max: 0}, ""},
		{"arg 0, -1", meta.Positional{Value: "arg", Min: 0, Max: -1}, "[arg ...]"},
		{"arg 0, 3", meta.Positional{Value: "arg", Min: 0, Max: 3}, "[arg [arg [arg]]]"},

		{"no name 0, 0", meta.Positional{Value: "", Min: 0, Max: 0}, ""},
		{"no name 0, -1", meta.Positional{Value: "", Min: 0, Max: -1}, "[ARGUMENT ...]"},
		{"no name 0, 3", meta.Positional{Value: "", Min: 0, Max: 3}, "[ARGUMENT [ARGUMENT [ARGUMENT]]]"},

		{"arg 2, 2", meta.Positional{Value: "arg", Min: 2, Max: 2}, "arg arg"},
		{"arg 2, 4", meta.Positional{Value: "arg", Min: 2, Max: 4}, "arg arg [arg [arg]]"},
		{"arg 2, -1", meta.Positional{Value: "arg", Min: 2, Max: -1}, "arg arg [arg ...]"},

		{"no name 2, 2", meta.Positional{Value: "", Min: 2, Max: 2}, "ARGUMENT ARGUMENT"},
		{"no name 2, 4", meta.Positional{Value: "", Min: 2, Max: 4}, "ARGUMENT ARGUMENT [ARGUMENT [ARGUMENT]]"},
		{"no name 2, -1", meta.Positional{Value: "", Min: 2, Max: -1}, "ARGUMENT ARGUMENT [ARGUMENT ...]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var builder strings.Builder
			if err := tt.pos.WriteSpecTo(&builder); err != nil {
				t.Errorf("Positional.WriteSpecTo() returned non-nil error")
			}
			if got := builder.String(); got != tt.want {
				t.Errorf("Positional.WriteSpecTo() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPositional_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		Positional meta.Positional
		Count      int

		wantErr string
	}{
		// taking 0 args
		{
			"no arguments",
			meta.Positional{Min: 0, Max: 0},
			0,
			"",
		},

		// taking 1 arg
		{
			"one argument, too few",
			meta.Positional{Min: 1, Max: 1},
			0,
			"wrong argument count: exactly 1 argument(s) required",
		},
		{
			"one argument, exactly enough",
			meta.Positional{Min: 1, Max: 1},
			1,
			"",
		},
		{
			"one argument, too many",
			meta.Positional{Min: 1, Max: 1},
			2,
			"wrong argument count: exactly 1 argument(s) required",
		},

		// taking 1 or 2 args
		{
			"1-2 arguments, too few",
			meta.Positional{Min: 1, Max: 2},
			0,
			"wrong argument count: between 1 and 2 argument(s) required",
		},
		{
			"1-2 arguments, enough",
			meta.Positional{Min: 1, Max: 2},
			1,
			"",
		},
		{
			"1-2 arguments, enough (2)",
			meta.Positional{Min: 1, Max: 2},
			2,
			"",
		},
		{
			"1-2 arguments, too many",
			meta.Positional{Min: 1, Max: 2},
			3,
			"wrong argument count: between 1 and 2 argument(s) required",
		},

		// taking 2 args
		{
			"two arguments, too few",
			meta.Positional{Min: 2, Max: 2},
			0,
			"wrong argument count: exactly 2 argument(s) required",
		},
		{
			"two arguments, too few (2)",
			meta.Positional{Min: 2, Max: 2},
			1,
			"wrong argument count: exactly 2 argument(s) required",
		},
		{
			"two arguments, enough",
			meta.Positional{Min: 2, Max: 2},
			2,
			"",
		},
		{
			"two arguments, too many",
			meta.Positional{Min: 2, Max: 2},
			3,
			"wrong argument count: exactly 2 argument(s) required",
		},

		// at least one argument
		{
			"at least 1 arguments, not enough",
			meta.Positional{Min: 1, Max: -1},
			0,
			"wrong argument count: at least 1 argument(s) required",
		},
		{
			"at least 1 arguments, enough",
			meta.Positional{Min: 1, Max: -1},
			1,
			"",
		},
		{
			"at least 1 arguments, enough (2)",
			meta.Positional{Min: 1, Max: -1},
			2,
			"",
		},
		{
			"at least 1 arguments, enough (3)",
			meta.Positional{Min: 1, Max: -1},
			3,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.Positional.Validate(tt.Count)
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tt.wantErr {
				t.Errorf("Positional.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
