package slice

import (
	"fmt"
	"reflect"
	"testing"
)

func TestContainsAny(t *testing.T) {
	type args struct {
		haystack []string
		needles  []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"single needle contained in haystack", args{
			haystack: []string{"a", "b", "c"},
			needles:  []string{"a"},
		}, true},
		{"single needle not contained in haystack", args{
			haystack: []string{"a", "b", "c"},
			needles:  []string{"d"},
		}, false},
		{"haystack contains a single needle", args{
			haystack: []string{"a", "b", "c"},
			needles:  []string{"f", "a", "e"},
		}, true},
		{"haystack contains no needle", args{
			haystack: []string{"a", "b", "c"},
			needles:  []string{"d", "e", "f"},
		}, false},
		{"empty haystack", args{
			haystack: nil,
			needles:  []string{"d", "e", "f"},
		}, false},
		{"empty needles", args{
			haystack: []string{"a", "b", "c"},
			needles:  nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAny(tt.args.haystack, tt.args.needles...); got != tt.want {
				t.Errorf("ContainsAny() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMatchesAny(t *testing.T) {

	// p makes one predicate for each element in haystack
	p := func(haystack ...string) (fs []func(string) bool) {
		for _, hay := range haystack {
			h := hay
			fs = append(fs, func(s string) bool {
				return s == h
			})
		}
		return
	}

	type args struct {
		haystack   []string
		predicates []func(string) bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"single needle contained in haystack", args{
			haystack:   []string{"a", "b", "c"},
			predicates: p("a"),
		}, true},
		{"single needle not contained in haystack", args{
			haystack:   []string{"a", "b", "c"},
			predicates: p("d"),
		}, false},
		{"haystack contains a single needle", args{
			haystack:   []string{"a", "b", "c"},
			predicates: p("f", "a", "e"),
		}, true},
		{"haystack contains no needle", args{
			haystack:   []string{"a", "b", "c"},
			predicates: p("d", "e", "f"),
		}, false},
		{"empty haystack", args{
			haystack:   nil,
			predicates: p("d", "e", "f"),
		}, false},
		{"empty needles", args{
			haystack:   []string{"a", "b", "c"},
			predicates: nil,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchesAny(tt.args.haystack, tt.args.predicates...); got != tt.want {
				t.Errorf("MatchesAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEquals(t *testing.T) {
	type args struct {
		first  []string
		second []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{

		{"equality of nil slices", args{nil, nil}, true},
		{"equality of nil and empty slice", args{nil, []string{}}, true},
		{"equality of empty and nil slice", args{[]string{}, nil}, true},
		{"equality of empty slices", args{[]string{}, []string{}}, true},

		{"inequality of empty and full slice", args{nil, []string{"a"}}, false},
		{"inequality of full and empty slice", args{[]string{"a"}, nil}, false},

		{"equality of full slices", args{[]string{"a"}, []string{"a"}}, true},
		{"equality of full slices (2)", args{[]string{"a", "b", "c"}, []string{"a", "b", "c"}}, true},

		{"inequality of full slices", args{[]string{"a"}, []string{"a", "b", "c"}}, false},
		{"inequality of full slices (2)", args{[]string{"a", "b", "c"}, []string{"a"}}, false},
		{"inequality of full slices (3)", args{[]string{"a", "b", "c"}, []string{"a", "d", "c"}}, false},
		{"inequality of full slices (4)", args{[]string{"a", "d", "c"}, []string{"a", "b", "c"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equals(tt.args.first, tt.args.second); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	type args struct {
		slice []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"nil slice", args{nil}, nil},

		{"non-empty slice", args{[]string{"a", "b", "c"}}, []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Copy(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleCopy() {
	// create an "old world" slice
	old_world := []string{"hello", "old", "world"}

	// make a copy and change old to new
	new_world := Copy(old_world)
	new_world[1] = "new"

	// old world remains the same
	fmt.Println(old_world)
	fmt.Println(new_world)

	// Output: [hello old world]
	// [hello new world]
}

func TestFilterI(t *testing.T) {
	type args struct {
		slice []int
		pred  func(int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "nil list",
			args: args{slice: nil, pred: func(i int) bool { panic("never reached") }},
			want: nil,
		},

		{
			name: "true filter",
			args: args{slice: []int{0, 1, 2, 3}, pred: func(i int) bool { return true }},
			want: []int{0, 1, 2, 3},
		},
		{
			name: "false filter",
			args: args{slice: []int{0, 1, 2, 3}, pred: func(i int) bool { return false }},
			want: []int{},
		},

		{
			name: "even filter",
			args: args{slice: []int{0, 1, 2, 3}, pred: func(i int) bool { return i%2 == 0 }},
			want: []int{0, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterI(tt.args.slice, tt.args.pred); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleFilterI() {

	// create a slice and filter it in-place!
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	filtered := FilterI(slice, func(i int) bool { return i%2 == 0 })

	// the filtered slice is as we would expect
	fmt.Println(filtered)

	// the original slice has been invalidated, elements not used have been zeroed out.
	fmt.Println(slice)

	// this has invalidated the outer slice

	// because we filtered in place, slice[0:6] is the same as filtered[0:6]
	// we show this by setting all of slice and printing it again
	slice[0] = -1
	slice[1] = -1
	slice[2] = -1
	slice[3] = -1
	slice[4] = -1
	slice[5] = -1
	fmt.Println(filtered)

	// Normally one would just
	//  slice = FilterI(slice, ...)
	// to prevent accidentally leaking memory.

	// Output: [0 2 4 6 8]
	// [0 2 4 6 8 0 0 0 0 0]
	// [-1 -1 -1 -1 -1]
}

func BenchmarkFilterI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FilterI([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, func(i int) bool { return i%2 == 0 })
	}
}

func TestRemoveZeros(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"remove from the nil slice", args{nil}, nil},
		{"remove from the empty array", args{[]string{}}, []string{}},
		{"remove from some places", args{[]string{"", "x", "y", "", "z"}}, []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveZeros(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveZeros() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkRemoveZeros(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RemoveZeros[struct{}](nil)
		RemoveZeros([]string{})
		RemoveZeros([]string{"", "x", "y", "", "z"})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"nil slice", args{nil}, nil},
		{"no duplicates", args{[]string{"a", "b", "c", "d"}}, []string{"a", "b", "c", "d"}},
		{"some duplicates", args{[]string{"b", "c", "c", "d", "a", "b", "c", "d"}}, []string{"a", "b", "c", "d"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicates(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}
