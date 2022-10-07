package collection

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// First returns the first value in S for which test returns true.
// When no such value exists, returns the zero value of T.
func First[T any](S []T, test func(T) bool) T {
	for _, v := range S {
		if test(v) {
			return v
		}
	}
	var v T
	return v
}

// Any returns true if there is some value in S for which test returns true
// and false otherwise.
func Any[T any](S []T, test func(T) bool) bool {
	return slices.IndexFunc(S, test) >= 0
}

// Filter modifies S in-place to those elements where filter returns true.
// Filter never re-allocates, invalidating the previous value of S.
// See also [FilterClone].
func Filter[T any](S []T, filter func(T) bool) []T {
	results := S[:0]
	for _, value := range S {
		if filter(value) {
			results = append(results, value)
		}
	}
	var zero T
	for i := len(results); i < len(S); i++ {
		S[i] = zero
	}
	return results
}

// FilterClone behaves like [Filter], except that it generates a new slice
// and keeps the original reference S valid.
func FilterClone[T any](S []T, filter func(T) bool) (results []T) {
	for _, value := range S {
		if filter(value) {
			results = append(results, value)
		}
	}
	return
}

// MapSlice generates a new slice by passing each element of S through f.
func MapSlice[T1, T2 any](S []T1, f func(T1) T2) []T2 {
	if S == nil {
		return nil
	}

	results := make([]T2, len(S))
	for i, v := range S {
		results[i] = f(v)
	}
	return results
}

// AsAny returns a new slice containing the same elements as S, but as any.
func AsAny[T any](S []T) []any {
	// NOTE(twiesing): This function is untested because MapSlice is tested.
	return MapSlice(S, func(t T) any { return t })
}

// Partition partitions elements of S into sub-slices by passing them through f.
// Order of elements within subslices is preserved.
func Partition[T any, P comparable](S []T, f func(T) P) map[P][]T {
	result := make(map[P][]T)
	for _, v := range S {
		key := f(v)
		result[key] = append(result[key], v)
	}
	return result
}

// NonSequential sorts values, and then removes sequential elements for which test() returns true.
// NonSequential does not re-allocate, but uses the existing slice.
func NonSequential[T constraints.Ordered](values []T, test func(prev, current T) bool) []T {
	if len(values) < 2 {
		return values
	}

	// sort the values and make a results array
	slices.Sort(values)
	results := values[:1]

	// do the filter loop
	prev := results[0]
	for _, current := range values[1:] {
		if !test(prev, current) {
			results = append(results, current)
		}
		prev = current
	}

	return results
}
