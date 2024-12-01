package set

import (
	"fmt"
	"math/rand"
	"testing"
)

func generateRandomInts(n int) []int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = rand.Intn(1000000)
	}
	return result
}

func BenchmarkSetOperations(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		randomInts := generateRandomInts(size)

		b.Run(fmt.Sprintf("New/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				New(randomInts...)
			}
		})

		set := New(randomInts...)

		b.Run(fmt.Sprintf("Add/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := New[int]()
				s.Add(randomInts...)
			}
		})

		b.Run(fmt.Sprintf("Contains/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				set.Contains(randomInts[i%size])
			}
		})

		b.Run(fmt.Sprintf("Delete/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := New(randomInts...)
				s.Delete(randomInts[:size/2]...)
			}
		})

		b.Run(fmt.Sprintf("Union/size=%d", size), func(b *testing.B) {
			other := New(generateRandomInts(size)...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				set.Union(other)
			}
		})

		b.Run(fmt.Sprintf("Intersection/size=%d", size), func(b *testing.B) {
			other := New(generateRandomInts(size)...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				set.Intersection(other)
			}
		})

		b.Run(fmt.Sprintf("Difference/size=%d", size), func(b *testing.B) {
			other := New(generateRandomInts(size)...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				set.Difference(other)
			}
		})

		b.Run(fmt.Sprintf("SymmetricDifference/size=%d", size), func(b *testing.B) {
			other := New(generateRandomInts(size)...)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				set.SymmetricDifference(other)
			}
		})

		b.Run(fmt.Sprintf("Copy/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				set.Copy()
			}
		})

		b.Run(fmt.Sprintf("Clear/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := New(randomInts...)
				s.Clear()
			}
		})

		b.Run(fmt.Sprintf("MarshalJSON/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				set.MarshalJSON()
			}
		})

		marshaledData, _ := set.MarshalJSON()
		b.Run(fmt.Sprintf("UnmarshalJSON/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := New[int]()
				s.UnmarshalJSON(marshaledData)
			}
		})

		b.Run(fmt.Sprintf("Filter/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				set.Filter(func(item int) bool {
					return item%2 == 0
				})
			}
		})

		b.Run(fmt.Sprintf("Any/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				set.Any(func(item int) bool {
					return item > 500000
				})
			}
		})

		b.Run(fmt.Sprintf("All/size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				set.All(func(item int) bool {
					return item >= 0
				})
			}
		})
	}
}
