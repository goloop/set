package set

import (
	"encoding/json"
	"testing"
)

// splitInts deterministically derives two int slices from a byte string, so
// fuzzing explores many (A, B) set pairs. The high bit of each byte routes it
// to A or B; the low bits (mod 16) keep the value space small enough that
// overlaps between A and B actually happen.
func splitInts(data []byte) (a, b []int) {
	for _, c := range data {
		v := int(c % 16)
		if c&0x80 == 0 {
			a = append(a, v)
		} else {
			b = append(b, v)
		}
	}
	return a, b
}

// FuzzSetAlgebra checks the algebraic laws that must hold for any two sets,
// regardless of contents. A violation means an operation is wrong.
func FuzzSetAlgebra(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1, 2, 3})
	f.Add([]byte{0x01, 0x81, 0x02, 0x82, 0x03})
	f.Add([]byte{0xFF, 0x7F, 0x10, 0x90, 0x11, 0x91})

	f.Fuzz(func(t *testing.T, data []byte) {
		ai, bi := splitInts(data)
		a := New(ai...)
		b := New(bi...)

		union := a.Union(b)
		inter := a.Intersection(b)
		da := a.Difference(b) // a \ b
		db := b.Difference(a) // b \ a
		sdiff := a.SymmetricDifference(b)

		// Union is a superset of both operands and is commutative.
		if !union.IsSuperset(a) || !union.IsSuperset(b) {
			t.Fatalf("union is not a superset of its operands")
		}
		if !union.Equal(b.Union(a)) {
			t.Fatalf("union not commutative")
		}

		// Intersection is a subset of both and is commutative.
		if !inter.IsSubset(a) || !inter.IsSubset(b) {
			t.Fatalf("intersection is not a subset of its operands")
		}
		if !inter.Equal(b.Intersection(a)) {
			t.Fatalf("intersection not commutative")
		}

		// |A∪B| + |A∩B| == |A| + |B| (inclusion–exclusion).
		if union.Len()+inter.Len() != a.Len()+b.Len() {
			t.Fatalf("inclusion-exclusion violated: |U|=%d |I|=%d |A|=%d |B|=%d",
				union.Len(), inter.Len(), a.Len(), b.Len())
		}

		// A\B is a subset of A and disjoint from B.
		if !da.IsSubset(a) || !da.IsDisjoint(b) {
			t.Fatalf("difference law violated")
		}

		// Symmetric difference equals (A\B) ∪ (B\A) and also equals
		// (A∪B) \ (A∩B).
		if !sdiff.Equal(da.Union(db)) {
			t.Fatalf("sdiff != (A\\B) ∪ (B\\A)")
		}
		if !sdiff.Equal(union.Difference(inter)) {
			t.Fatalf("sdiff != (A∪B) \\ (A∩B)")
		}

		// (A\B) and (B\A) are disjoint.
		if !da.IsDisjoint(db) {
			t.Fatalf("A\\B and B\\A must be disjoint")
		}

		// Equal is reflexive and copies compare equal.
		if !a.Equal(a.Copy()) {
			t.Fatalf("a != copy of a")
		}

		// Subset/superset duality.
		if a.IsSubset(b) != b.IsSuperset(a) {
			t.Fatalf("IsSubset/IsSuperset duality violated")
		}

		// Disjoint is symmetric and consistent with empty intersection.
		if a.IsDisjoint(b) != b.IsDisjoint(a) {
			t.Fatalf("IsDisjoint not symmetric")
		}
		if a.IsDisjoint(b) != inter.IsEmpty() {
			t.Fatalf("IsDisjoint inconsistent with empty intersection")
		}
	})
}

