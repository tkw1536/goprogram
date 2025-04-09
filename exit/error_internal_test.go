package exit

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

//spellchecker:words nolint testpackage

func TestAsError(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			if got := AsError(tt.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsError() = %v, want %v", got, tt.want)
			}
		})
	}
}

var errTestInner = errors.New("inner error")

func TestError_WrapError(t *testing.T) {
	t.Parallel()

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
		{"wraps an inner error", fields{ExitCode: 1, Message: "something went wrong"}, args{errTestInner}, Error{ExitCode: 1, Message: "something went wrong: inner error", err: errTestInner}},
		{"wraps a nil error", fields{ExitCode: 1, Message: "something went wrong"}, args{nil}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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

func TestError_Wrap(t *testing.T) {
	t.Parallel()

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
		{"wraps an inner error", fields{ExitCode: 1, Message: "something went wrong"}, args{errTestInner}, Error{ExitCode: 1, Message: "something went wrong: inner error", err: errTestInner}},
		{"wraps a nil error", fields{ExitCode: 1, Message: "something went wrong"}, args{nil}, Error{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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
