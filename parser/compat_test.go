package parser

import (
	"reflect"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/goprogram/meta"
)

func TestNewFlag(t *testing.T) {
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := NewFlag(tt.opt)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFlag() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
