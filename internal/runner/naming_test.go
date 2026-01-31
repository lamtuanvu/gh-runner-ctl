package runner

import (
	"reflect"
	"testing"
)

func TestNextNumbers(t *testing.T) {
	tests := []struct {
		name     string
		existing []int
		count    int
		want     []int
	}{
		{"empty, add 3", nil, 3, []int{1, 2, 3}},
		{"1,2,5 exist, add 2", []int{1, 2, 5}, 2, []int{3, 4}},
		{"1,2,3 exist, add 1", []int{1, 2, 3}, 1, []int{4}},
		{"2,4 exist, add 3", []int{2, 4}, 3, []int{1, 3, 5}},
		{"none, add 0", nil, 0, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextNumbers(tt.existing, tt.count)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextNumbers(%v, %d) = %v, want %v", tt.existing, tt.count, got, tt.want)
			}
		})
	}
}

func TestHighestNumbers(t *testing.T) {
	tests := []struct {
		name     string
		existing []int
		count    int
		want     []int
	}{
		{"5,3,1,2,4, take 2", []int{5, 3, 1, 2, 4}, 2, []int{5, 4}},
		{"1,2,3, take 5", []int{1, 2, 3}, 5, []int{3, 2, 1}},
		{"10,1, take 1", []int{10, 1}, 1, []int{10}},
		{"empty, take 1", nil, 1, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HighestNumbers(tt.existing, tt.count)
			if len(got) != len(tt.want) {
				t.Errorf("HighestNumbers(%v, %d) = %v, want %v", tt.existing, tt.count, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("HighestNumbers(%v, %d) = %v, want %v", tt.existing, tt.count, got, tt.want)
					return
				}
			}
		})
	}
}