// FuzzJSONRoundTrip checks that marshalling then unmarshalling any set yields
// an equal set.
func FuzzJSONRoundTrip(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1, 2, 3, 3, 2, 1})
	f.Add([]byte{0xFF, 0x00, 0x7F, 0x80})

	f.Fuzz(func(t *testing.T, data []byte) {
		ai, _ := splitInts(append(data, 0)) // route everything to A
		s := New(ai...)

		encoded, err := json.Marshal(s)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}

		var back Set[int]
		if err := json.Unmarshal(encoded, &back); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !back.Equal(s) {
			t.Fatalf("round-trip mismatch: got %v, want %v",
				back.Elements(), s.Elements())
		}
	})
}

// FuzzModelAgainstMap is a stateful, model-based test: it replays a random
// sequence of operations against both the Set under test and a reference
// map[int]struct{}, asserting after every step that the two agree. Any
// divergence in Add/Delete/Contains/Len/Pop semantics surfaces here.
func FuzzModelAgainstMap(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0x00, 1, 0x00, 2, 0x40, 3})
	f.Add([]byte{0x80, 5, 0xC0, 5, 0x00, 7, 0x80, 7})

	f.Fuzz(func(t *testing.T, ops []byte) {
		s := New[int]()
		ref := make(map[int]struct{})

		// Each step consumes two bytes: an opcode and a value. Values are kept
		// in a small range so adds and deletes actually collide.
		for i := 0; i+1 < len(ops); i += 2 {
			val := int(ops[i+1] % 24)
			switch ops[i] >> 6 { // top two bits -> one of four ops
			case 0: // add (also the bias, since it is the common case)
				s.Add(val)
				ref[val] = struct{}{}
			case 1: // delete
				s.Delete(val)
				delete(ref, val)
			case 2: // contains must agree
				_, want := ref[val]
				if s.Contains(val) != want {
					t.Fatalf("Contains(%d)=%v, ref=%v", val, s.Contains(val), want)
				}
			case 3: // pop must remove exactly one known element
				got, ok := s.Pop()
				if ok != (len(ref) > 0) {
					t.Fatalf("Pop ok=%v but ref size=%d", ok, len(ref))
				}
				if ok {
					if _, present := ref[got]; !present {
						t.Fatalf("Pop returned %d not in reference", got)
					}
					delete(ref, got)
				}
			}

			if s.Len() != len(ref) {
				t.Fatalf("Len=%d, ref=%d after op %d", s.Len(), len(ref), ops[i]>>6)
			}
		}

		// Final state must match the reference exactly, both directions.
		if s.Len() != len(ref) {
			t.Fatalf("final Len=%d, ref=%d", s.Len(), len(ref))
		}
		for v := range ref {
			if !s.Contains(v) {
				t.Fatalf("reference has %d, set does not", v)
			}
		}
		for _, v := range s.Elements() {
			if _, ok := ref[v]; !ok {
				t.Fatalf("set has %d, reference does not", v)
			}
		}
	})
}

// FuzzSubsetConsistency cross-checks the relation predicates against a
// brute-force definition built only from Contains, so a bug in the optimized
// size-based shortcuts would surface.
func FuzzSubsetConsistency(f *testing.F) {
	f.Add([]byte{1, 2, 0x81, 0x82})
	f.Add([]byte{0x01, 0x02, 0x03})

	f.Fuzz(func(t *testing.T, data []byte) {
		ai, bi := splitInts(data)
		a := New(ai...)
		b := New(bi...)

		// Brute-force subset: every element of a is in b.
		bruteSubset := true
		for v := range a.m {
			if !b.Contains(v) {
				bruteSubset = false
				break
			}
		}
		if a.IsSubset(b) != bruteSubset {
			t.Fatalf("IsSubset=%v but brute-force=%v", a.IsSubset(b), bruteSubset)
		}

		// Proper subset is subset and not equal.
		wantProper := bruteSubset && !a.Equal(b)
		if a.IsProperSubset(b) != wantProper {
			t.Fatalf("IsProperSubset=%v, want %v", a.IsProperSubset(b), wantProper)
		}
	})
}
