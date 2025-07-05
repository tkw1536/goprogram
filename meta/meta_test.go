//spellchecker:words meta
package meta_test

//spellchecker:words strings testing
import (
	"strings"
	"testing"

	"go.tkw01536.de/goprogram/meta"
)

//spellchecker:words positionals

func TestUsage_WriteMessageTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		meta meta.Meta
		want string
	}{
		{
			"main executable page",
			meta.Meta{
				Executable:  "cmd",
				Description: "do something interesting",

				GlobalFlags: []meta.Flag{
					{
						Required: true,

						Short: []string{"g"},
						Long:  []string{"global"},

						Value:   "name",
						Usage:   "a global argument",
						Default: "",
					},
					{
						Required: false,

						Short:   []string{"q"},
						Long:    []string{"quiet"},
						Usage:   "be quiet",
						Default: "false",
					},
				},
				Commands: []string{"a", "b", "c"},
			},
			"Usage: cmd --global|-g name [--quiet|-q] [--] COMMAND [ARGS...]\n\ndo something interesting\n\n   -g, --global name\n      a global argument\n\n   -q, --quiet\n      be quiet (default false)\n\n   COMMAND [ARGS...]\n      Command to call. One of \"a\", \"b\", \"c\". See individual commands for more help.",
		},
		{
			"sub executable page",
			meta.Meta{
				Executable:  "cmd",
				Command:     "sub",
				Description: "do something local",

				GlobalFlags: []meta.Flag{
					{
						Required: true,

						Short: []string{"g"},
						Long:  []string{"global"},

						Value:   "name",
						Usage:   "a global argument",
						Default: "",
					},
					{
						Required: false,

						Short:   []string{"q"},
						Long:    []string{"quiet"},
						Usage:   "be quiet",
						Default: "false",
					},
				},
				CommandFlags: []meta.Flag{
					{
						Required: true,

						Short: []string{"d"},
						Long:  []string{"dud"},

						Value:   "dud",
						Usage:   "a local argument",
						Default: "",
					},
					{
						Required: false,

						Short:   []string{"s"},
						Long:    []string{"silent"},
						Usage:   "be silent",
						Default: "true",
					},
				},
				Positionals: []meta.Positional{
					{
						Value: "op",
						Usage: "operations to make",
						Min:   1,
						Max:   -1,
					},
				},
			},
			"Usage: cmd --global|-g name [--quiet|-q] [--] sub --dud|-d dud [--silent|-s] [--] op [op ...]\n\ndo something local\n\nGlobal Arguments:\n\n   -g, --global name\n      a global argument\n\n   -q, --quiet\n      be quiet (default false)\n\nCommand Arguments:\n\n   -d, --dud dud\n      a local argument\n\n   -s, --silent\n      be silent (default true)\n\n   op [op ...]\n      operations to make",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var builder strings.Builder
			if err := tt.meta.WriteMessageTo(&builder); err != nil {
				t.Errorf("Usage.WriteMessageTo() returned non-nil error")
			}
			if got := builder.String(); got != tt.want {
				t.Errorf("Usage.WriteMessageTo() = %q, want %q", got, tt.want)
			}
		})
	}
}
