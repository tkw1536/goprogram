//spellchecker:words parser
package parser_test

//spellchecker:words reflect testing github jessevdk flags goprogram meta
import (
	"reflect"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/goprogram/meta"
	"github.com/tkw1536/goprogram/parser"
)

func TestNewFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opt  *flags.Option
		want meta.Flag
	}{
		{
			"simple option without default",
			&flags.Option{
				Required: true,

				ShortName: 's',
				LongName:  "long",

				ValueName:   "test",
				Description: "something",
				Default:     nil,
			},
			meta.Flag{
				Required: true,

				Short:     []string{"s"},
				Long:      []string{"long"},
				FieldName: "",

				Value:   "test",
				Usage:   "something",
				Default: "",

				Choices: nil,
			},
		},

		{
			"simple option with default",
			&flags.Option{
				Required: false,

				LongName: "long",

				ValueName:   "test",
				Description: "something",
				Default:     []string{"a"},
			},
			meta.Flag{
				Required: false,

				Long: []string{"long"},

				Value:   "test",
				Usage:   "something",
				Default: "a",
				Choices: nil,
			},
		},

		{
			"option with choices",
			&flags.Option{
				Required: false,

				LongName: "long",

				ValueName:   "test",
				Description: "something",
				Default:     []string{"a"},

				Choices: []string{"a", "b", "c"},
			},
			meta.Flag{
				Required: false,

				Long: []string{"long"},

				Value:   "test",
				Usage:   "something",
				Default: "a",

				Choices: []string{"a", "b", "c"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := parser.NewFlag(tt.opt)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFlag() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
