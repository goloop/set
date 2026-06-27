package set

import (
	"reflect"
	"sort"
	"testing"
)

func TestFreeMapChangesType(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	s := New(User{"John", 20}, User{"Jane", 30})

	names := Map(s, func(u User) string { return u.Name })
	got := names.Elements()
	sort.Strings(got)
	if !reflect.DeepEqual(got, []string{"Jane", "John"}) {
		t.Fatalf("Map names = %v, want [Jane John]", got)
	}
}

// Free Map collapses collisions in the target type.
func TestFreeMapCollapses(t *testing.T) {
	s := New(1, 2, 3, 4, 5)
	parity := Map(s, func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	})
	if parity.Len() != 2 {
		t.Fatalf("Len = %d, want 2 (even/odd)", parity.Len())
	}
	if !parity.ContainsAll("even", "odd") {
		t.Fatal("parity set must contain both labels")
	}
}

func TestFreeReduceZeroStart(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	s := New(User{"John", 20}, User{"Jane", 30})
	total := Reduce(s, func(acc int, u User) int { return acc + u.Age })
	if total != 50 {
		t.Fatalf("Reduce total = %d, want 50", total)
	}
}

// BUG-08 regression: Fold supplies a real initial value, so a product works.
func TestFoldProduct(t *testing.T) {
	s := New(2, 3, 4)
	product := Fold(s, 1, func(acc, v int) int { return acc * v })
	if product != 24 {
		t.Fatalf("Fold product = %d, want 24", product)
	}

	// The same operation via Reduce (zero start) is wrong on purpose — it
	// demonstrates why Fold exists.
	bad := Reduce(s, func(acc, v int) int { return acc * v })
	if bad != 0 {
		t.Fatalf("Reduce product = %d, want 0 (zero start swallows product)", bad)
	}
}

func TestFoldEmpty(t *testing.T) {
	s := New[int]()
	if got := Fold(s, 42, func(acc, v int) int { return acc + v }); got != 42 {
		t.Fatalf("Fold on empty = %d, want initial 42", got)
	}
}

func TestFreeSortedNaturalOrder(t *testing.T) {
	ints := Sorted(New(3, 1, 2))
	if !reflect.DeepEqual(ints, []int{1, 2, 3}) {
		t.Fatalf("Sorted ints = %v", ints)
	}

	strs := Sorted(New("banana", "apple", "cherry"))
	if !reflect.DeepEqual(strs, []string{"apple", "banana", "cherry"}) {
		t.Fatalf("Sorted strings = %v", strs)
	}

	if got := Sorted(New[int]()); len(got) != 0 {
		t.Fatalf("Sorted empty = %v, want empty", got)
	}
}

// --- String-keyed dedup (typical real use) --------------------------------

func TestStringSet(t *testing.T) {
	s := New("a", "b", "a", "c", "b")
	if s.Len() != 3 {
		t.Fatalf("Len = %d, want 3", s.Len())
	}
	got := Sorted(s)
	if !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Fatalf("got %v", got)
	}
}
