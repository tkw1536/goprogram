package meta

import (
	"strings"
	"testing"
)

//spellchecker:words nolint testpackage

func TestMeta_writeCommandsTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		meta Meta
		want string
	}{
		{"no commands", Meta{Commands: nil}, ""},
		{"single command", Meta{Commands: []string{"a"}}, `"a"`},
		{"multiple commands", Meta{Commands: []string{"a", "b", "c"}}, `"a", "b", "c"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var builder strings.Builder
			if err := tt.meta.writeCommandsTo(&builder); err != nil {
				t.Errorf("writeCommandsTo() returned non-nil error")
			}
			if got := builder.String(); got != tt.want {
				t.Errorf("writeCommandsTo() = %v, want %v", got, tt.want)
			}
		})
	}
}
