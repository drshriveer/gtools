package benches

import (
	"slices"
	"testing"

	set2 "github.com/drshriveer/gtools/set"
)

// Takeaway--  basically anything <= 15 items should be just a normal for loop
//
//	otherwise use binary search.
func BenchmarkSetVsListOfIndexes(b *testing.B) {
	numItems := 15 // <- anything less than 20 items is faster to use slice. but binary is always faster.
	target := 12
	slice := make([]int, numItems)
	for i := range slice {
		slice[i] = i
	}
	st := set2.Make(slice...)
	b.Run("set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = st.Has(target)
		}
	})
	b.Run("indexed items", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range slice {
				if v == target {
					break
				}
			}
		}
	})
	b.Run("indexed items; binary", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slices.BinarySearch(slice, target)
		}
	})
}
