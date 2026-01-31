package runner

import "sort"

// NextNumbers finds the next `count` available runner numbers by filling gaps
// in the existing number sequence starting from 1.
func NextNumbers(existing []int, count int) []int {
	used := make(map[int]bool, len(existing))
	for _, n := range existing {
		used[n] = true
	}

	var result []int
	candidate := 1
	for len(result) < count {
		if !used[candidate] {
			result = append(result, candidate)
		}
		candidate++
	}
	return result
}

// HighestNumbers returns the highest `count` numbers from the given slice,
// sorted descending. Used by `down` to remove highest-numbered runners first.
func HighestNumbers(existing []int, count int) []int {
	sorted := make([]int, len(existing))
	copy(sorted, existing)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))

	if count > len(sorted) {
		count = len(sorted)
	}
	return sorted[:count]
}
