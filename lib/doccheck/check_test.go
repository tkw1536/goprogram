// Package doccheck checks documentation strings at runtime for proper formatting.
//
// Checking is disabled by default, but can be enabled by building with the "doccheck" tag.
package doccheck

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantMessage string
	}{
		// failed checks
		{"empty part", "hello::world", "message \"hello::world\" failed validation: part \"\": empty"},
		{"part must start with ' 's", "hello:world", "message \"hello:world\" failed validation: part \"world\": missing a leading space"},
		{"may not have extra spaces", "hello: world  ", "message \"hello: world  \" failed validation: part \" world  \": contains extra space"},
		{"may not have periods", "world.", "message \"world.\" failed validation: part \"world.\": ends with a period"},
		{"may not start with upper case", "Hello World", "message \"Hello World\" failed validation: part \"Hello World\": starts with upper case"},
		{"may not start with upper case (2)", "HeLLo World", "message \"HeLLo World\" failed validation: part \"HeLLo World\": starts with upper case"},

		// passed checks
		{"empty string passes", "", ""},
		{"string with multiple parts passes", "something: something else: something else again", ""},
		{"string with entire uppercase word passes", "SOMETHING: something else: something else again", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)

			var gotMessage string
			if err != nil {
				gotMessage = err.Error()
			}

			if gotMessage != tt.wantMessage {
				t.Errorf("Validate() error = %q, wantErr %q", gotMessage, tt.wantMessage)
			}
		})
	}
}
