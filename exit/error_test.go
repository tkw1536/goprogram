//spellchecker:words exit
package exit

//spellchecker:words errors reflect testing github pkglib testlib
import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/tkw1536/pkglib/testlib"
)

func TestAsError(t *testing.T) {
	var errStuff = Error{ExitCode: ExitGeneric, Message: "stuff"}
	var errStuffWrapped = fmt.Errorf("wrapping: %w", errStuff)
	var errWrapped = Error{ExitCode: ExitGeneric, Message: "wrapping: stuff", err: errStuffWrapped}

	tests := []struct {
		name string
		err  error
		want Error
	}{
		{
			name: "nil error returns zero value",
			err:  nil,
			want: Error{},
		},
		{
			name: "Error object returns itself",
			err:  errStuff,
			want: errStuff,
		},
		{
			name: "Wrapped error returns same exit code",
			err:  errStuffWrapped,
			want: errWrapped,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsError(tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsErrorPanic(t *testing.T) {
	_, gotPanic := testlib.DoesPanic(func() { _ = AsError(errors.New("not an error")) })
	wantPanic := interface{}("AsError: err must be nil or wrap type Error")
	if wantPanic != gotPanic {
		t.Errorf("AsError: got panic = %v, want = %v", gotPanic, wantPanic)
	}
}

func TestError_WithMessage(t *testing.T) {
	type fields struct {
		ExitCode ExitCode
		Message  string
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Error
	}{
		{"replaces empty message", fields{}, args{message: "hello world"}, Error{Message: "hello world"}},
		{"replaces non-empty message", fields{Message: "not empty"}, args{message: "hello world"}, Error{Message: "hello world"}},

		{"keeps exit code 1", fields{ExitCode: 1}, args{message: "hello world"}, Error{ExitCode: 1, Message: "hello world"}},
		{"keeps exit code 2", fields{ExitCode: 2}, args{message: "hello world"}, Error{ExitCode: 2, Message: "hello world"}},
		{"keeps exit code 3", fields{ExitCode: 3}, args{message: "hello world"}, Error{ExitCode: 3, Message: "hello world"}},
		{"keeps exit code 4", fields{ExitCode: 4}, args{message: "hello world"}, Error{ExitCode: 4, Message: "hello world"}},
		{"keeps exit code 5", fields{ExitCode: 5}, args{message: "hello world"}, Error{ExitCode: 5, Message: "hello world"}},

		{"does not substitute strings in old message", fields{Message: "old %s"}, args{message: "hello world"}, Error{Message: "hello world"}},
		{"does not substitute strings in new message", fields{Message: "old message"}, args{message: "hello world %s"}, Error{Message: "hello world %s"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Error{
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
	type fields struct {
		ExitCode ExitCode
		Message  string
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Error
	}{
		{"keeps message without format", fields{Message: "hello world"}, args{}, Error{Message: "hello world"}},
		{"replaces message", fields{Message: "hello %s"}, args{[]interface{}{"world"}}, Error{Message: "hello world"}},

		{"keeps exit code 1", fields{ExitCode: 1, Message: "%s"}, args{[]interface{}{"hello world"}}, Error{ExitCode: 1, Message: "hello world"}},
		{"keeps exit code 2", fields{ExitCode: 2, Message: "%s"}, args{[]interface{}{"hello world"}}, Error{ExitCode: 2, Message: "hello world"}},
		{"keeps exit code 3", fields{ExitCode: 3, Message: "%s"}, args{[]interface{}{"hello world"}}, Error{ExitCode: 3, Message: "hello world"}},
		{"keeps exit code 4", fields{ExitCode: 4, Message: "%s"}, args{[]interface{}{"hello world"}}, Error{ExitCode: 4, Message: "hello world"}},
		{"keeps exit code 5", fields{ExitCode: 5, Message: "%s"}, args{[]interface{}{"hello world"}}, Error{ExitCode: 5, Message: "hello world"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Error{
				ExitCode: tt.fields.ExitCode,
				Message:  tt.fields.Message,
			}
			if got := err.WithMessageF(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.WithMessageF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Wrap(t *testing.T) {
	var inner = errors.New("inner error")
	type fields struct {
		ExitCode ExitCode
		Message  string
		err      error
	}
	type args struct {
		inner error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Error
	}{
		{"wraps an inner error", fields{ExitCode: 1, Message: "something went wrong"}, args{inner}, Error{ExitCode: 1, Message: "something went wrong: inner error", err: inner}},
		{"wraps a nil error", fields{ExitCode: 1, Message: "something went wrong"}, args{nil}, Error{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Error{
				ExitCode: tt.fields.ExitCode,
				Message:  tt.fields.Message,
				err:      tt.fields.err,
			}
			if got := err.Wrap(tt.args.inner); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_WrapError(t *testing.T) {
	var inner = errors.New("inner error")
	type fields struct {
		ExitCode ExitCode
		Message  string
		err      error
	}
	type args struct {
		inner error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
	}{
		{"wraps an inner error", fields{ExitCode: 1, Message: "something went wrong"}, args{inner}, Error{ExitCode: 1, Message: "something went wrong: inner error", err: inner}},
		{"wraps a nil error", fields{ExitCode: 1, Message: "something went wrong"}, args{nil}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Error{
				ExitCode: tt.fields.ExitCode,
				Message:  tt.fields.Message,
				err:      tt.fields.err,
			}
			if got := err.WrapError(tt.args.inner); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error.Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleError_DeferWrap() {
	var genericError = Error{ExitCode: ExitGeneric, Message: "generic error"}

	// something returns the error it is passed
	something := func(in error) (err error) {
		// ensure that err is of type Error!
		// this only updates error which are not yet of type Error.
		defer genericError.DeferWrap(&err)

		return in
	}

	fmt.Println(something(nil))
	fmt.Println(something(errors.New("something went wrong")))
	fmt.Println(something(Error{ExitCode: ExitGeneric, Message: "specific error"}))

	// output: <nil>
	// generic error: something went wrong
	// specific error
}
