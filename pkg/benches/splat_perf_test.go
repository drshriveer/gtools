package benches_test

import "testing"

func BenchmarkSplats(b *testing.B) {
	v1, v2, v3, v4, v5 := 1, 2, 3, 4, 5
	vslice := []int{v1, v2, v3, v4, v5}
	b.Run("straight splat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = splatInput(v1, v2, v3, v4, v5)
		}
	},
	)
	b.Run("slice to splat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = splatInput(vslice...)
		}
	},
	)
	b.Run("straight slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sliceInput(vslice)
		}
	},
	)
	b.Run("constructed slice to splat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sliceInput([]int{v1, v2, v3, v4, v5})
		}
	},
	)
}

func splatInput(in ...int) int {
	sum := 0
	for _, v := range in {
		sum += v
	}
	return sum
}

func sliceInput(in []int) int {
	sum := 0
	for _, v := range in {
		sum += v
	}
	return sum
}
