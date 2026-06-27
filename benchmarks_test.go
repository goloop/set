package set

import (
	"strconv"
	"testing"
)

var sizes = []int{10, 100, 1000, 10000}

func benchName(prefix string, n int) string {
	return prefix + "-" + strconv.Itoa(n)
}

func seedInts(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	return out
}

func BenchmarkAdd(b *testing.B) {
	for _, n := range sizes {
		data := seedInts(n)
		b.Run(benchName("ints", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				s := NewWithCapacity[int](n)
				s.Add(data...)
			}
		})
	}
}

func BenchmarkContains(b *testing.B) {
	for _, n := range sizes {
		s := New(seedInts(n)...)
		b.Run(benchName("ints", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = s.Contains(i % n)
			}
		})
	}
}

func BenchmarkUnion(b *testing.B) {
	for _, n := range sizes {
		x := New(seedInts(n)...)
		y := New(seedInts(n)...)
		y.Add(seedInts(n)...) // overlap
		b.Run(benchName("ints", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = x.Union(y)
			}
		})
	}
}

func BenchmarkIntersection(b *testing.B) {
	for _, n := range sizes {
		x := New(seedInts(n)...)
		half := make([]int, n/2)
		for i := range half {
			half[i] = i
		}
		y := New(half...)
		b.Run(benchName("ints", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = x.Intersection(y)
			}
		})
	}
}

func BenchmarkSymmetricDifference(b *testing.B) {
	for _, n := range sizes {
		x := New(seedInts(n)...)
		shifted := make([]int, n)
		for i := range shifted {
			shifted[i] = i + n/2
		}
		y := New(shifted...)
		b.Run(benchName("ints", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = x.SymmetricDifference(y)
			}
		})
	}
}

func BenchmarkSortedNatural(b *testing.B) {
	for _, n := range sizes {
		s := New(seedInts(n)...)
		b.Run(benchName("ints", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = Sorted(s)
			}
		})
	}
}
