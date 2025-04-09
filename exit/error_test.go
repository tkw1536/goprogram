//spellchecker:words exit
package exit_test

//spellchecker:words errors reflect testing github pkglib testlib
import (
	"errors"
	"reflect"
	"testing"

	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/testlib"
)

var errNotAnError = errors.New("test: not an error")

func TestAsErrorPanic(t *testing.T) {
	t.Parallel()

	_, gotPanic := testlib.DoesPanic(func() { _ = exit.AsError(errNotAnError) })
	wantPanic := interface{}("AsError: err must be nil or wrap type Error")
	if wantPanic != gotPanic {
		t.Errorf("AsError: got panic = %v, want = %v", gotPanic, wantPanic)
	}
}

func TestError_WithMessage(t *testing.T) {
	t.Parallel()

	type fields struct {
		ExitCode exit.ExitCode
		Message  string
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   exit.Error
	}{
		{"replaces empty message", fields{}, args{message: "hello world"}, exit.Error{Message: "hello world"}},
		{"replaces non-empty message", fields{Message: "not empty"}, args{message: "hello world"}, exit.Error{Message: "hello world"}},

		{"keeps exit code 1", fields{ExitCode: 1}, args{message: "hello world"}, exit.Error{ExitCode: 1, Message: "hello world"}},
		{"keeps exit code 2", fields{ExitCode: 2}, args{message: "hello world"}, exit.Error{ExitCode: 2, Message: "hello world"}},
		{"keeps exit code 3", fields{ExitCode: 3}, args{message: "hello world"}, exit.Error{ExitCode: 3, Message: "hello world"}},
		{"keeps exit code 4", fields{ExitCode: 4}, args{message: "hello world"}, exit.Error{ExitCode: 4, Message: "hello world"}},
		{"keeps exit code 5", fields{ExitCode: 5}, args{message: "hello world"}, exit.Error{ExitCode: 5, Message: "hello world"}},

		{"does not substitute strings in old message", fields{Message: "old %s"}, args{message: "hello world"}, exit.Error{Message: "hello world"}},
		{"does not substitute strings in new message", fields{Message: "old message"}, args{message: "hello world %s"}, exit.Error{Message: "hello world %s"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := exit.Error{
				ExitCode: tt.fields.ExitCode,
				Message:  tt.fields.Message,
			}
			if got := err.WithMessage(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.WithMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WithMessageF(t *testing.T) {
	t.Parallel()

	type fields struct {
		ExitCode exit.ExitCode
		Message  string
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   exit.Error
	}{
		{"keeps message without format", fields{Message: "hello world"}, args{}, exit.Error{Message: "hello world"}},
		{"replaces message", fields{Message: "hello %s"}, args{[]interface{}{"world"}}, exit.Error{Message: "hello world"}},

		{"keeps exit code 1", fields{ExitCode: 1, Message: "%s"}, args{[]interface{}{"hello world"}}, exit.Error{ExitCode: 1, Message: "hello world"}},
		{"keeps exit code 2", fields{ExitCode: 2, Message: "%s"}, args{[]interface{}{"hello world"}}, exit.Error{ExitCode: 2, Message: "hello world"}},
		{"keeps exit code 3", fields{ExitCode: 3, Message: "%s"}, args{[]interface{}{"hello world"}}, exit.Error{ExitCode: 3, Message: "hello world"}},
		{"keeps exit code 4", fields{ExitCode: 4, Message: "%s"}, args{[]interface{}{"hello world"}}, exit.Error{ExitCode: 4, Message: "hello world"}},
		{"keeps exit code 5", fields{ExitCode: 5, Message: "%s"}, args{[]interface{}{"hello world"}}, exit.Error{ExitCode: 5, Message: "hello world"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := exit.Error{
				ExitCode: tt.fields.ExitCode,
				Message:  tt.fields.Message,
			}
			if got := err.WithMessageF(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.WithMessageF() = %v, want %v", got, tt.want)
			}
		})
	}
}
