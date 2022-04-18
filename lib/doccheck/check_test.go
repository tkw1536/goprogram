// Package doccheck checks documentation strings at runtime for proper formatting.
//
// Checking is disabled by default, but can be enabled by building with the "doccheck" tag.
package doccheck

import (
	"reflect"
	"testing"

	"github.com/tkw1536/goprogram/lib/testlib"
)

var validateTests = []struct {
	name      string
	input     string
	wantError interface{}
}{
	// failed checks
	{"empty part", "hello::world", &CheckError{Message: "hello::world", Part: "", Index: 1, Failure: "empty"}},
	{"part must start with ' 's", "hello:world", &CheckError{Message: "hello:world", Part: "world", Index: 1, Failure: "missing a leading space"}},
	{"may not have extra spaces", "hello: world  ", &CheckError{Message: "hello: world  ", Part: " world  ", Index: 1, Failure: "contains extra space"}},
	{"may not have periods", "world.", &CheckError{Message: "world.", Part: "world.", Index: 0, Failure: "ends with a period"}},
	{"may not start with upper case", "Hello World", &CheckError{Message: "Hello World", Part: "Hello World", Index: 0, Failure: "starts with upper case"}},
	{"may not start with upper case (2)", "HeLLo World", &CheckError{Message: "HeLLo World", Part: "HeLLo World", Index: 0, Failure: "starts with upper case"}},

	// passed checks
	{"empty string passes", "", nil},
	{"string with multiple parts passes", "something: something else: something else again", nil},
	{"string with entire uppercase word passes", "SOMETHING: something else: something else again", nil},
}

func TestValidate(t *testing.T) {
	for _, tt := range validateTests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := Validate(tt.input)

			if !reflect.DeepEqual(gotErr, tt.wantError) {
				t.Errorf("Validate() error = %#v, want = %#v", gotErr, tt.wantError)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	for _, tt := range validateTests {
		t.Run(tt.name, func(t *testing.T) {
			var wantPanic bool
			var wantError interface{}

			if enabled {
				wantPanic = tt.wantError != nil
				wantError = tt.wantError
			}

			gotPanic, gotError := testlib.DoesPanic(func() {
				Check(tt.input)
			})

			if gotPanic != wantPanic {
				t.Errorf("Check() got panic = %v, want = %v", gotPanic, wantPanic)
			}

			if !reflect.DeepEqual(gotError, wantError) {
				t.Errorf("Check() got error = %v, want = %v", gotError, wantError)
			}
		})
	}
}
