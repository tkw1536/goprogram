package docfmt

import (
	"reflect"
	"testing"

	"github.com/tkw1536/goprogram/lib/testlib"
)

func TestCheck(t *testing.T) {
	for _, tt := range partTests {
		t.Run(tt.name, func(t *testing.T) {
			var wantPanic bool
			var wantError interface{}

			if enabled {
				wantPanic = tt.wantError != nil
				if wantPanic {
					wantError = &CheckError{
						Message: tt.input,
						Results: tt.wantError,
					}
				}
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
