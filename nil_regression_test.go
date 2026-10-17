package set

import "testing"

func TestNilReadOperations(t *testing.T) {
	var s *Set[int]
	if !s.IsDisjoint(New(1)) {
		t.Fatal("nil set intersects")
	}
	if Map(s, func(v int) int { return v }).Len() != 0 {
		t.Fatal("map of nil set")
	}
	if Reduce(s, func(a, b int) int { return a + b }) != 0 {
		t.Fatal("reduce nil")
	}
	if Fold(s, 42, func(a, b int) int { return a + b }) != 42 {
		t.Fatal("fold nil")
	}
}

var intersectionResult *Set[int]

func BenchmarkIntersectionUneven(b *testing.B) {
	large, small := New[int](), New[int]()
	for i := 0; i < 65536; i++ {
		large.Add(i)
	}
	for i := 0; i < 16; i++ {
		small.Add(i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intersectionResult = large.Intersection(small)
	}
}
