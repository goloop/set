package set

import (
	"reflect"
	"sort"
	"testing"
)

// asSorted returns the elements of s as an ascending []int, so tests can
// compare against a deterministic expectation despite the unordered set.
func asSortedInt(s *Set[int]) []int {
	out := s.Elements()
	sort.Ints(out)
	return out
}

func eqInts(t *testing.T, got, want []int) {
	t.Helper()
	if len(got) == 0 && len(want) == 0 {
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// --- Construction & basic mutation ---------------------------------------

func TestNewDeduplicates(t *testing.T) {
	s := New(1, 2, 2, 3, 3, 3)
	if s.Len() != 3 {
		t.Fatalf("Len = %d, want 3", s.Len())
	}
	eqInts(t, asSortedInt(s), []int{1, 2, 3})
}

func TestNewEmpty(t *testing.T) {
	s := New[int]()
	if !s.IsEmpty() || s.Len() != 0 {
		t.Fatalf("empty set reports Len=%d IsEmpty=%v", s.Len(), s.IsEmpty())
	}
	if got := s.Elements(); len(got) != 0 {
		t.Fatalf("Elements on empty = %v, want empty", got)
	}
}

func TestNewWithCapacity(t *testing.T) {
	// Negative capacity must not panic and must behave like zero.
	s := NewWithCapacity(-5, 1, 2, 3)
	eqInts(t, asSortedInt(s), []int{1, 2, 3})

	// Capacity smaller than the item count must still admit all items.
	s = NewWithCapacity(1, 1, 2, 3, 4)
	eqInts(t, asSortedInt(s), []int{1, 2, 3, 4})
}

func TestAddIsIdempotent(t *testing.T) {
	s := New[int]()
	s.Add(1)
	s.Add(1)
	s.Add(1, 1, 1)
	if s.Len() != 1 {
		t.Fatalf("Len = %d, want 1 after idempotent adds", s.Len())
	}
}

func TestDeleteMissingIsNoop(t *testing.T) {
	s := New(1, 2, 3)
	s.Delete(99, 100) // not present
	s.Delete(2)
	eqInts(t, asSortedInt(s), []int{1, 3})
}

func TestClearReleasesAll(t *testing.T) {
	s := New(1, 2, 3)
	s.Clear()
	if !s.IsEmpty() {
		t.Fatalf("Clear left %d elements", s.Len())
	}
	// Must remain usable after Clear.
	s.Add(7)
	eqInts(t, asSortedInt(s), []int{7})
}

func TestOverwrite(t *testing.T) {
	s := New(1, 2, 3)
	s.Overwrite(5, 6, 7)
	eqInts(t, asSortedInt(s), []int{5, 6, 7})

	s.Overwrite() // overwrite with nothing -> empty
	if !s.IsEmpty() {
		t.Fatalf("Overwrite() left %d elements", s.Len())
	}
}

func TestAppendIgnoresNil(t *testing.T) {
	s := New(1, 2)
	s.Append(nil, New(2, 3), nil, New(4))
	eqInts(t, asSortedInt(s), []int{1, 2, 3, 4})
}

// --- Membership -----------------------------------------------------------

func TestContainsVariants(t *testing.T) {
	s := New(1, 2, 3)

	if !s.Contains(1) || s.Contains(9) {
		t.Fatal("Contains basic mismatch")
	}

	// Vacuous truth: ContainsAll of nothing is true, ContainsAny of nothing
	// is false.
	if !s.ContainsAll() {
		t.Fatal("ContainsAll() with no args must be true")
	}
	if s.ContainsAny() {
		t.Fatal("ContainsAny() with no args must be false")
	}

	if !s.ContainsAll(1, 2) || s.ContainsAll(1, 9) {
		t.Fatal("ContainsAll mismatch")
	}
	if !s.ContainsAny(9, 2) || s.ContainsAny(8, 9) {
		t.Fatal("ContainsAny mismatch")
	}
}

// --- Pop ------------------------------------------------------------------

func TestPopDrainsSet(t *testing.T) {
	s := New(1, 2, 3)
	seen := New[int]()
	for {
		v, ok := s.Pop()
		if !ok {
			break
		}
		if seen.Contains(v) {
			t.Fatalf("Pop returned duplicate %d", v)
		}
		seen.Add(v)
	}
	if !s.IsEmpty() {
		t.Fatalf("set not empty after draining, Len=%d", s.Len())
	}
	eqInts(t, asSortedInt(seen), []int{1, 2, 3})
}

func TestPopEmpty(t *testing.T) {
	s := New[int]()
	v, ok := s.Pop()
	if ok || v != 0 {
		t.Fatalf("Pop on empty = (%d, %v), want (0, false)", v, ok)
	}
}

// --- Copy independence ----------------------------------------------------

func TestCopyIsIndependent(t *testing.T) {
	s := New(1, 2, 3)
	c := s.Copy()

	c.Add(4)
	s.Delete(1)

	eqInts(t, asSortedInt(s), []int{2, 3})
	eqInts(t, asSortedInt(c), []int{1, 2, 3, 4})
}

// --- Comparable element identity (the old BUG-02 territory) ---------------

// Two structs that the old reflection-hash collapsed into one because it
// concatenated field bytes without separators. With == identity they are
// correctly distinct.
func TestStructFieldSwapDistinct(t *testing.T) {
	type P struct{ A, B string }
	s := New(P{"ab", "c"}, P{"a", "bc"})
	if s.Len() != 2 {
		t.Fatalf("Len = %d, want 2: field-swapped structs must stay distinct", s.Len())
	}
	if !s.Contains(P{"ab", "c"}) || !s.Contains(P{"a", "bc"}) {
		t.Fatal("both distinct structs must be present")
	}
}

// Equal structs must collapse to one element.
func TestStructEqualCollapses(t *testing.T) {
	type P struct{ A, B string }
	s := New(P{"x", "y"}, P{"x", "y"})
	if s.Len() != 1 {
		t.Fatalf("Len = %d, want 1 for equal structs", s.Len())
	}
}

// Arrays are comparable; field-swapped arrays must stay distinct too.
func TestArrayElementsDistinct(t *testing.T) {
	s := New([2]int{1, 23}, [2]int{12, 3})
	if s.Len() != 2 {
		t.Fatalf("Len = %d, want 2 for distinct arrays", s.Len())
	}
}

// Pointer identity: two pointers to equal values are distinct elements,
// the same pointer is one element.
func TestPointerIdentity(t *testing.T) {
	a, b := 1, 1
	pa, pb := &a, &b
	s := New(pa, pb, pa)
	if s.Len() != 2 {
		t.Fatalf("Len = %d, want 2 (two distinct pointers)", s.Len())
	}
}

// --- Set algebra correctness ---------------------------------------------

func TestUnion(t *testing.T) {
	a := New(1, 2, 3)
	b := New(3, 4, 5)
	c := New(5, 6)

	got := a.Union(b, c)
	eqInts(t, asSortedInt(got), []int{1, 2, 3, 4, 5, 6})

	// Inputs unchanged.
	eqInts(t, asSortedInt(a), []int{1, 2, 3})
	eqInts(t, asSortedInt(b), []int{3, 4, 5})
}

func TestUnionNoArgsIsCopy(t *testing.T) {
	a := New(1, 2, 3)
	got := a.Union()
	if got == a {
		t.Fatal("Union() must return a new set, not the receiver")
	}
	eqInts(t, asSortedInt(got), []int{1, 2, 3})
}

func TestIntersection(t *testing.T) {
	a := New(1, 2, 3, 4)
	b := New(3, 4, 5, 6)
	c := New(4, 5, 6, 7)

	eqInts(t, asSortedInt(a.Intersection(b)), []int{3, 4})
	eqInts(t, asSortedInt(a.Intersection(b, c)), []int{4})
	eqInts(t, asSortedInt(a.Inter(b)), []int{3, 4})
}

func TestIntersectionDisjoint(t *testing.T) {
	a := New(1, 2, 3)
	b := New(4, 5, 6)
	if got := a.Intersection(b); !got.IsEmpty() {
		t.Fatalf("disjoint intersection = %v, want empty", asSortedInt(got))
	}
}

func TestIntersectionNoArgsIsCopy(t *testing.T) {
	a := New(1, 2, 3)
	got := a.Intersection()
	eqInts(t, asSortedInt(got), []int{1, 2, 3})
	if got == a {
		t.Fatal("Intersection() must return a new set")
	}
}

func TestDifference(t *testing.T) {
	a := New(1, 2, 3, 4)
	b := New(3, 4)
	c := New(4, 5)

	eqInts(t, asSortedInt(a.Difference(b)), []int{1, 2})
	eqInts(t, asSortedInt(a.Difference(b, c)), []int{1, 2})
	eqInts(t, asSortedInt(a.Diff(b)), []int{1, 2})

	// Difference with nothing is a copy.
	eqInts(t, asSortedInt(a.Difference()), []int{1, 2, 3, 4})
}

func TestSymmetricDifferenceTwoSets(t *testing.T) {
	a := New(1, 2, 3)
	b := New(3, 4, 5)
	eqInts(t, asSortedInt(a.SymmetricDifference(b)), []int{1, 2, 4, 5})
	eqInts(t, asSortedInt(a.Sdiff(b)), []int{1, 2, 4, 5})
}

// Multi-set symmetric difference is parity of membership across all sets.
func TestSymmetricDifferenceParity(t *testing.T) {
	// {1} appears 3 times -> odd -> present.
	got := New(1).SymmetricDifference(New(1), New(1))
	eqInts(t, asSortedInt(got), []int{1})

	// {1} appears 2 times -> even -> absent.
	got = New(1).SymmetricDifference(New(1))
	if !got.IsEmpty() {
		t.Fatalf("Sdiff parity (even) = %v, want empty", asSortedInt(got))
	}

	got = New(1, 2).SymmetricDifference(New(1, 3), New(1, 4))
	// 1: x3 odd, 2: x1 odd, 3: x1 odd, 4: x1 odd -> all present.
	eqInts(t, asSortedInt(got), []int{1, 2, 3, 4})
}

// --- Relations ------------------------------------------------------------

// BUG-03 regression: a set is a (non-proper) subset and superset of itself.
func TestSubsetSupersetEqualSets(t *testing.T) {
	a := New(1, 2, 3)
	b := New(1, 2, 3)

	if !a.IsSubset(b) {
		t.Fatal("equal set must be a subset (A ⊆ A)")
	}
	if !a.IsSuperset(b) {
		t.Fatal("equal set must be a superset (A ⊇ A)")
	}
	// Mirror methods must agree for equal sets.
	if a.IsSubset(b) != a.IsSuperset(b) {
		t.Fatal("IsSubset and IsSuperset disagree on equal sets")
	}
	// Proper variants must be false for equal sets.
	if a.IsProperSubset(b) || a.IsProperSuperset(b) {
		t.Fatal("equal sets must not be proper subset/superset")
	}
}

func TestProperSubsetSuperset(t *testing.T) {
	small := New(1, 2)
	big := New(1, 2, 3)

	if !small.IsProperSubset(big) {
		t.Fatal("small ⊊ big expected")
	}
	if !big.IsProperSuperset(small) {
		t.Fatal("big ⊋ small expected")
	}
	if small.IsProperSuperset(big) || big.IsProperSubset(small) {
		t.Fatal("reverse proper relations must be false")
	}
}

func TestSubsetOfLargerByCountButDifferent(t *testing.T) {
	// Same size but different content: neither subset nor superset.
	a := New(1, 2, 3)
	b := New(1, 2, 4)
	if a.IsSubset(b) || a.IsSuperset(b) {
		t.Fatal("equal-size, different content must not be subset/superset")
	}

	// A smaller set that is not contained must not be a subset.
	c := New(9, 2)
	if c.IsSubset(New(1, 2, 3)) {
		t.Fatal("non-contained smaller set must not be a subset")
	}
}

func TestEqual(t *testing.T) {
	if !New(1, 2, 3).Equal(New(3, 2, 1)) {
		t.Fatal("permuted equal sets must be Equal")
	}
	if New(1, 2).Equal(New(1, 2, 3)) {
		t.Fatal("different-size sets must not be Equal")
	}
	if New(1, 2, 3).Equal(New(1, 2, 4)) {
		t.Fatal("same-size different sets must not be Equal")
	}
	if !New[int]().Equal(New[int]()) {
		t.Fatal("two empty sets must be Equal")
	}
}

func TestIsDisjoint(t *testing.T) {
	if !New(1, 2).IsDisjoint(New(3, 4)) {
		t.Fatal("disjoint sets expected")
	}
	if New(1, 2).IsDisjoint(New(2, 3)) {
		t.Fatal("overlapping sets must not be disjoint")
	}
	if !New[int]().IsDisjoint(New(1, 2)) {
		t.Fatal("empty set is disjoint from everything")
	}
}

// nil-other handling across the relation methods.
func TestRelationsWithNil(t *testing.T) {
	empty := New[int]()
	nonEmpty := New(1)

	if !empty.Equal(nil) {
		t.Fatal("empty.Equal(nil) must be true")
	}
	if nonEmpty.Equal(nil) {
		t.Fatal("nonEmpty.Equal(nil) must be false")
	}
	if !empty.IsSubset(nil) {
		t.Fatal("empty ⊆ ∅ must be true")
	}
	if nonEmpty.IsSubset(nil) {
		t.Fatal("nonEmpty ⊆ ∅ must be false")
	}
	if !nonEmpty.IsSuperset(nil) {
		t.Fatal("every set ⊇ ∅")
	}
	if !nonEmpty.IsProperSuperset(nil) {
		t.Fatal("nonEmpty ⊋ ∅ must be true")
	}
	if empty.IsProperSuperset(nil) {
		t.Fatal("∅ ⊋ ∅ must be false")
	}
	if !nonEmpty.IsDisjoint(nil) {
		t.Fatal("everything is disjoint from ∅")
	}
}

// --- Functional -----------------------------------------------------------

func TestFilterAndFiltered(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	pred := func(v int) bool { return v > 3 }

	eqInts(t, asSortedInt(s.Filter(pred)), []int{4, 5})

	got := s.Filtered(pred)
	sort.Ints(got)
	eqInts(t, got, []int{4, 5})
}

func TestMapMethodSameType(t *testing.T) {
	s := New(1, 2, 3)
	eqInts(t, asSortedInt(s.Map(func(v int) int { return v * 2 })), []int{2, 4, 6})
}

// Map collapses collisions.
func TestMapCollapsesCollisions(t *testing.T) {
	s := New(1, 2, 3, 4)
	got := s.Map(func(v int) int { return v % 2 }) // -> {1, 0}
	eqInts(t, asSortedInt(got), []int{0, 1})
}

func TestReduceMethod(t *testing.T) {
	s := New(1, 2, 3, 4)
	if sum := s.Reduce(func(acc, v int) int { return acc + v }); sum != 10 {
		t.Fatalf("Reduce sum = %d, want 10", sum)
	}
}

// BUG-04 regression: All over an empty set is vacuously true; Any is false.
func TestAnyAllEmpty(t *testing.T) {
	empty := New[int]()
	if empty.Any(func(int) bool { return true }) {
		t.Fatal("Any on empty must be false")
	}
	if !empty.All(func(int) bool { return false }) {
		t.Fatal("All on empty must be true (vacuous truth)")
	}
}

func TestAnyAll(t *testing.T) {
	s := New(2, 4, 6)
	if !s.All(func(v int) bool { return v%2 == 0 }) {
		t.Fatal("All even expected true")
	}
	if s.All(func(v int) bool { return v > 4 }) {
		t.Fatal("All >4 expected false")
	}
	if !s.Any(func(v int) bool { return v > 4 }) {
		t.Fatal("Any >4 expected true")
	}
	if s.Any(func(v int) bool { return v > 100 }) {
		t.Fatal("Any >100 expected false")
	}
}

// --- Sorted ---------------------------------------------------------------

func TestSortedMethodComparator(t *testing.T) {
	s := New(3, 1, 4, 1, 5, 9, 2, 6)
	asc := s.Sorted(func(a, b int) int { return a - b })
	eqInts(t, asc, []int{1, 2, 3, 4, 5, 6, 9})

	desc := s.Sorted(func(a, b int) int { return b - a })
	eqInts(t, desc, []int{9, 6, 5, 4, 3, 2, 1})
}

func TestSortedIsDeterministic(t *testing.T) {
	s := New(5, 3, 8, 1, 9, 2, 7)
	cmp := func(a, b int) int { return a - b }
	first := s.Sorted(cmp)
	for i := 0; i < 50; i++ {
		if !reflect.DeepEqual(s.Sorted(cmp), first) {
			t.Fatal("Sorted not deterministic across calls")
		}
	}
}

// --- Iter -----------------------------------------------------------------

func TestIterVisitsAll(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	seen := New[int]()
	for v := range s.Iter() {
		seen.Add(v)
	}
	if !seen.Equal(s) {
		t.Fatalf("Iter visited %v, want %v", asSortedInt(seen), asSortedInt(s))
	}
}

func TestIterEarlyBreak(t *testing.T) {
	s := New(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	count := 0
	for range s.Iter() {
		count++
		if count == 3 {
			break
		}
	}
	if count != 3 {
		t.Fatalf("Iter early break visited %d, want 3", count)
	}
}
