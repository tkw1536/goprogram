// Package reflectx provides extensions to the reflect package
package reflectx

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypeOf(t *testing.T) {
	tests := []struct {
		name string
		got  reflect.Type
		want reflect.Type
	}{
		{
			"string",
			TypeOf[string](),
			reflect.TypeOf(string("")),
		},
		{
			"int",
			TypeOf[int](),
			reflect.TypeOf(int(0)),
		},
		{
			"slice",
			TypeOf[[]string](),
			reflect.TypeOf([]string(nil)),
		},
		{
			"array",
			TypeOf[[0]string](),
			reflect.TypeOf([0]string{}),
		},
		{
			"chan",
			TypeOf[chan string](),
			reflect.TypeOf((chan string)(nil)),
		},
		{
			"func",
			TypeOf[func(string) string](),
			reflect.TypeOf(func(string) string { return "" }),
		},
		{
			"map",
			TypeOf[map[string]string](),
			reflect.TypeOf(map[string]string(nil)),
		},
		{
			"struct",
			TypeOf[struct{ Thing string }](),
			reflect.TypeOf(struct{ Thing string }{}),
		},
		{
			"pointer",
			TypeOf[*struct{ Thing string }](),
			reflect.TypeOf(&struct{ Thing string }{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("TypeOf() = %v, want %v", tt.got, tt.want)
			}
		})
	}
}

// Iterate over the fields of a struct
func ExampleIterateFields() {

	type Embed struct {
		EmbeddedField string // field in an embedded struct
	}

	type SomeStruct struct {
		Field   string // regular field
		string         // embedded non-struct, not called recursively
		Embed          // an embed
		Another string //
	}

	fmt.Println(
		"returned:",
		IterateFields(TypeOf[SomeStruct](), func(f reflect.StructField, index int) (cancel bool) {
			fmt.Println("encounted field", f.Name, "with index", index)
			return false // do not cancel!
		}),
	)

	// Output: encounted field Field with index 0
	// encounted field string with index 1
	// encounted field Embed with index 2
	// encounted field Another with index 3
	// returned: false
}

// Iterate over the fields of a struct
func ExampleIterateAllFields() {

	type Embed struct {
		EmbeddedField string // field in an embedded struct
	}

	type SomeStruct struct {
		Field   string // regular field
		string         // embedded non-struct, not called recursively
		Embed          // an embed
		Another string //
	}

	fmt.Println(
		"returned:",
		IterateAllFields(TypeOf[SomeStruct](), func(f reflect.StructField, index ...int) (cancel bool) {
			fmt.Println("encounted field", f.Name, "with index", index)
			return false // do not cancel!
		}),
	)

	// Output: encounted field Field with index [0]
	// encounted field string with index [1]
	// encounted field EmbeddedField with index [2 0]
	// encounted field Another with index [3]
	// returned: false
}

// stop the iteration over a subset of fields only
func ExampleIterateAllFields_cancel() {

	type Embed struct {
		EmbeddedField string // field in an embedded struct
	}

	type SomeStruct struct {
		Field   string // regular field
		string         // embedded non-struct, not called recursively
		Embed          // an embed
		Another string //
	}

	fmt.Println(
		"returned:",
		IterateAllFields(TypeOf[SomeStruct](), func(f reflect.StructField, index ...int) (cancel bool) {
			fmt.Println("encounted field", f.Name, "with index", index)
			return f.Name == "EmbeddedField" // cancel on embedded field
		}),
	)

	// Output: encounted field Field with index [0]
	// encounted field string with index [1]
	// encounted field EmbeddedField with index [2 0]
	// returned: true
}
