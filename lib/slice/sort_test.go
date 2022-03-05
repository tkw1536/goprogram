package slice

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  []int
	}{
		{"nil sort", nil, nil},
		{"empty sort", []int{}, []int{}},
		{"simple sort", []int{7, 3, 0, 4, 9, 2, 1, 5, 6, 8}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"duplicate sort", []int{7, 7, 3, 4, 4, 4, 3, 3, 7}, []int{3, 3, 3, 4, 4, 4, 7, 7, 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Sort(tt.slice)
			if !reflect.DeepEqual(tt.slice, tt.want) {
				t.Errorf("Sort() = %v, want %v", tt.slice, tt.want)
			}
		})
	}
}

func TestSortFunc(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		less  func(int, int) bool
		want  []int
	}{
		{"nil sort", nil, nil, nil},
		{"empty sort", []int{}, func(i1, i2 int) bool { panic("never reached") }, []int{}},
		{"normal sort", []int{7, 3, 0, 4, 9, 2, 1, 5, 6, 8}, func(i, j int) bool { return i < j }, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"reverse sort", []int{7, 3, 0, 4, 9, 2, 1, 5, 6, 8}, func(i, j int) bool { return j < i }, []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortFunc(tt.slice, tt.less)
			if !reflect.DeepEqual(tt.slice, tt.want) {
				t.Errorf("SortFunc() = %v, want %v", tt.slice, tt.want)
			}
		})
	}
}
