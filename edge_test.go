package set

import "testing"

// Aliases must delegate to their canonical methods.
func TestAliases(t *testing.T) {
	a := New(1, 2)
	b := New(1, 2, 3)

	if a.IsSub(b) != a.IsSubset(b) {
		t.Fatal("IsSub must mirror IsSubset")
	}
	if b.IsSup(a) != b.IsSuperset(a) {
		t.Fatal("IsSup must mirror IsSuperset")
	}
	if !a.Inter(b).Equal(a.Intersection(b)) {
		t.Fatal("Inter must mirror Intersection")
	}
	if !a.Diff(b).Equal(a.Difference(b)) {
		t.Fatal("Diff must mirror Difference")
	}
	if !a.Sdiff(b).Equal(a.SymmetricDifference(b)) {
		t.Fatal("Sdiff must mirror SymmetricDifference")
	}
}

// A nil set among the variadic arguments must be handled coherently:
//   - Intersection with nil yields the empty set (∩ ∅ = ∅).
//   - Difference and SymmetricDifference treat nil as the empty set (no-op).
func TestVariadicNilHandling(t *testing.T) {
	a := New(1, 2, 3)

	if got := a.Intersection(nil); !got.IsEmpty() {
		t.Fatalf("Intersection(nil) = %v, want empty", got.Elements())
	}
	if got := a.Intersection(New(2, 3), nil); !got.IsEmpty() {
		t.Fatalf("Intersection(.., nil) = %v, want empty", got.Elements())
	}

	if got := a.Difference(nil); !got.Equal(a) {
		t.Fatalf("Difference(nil) must be a copy, got %v", got.Elements())
	}

	if got := a.SymmetricDifference(nil); !got.Equal(a) {
		t.Fatalf("SymmetricDifference(nil) must be a copy, got %v", got.Elements())
	}

	// Union already tolerates nil via Append.
	if got := a.Union(nil); !got.Equal(a) {
		t.Fatalf("Union(nil) must equal a, got %v", got.Elements())
	}
}

// Proper subset/superset against the empty set and unrelated sets.
func TestProperRelationsEdges(t *testing.T) {
	empty := New[int]()
	a := New(1, 2)

	// ∅ is a proper subset of any non-empty set.
	if !empty.IsProperSubset(a) {
		t.Fatal("∅ ⊊ a must be true for non-empty a")
	}
	// ∅ is not a proper subset of ∅.
	if empty.IsProperSubset(empty) {
		t.Fatal("∅ ⊊ ∅ must be false")
	}
	// A set unrelated by containment is neither proper subset nor superset.
	b := New(2, 3)
	if a.IsProperSubset(b) || a.IsProperSuperset(b) {
		t.Fatal("unrelated sets must not be proper subset/superset")
	}
	// Larger-or-equal cannot be a proper subset.
	if a.IsProperSubset(New(1)) {
		t.Fatal("larger set cannot be a proper subset of a smaller one")
	}
	// No set is a proper subset of the empty set (nil other).
	if a.IsProperSubset(nil) || empty.IsProperSubset(nil) {
		t.Fatal("nothing is a proper subset of ∅")
	}
}
