package internal_test

import (
	"sync"
	"testing"

	"github.com/drshriveer/gcommon/pkg/gsync"
	"github.com/drshriveer/gcommon/pkg/gsync/internal"
)

// BenchmarkSyncWaitGroup_Add
// BenchmarkSyncWaitGroup_Add/single
// BenchmarkSyncWaitGroup_Add/single-10    	171599598	         6.923 ns/op
// BenchmarkSyncWaitGroup_Add/parallel
// BenchmarkSyncWaitGroup_Add/parallel-10  	17776362	        67.32 ns/op
// BenchmarkSelectableWaitGroup1_Add
// BenchmarkSelectableWaitGroup1_Add/single
// BenchmarkSelectableWaitGroup1_Add/single-10         	64213616	        18.62 ns/op
// BenchmarkSelectableWaitGroup1_Add/parallel
// BenchmarkSelectableWaitGroup1_Add/parallel-10       	10303479	       118.9 ns/op
// BenchmarkSelectableWaitGroup2_Add
// BenchmarkSelectableWaitGroup2_Add/single
// BenchmarkSelectableWaitGroup2_Add/single-10         	86082441	        13.94 ns/op
// BenchmarkSelectableWaitGroup2_Add/parallel
// BenchmarkSelectableWaitGroup2_Add/parallel-10       	 6170856	       208.1 ns/op
// BenchmarkSelectableWaitGroup3_Add
// BenchmarkSelectableWaitGroup3_Add/single
// BenchmarkSelectableWaitGroup3_Add/single-10         	56798234	        20.82 ns/op
// BenchmarkSelectableWaitGroup3_Add/parallel
// BenchmarkSelectableWaitGroup3_Add/parallel-10       	 3100204	       384.2 ns/op
// BenchmarkSelectableWaitGroup4_Add
// BenchmarkSelectableWaitGroup4_Add/single
// BenchmarkSelectableWaitGroup4_Add/single-10         	174314197	         6.879 ns/op
// BenchmarkSelectableWaitGroup4_Add/parallel
// BenchmarkSelectableWaitGroup4_Add/parallel-10       	17471242	        70.04 ns/op
func BenchmarkSyncWaitGroup_Add(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		wg := &sync.WaitGroup{}
		for i := 0; i < b.N; i++ {
			wg.Add(1)
		}
	})
	b.Run("parallel", func(b *testing.B) {
		wg := &sync.WaitGroup{}
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wg.Add(1)
			}
		})
	})
}

func BenchmarkSelectableWaitGroup1_Add(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		wg := internal.NewSelectableWaitGroup1()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
		}
	})
	b.Run("parallel", func(b *testing.B) {
		wg := internal.NewSelectableWaitGroup1()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wg.Add(1)
			}
		})
	})
}

func BenchmarkSelectableWaitGroup2_Add(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		wg := internal.NewSelectableWaitGroup2()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
		}
	})

	b.Run("parallel", func(b *testing.B) {
		wg := internal.NewSelectableWaitGroup2()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wg.Add(1)
			}
		})
	})
}

func BenchmarkSelectableWaitGroup3_Add(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		wg := internal.NewSelectableWaitGroup3()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
		}
	})

	b.Run("parallel", func(b *testing.B) {
		wg := internal.NewSelectableWaitGroup3()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wg.Add(1)
			}
		})
	})
}

func BenchmarkSelectableWaitGroup4_Add(b *testing.B) {
	b.Run("single", func(b *testing.B) {
		wg := gsync.NewSelectableWaitGroup()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
		}
	})

	b.Run("parallel", func(b *testing.B) {
		wg := gsync.NewSelectableWaitGroup()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wg.Add(1)
			}
		})
	})
}
