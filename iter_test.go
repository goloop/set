package set

import (
	"slices"
	"sort"
	"testing"
)

// The zero value of a Set must be usable: reads return empty results and the
// first mutation lazily allocates the backing map instead of panicking.
func TestZeroValueUsable(t *testing.T) {
	var s Set[int]

	// Reads on the zero value must not panic.
	if s.Len() != 0 || !s.IsEmpty() {
		t.Fatalf("zero value not empty: Len=%d", s.Len())
	}
	if s.Contains(1) {
		t.Fatal("zero value must not contain anything")
	}
	if got := s.Elements(); len(got) != 0 {
		t.Fatalf("zero value Elements = %v", got)
	}
	for range s.Iter() {
		t.Fatal("zero value Iter must yield nothing")
	}

	// First mutation must work via lazy initialization.
	s.Add(1, 2, 3)
	if s.Len() != 3 {
		t.Fatalf("after Add Len = %d, want 3", s.Len())
	}
}

func TestZeroValueMutatorsLazyInit(t *testing.T) {
	// Append into a zero value.
	var a Set[int]
	a.Append(New(1, 2))
	eqInts(t, asSortedInt(&a), []int{1, 2})

	// Overwrite a zero value.
	var b Set[int]
	b.Overwrite(7, 8)
	eqInts(t, asSortedInt(&b), []int{7, 8})

	// AddSeq into a zero value.
	var c Set[int]
	c.AddSeq(slices.Values([]int{4, 4, 5}))
	eqInts(t, asSortedInt(&c), []int{4, 5})

	// Empty mutations on a zero value must stay empty and not panic.
	var d Set[int]
	d.Add()
	d.Append(nil)
	d.Append(New[int]())
	if !d.IsEmpty() {
		t.Fatalf("zero value with empty mutations not empty: Len=%d", d.Len())
	}
}

func TestAddSeq(t *testing.T) {
	s := New(1)
	s.AddSeq(slices.Values([]int{2, 3, 3, 1}))
	eqInts(t, asSortedInt(s), []int{1, 2, 3})

	// Round-trip Iter -> AddSeq reconstructs the set.
	dst := New[int]()
	dst.AddSeq(s.Iter())
	if !dst.Equal(s) {
		t.Fatalf("Iter->AddSeq mismatch: %v vs %v", asSortedInt(dst), asSortedInt(s))
	}
}

func TestCollect(t *testing.T) {
	got := Collect(slices.Values([]int{3, 1, 2, 2, 3}))
	eqInts(t, asSortedInt(got), []int{1, 2, 3})

	// Collect over a set's own iterator is a faithful copy.
	src := New("a", "b", "c")
	cp := Collect(src.Iter())
	if !cp.Equal(src) {
		t.Fatal("Collect(src.Iter()) must equal src")
	}

	if got := Collect(slices.Values([]int{})); !got.IsEmpty() {
		t.Fatalf("Collect of empty = %v, want empty", got.Elements())
	}
}

// Sanity: AddSeq consuming maps.Keys-style iteration produces the key set.
func TestCollectFromMapKeys(t *testing.T) {
	m := map[string]int{"x": 1, "y": 2, "z": 3}
	keys := Collect(func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	})
	got := keys.Elements()
	sort.Strings(got)
	if len(got) != 3 || got[0] != "x" || got[1] != "y" || got[2] != "z" {
		t.Fatalf("collected keys = %v", got)
	}
}
