package set

import (
	"encoding/json"
	"fmt"
	"iter"
	"slices"
)

// Set is a generic, unordered collection of unique elements of a comparable
// type T. It is backed by a Go map, so an element's identity is exactly the
// language's own equality (==): two elements are "the same" if and only if
// they are == to each other. This means there is no hashing, no reflection
// and no possibility of silently losing elements to a hash collision; the
// runtime map decides uniqueness.
//
// The comparable constraint covers all the usual element types — the numeric
// kinds, string, bool, pointers, channels, interfaces, and any struct or
// array whose fields are themselves comparable. Slices, maps and functions
// are not comparable and therefore cannot be Set elements directly; when you
// need to deduplicate such values, derive a comparable key from them (for
// example a string, or a struct of comparable fields) and build a Set of that
// key.
//
// Set is not safe for concurrent use by multiple goroutines, exactly like the
// built-in map it is built upon. If a Set is shared across goroutines and at
// least one of them mutates it, the callers are responsible for
// synchronization, e.g. with a sync.Mutex or sync.RWMutex.
//
// The zero value of Set is not ready for use; create a Set with New or
// NewWithCapacity.
type Set[T comparable] struct {
	m map[T]struct{}
}

// New creates a new Set containing the given items. Duplicate items collapse
// into a single element, so New(1, 2, 2, 3) holds exactly 1, 2 and 3.
//
// Example usage:
//
//	empty := set.New[int]()        // an empty set of int
//	s := set.New(1, 2, 3, 4)       // a set of int with four elements
//	letters := set.New("a", "b")   // a set of string
func New[T comparable](items ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{}, len(items))}
	s.Add(items...)
	return s
}

// NewWithCapacity creates a new Set with room pre-allocated for at least
// capacity elements before the underlying map needs to grow. Any items passed
// are added after the allocation. Use it when you know roughly how many
// elements you are about to insert to avoid repeated map resizes.
//
// A negative capacity is treated as zero.
func NewWithCapacity[T comparable](capacity int, items ...T) *Set[T] {
	if capacity < 0 {
		capacity = 0
	}
	if capacity < len(items) {
		capacity = len(items)
	}

	s := &Set[T]{m: make(map[T]struct{}, capacity)}
	s.Add(items...)
	return s
}

// Add inserts the given items into the set. Items already present are left
// untouched, so Add is idempotent.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3, 4) // s is 1, 2, 3 and 4
func (s *Set[T]) Add(items ...T) {
	for _, v := range items {
		s.m[v] = struct{}{}
	}
}

// Delete removes the given items from the set. Items that are not present are
// ignored.
//
// Example usage:
//
//	s := set.New(1, 2, 3, 4)
//	s.Delete(1, 3) // s is 2 and 4
func (s *Set[T]) Delete(items ...T) {
	for _, v := range items {
		delete(s.m, v)
	}
}

// Clear removes all elements from the set, leaving it empty. The underlying
// capacity is released.
//
// Example usage:
//
//	s := set.New(1, 2, 3)
//	s.Clear() // s is now empty
func (s *Set[T]) Clear() {
	clear(s.m)
}

// Overwrite replaces the entire contents of the set with the given items, as
// if by Clear followed by Add.
//
// Example usage:
//
//	s := set.New(1, 2, 3)
//	s.Overwrite(5, 6, 7) // s is now 5, 6 and 7
func (s *Set[T]) Overwrite(items ...T) {
	clear(s.m)
	s.Add(items...)
}

// Append adds every element of each of the given sets into this set, mutating
// it in place. It is the in-place counterpart of Union.
//
// Example usage:
//
//	s1 := set.New(1, 2, 3)
//	s2 := set.New(4, 5, 6)
//	s1.Append(s2) // s1 is now 1, 2, 3, 4, 5 and 6
func (s *Set[T]) Append(others ...*Set[T]) {
	for _, other := range others {
		if other == nil {
			continue
		}
		for v := range other.m {
			s.m[v] = struct{}{}
		}
	}
}

// Contains reports whether the item is present in the set.
//
// Example usage:
//
//	s := set.New(1, 2, 3, 4)
//	s.Contains(1) // true
//	s.Contains(5) // false
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}

// ContainsAll reports whether every one of the given items is present in the
// set. It returns true for an empty argument list (vacuously true).
//
// Example usage:
//
//	s := set.New(1, 2, 3, 4)
//	s.ContainsAll(1, 2) // true
//	s.ContainsAll(1, 9) // false
func (s *Set[T]) ContainsAll(items ...T) bool {
	for _, v := range items {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// ContainsAny reports whether at least one of the given items is present in
// the set. It returns false for an empty argument list.
//
// Example usage:
//
//	s := set.New(1, 2, 3, 4)
//	s.ContainsAny(9, 2) // true
//	s.ContainsAny(8, 9) // false
func (s *Set[T]) ContainsAny(items ...T) bool {
	for _, v := range items {
		if _, ok := s.m[v]; ok {
			return true
		}
	}
	return false
}

// Len returns the number of elements in the set.
//
// Example usage:
//
//	s := set.New(1, 2, 3, 4)
//	s.Len() // 4
func (s *Set[T]) Len() int {
	return len(s.m)
}

// IsEmpty reports whether the set has no elements.
//
// Example usage:
//
//	set.New[int]().IsEmpty()  // true
//	set.New(1).IsEmpty()      // false
func (s *Set[T]) IsEmpty() bool {
	return len(s.m) == 0
}

// Elements returns a slice with all elements of the set. The order is not
// specified and may differ between calls; use Sorted when you need a stable
// order.
//
// Example usage:
//
//	s := set.New(1, 2, 3, 4)
//	e := s.Elements() // some permutation of 1, 2, 3, 4
func (s *Set[T]) Elements() []T {
	result := make([]T, 0, len(s.m))
	for v := range s.m {
		result = append(result, v)
	}
	return result
}

// Iter returns an iterator over the elements of the set for use with range.
// The iteration order is not specified. It is not safe to add to or delete
// from the set while iterating over it.
//
// Example usage:
//
//	s := set.New(1, 2, 3)
//	for v := range s.Iter() {
//	    fmt.Println(v)
//	}
func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
	}
}

// Sorted returns all elements of the set as a slice ordered by the given
// comparison function, which must follow the same contract as the standard
// library cmp.Compare: it returns a negative number when a should sort before
// b, a positive number when a should sort after b, and zero when they are
// equal. The sort is stable, although for a set the elements are unique so
// stability only matters when cmp treats distinct elements as equal.
//
// For element types whose natural ordering you want, use the package-level
// Sorted function instead, which requires no comparison function.
//
// Example usage:
//
//	s := set.New(3, 1, 2)
//	s.Sorted(func(a, b int) int { return a - b }) // 1, 2, 3
func (s *Set[T]) Sorted(cmp func(a, b T) int) []T {
	result := s.Elements()
	slices.SortStableFunc(result, cmp)
	return result
}

// Filtered returns a slice with the elements that satisfy the predicate fn.
// The order is not specified. Use Filter to obtain a new Set instead of a
// slice.
//
// Example usage:
//
//	s := set.New(1, 2, 3, 4, 5)
//	s.Filtered(func(v int) bool { return v > 3 }) // 4 and 5
func (s *Set[T]) Filtered(fn func(item T) bool) []T {
	result := make([]T, 0, len(s.m))
	for v := range s.m {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// Pop removes an arbitrary element from the set and returns it together with
// true. If the set is empty it returns the zero value of T and false. Because
// the set is unordered there is no guarantee which element is returned.
//
// Example usage:
//
//	s := set.New(1, 2, 3)
//	v, ok := s.Pop() // v is one of 1, 2, 3; ok is true
func (s *Set[T]) Pop() (T, bool) {
	for v := range s.m {
		delete(s.m, v)
		return v, true
	}

	var zero T
	return zero, false
}

// Copy returns a new set holding the same elements as this one. The two sets
// are independent: mutating one does not affect the other.
//
// Example usage:
//
//	s := set.New(1, 2, 3)
//	c := s.Copy() // c is an independent set 1, 2, 3
func (s *Set[T]) Copy() *Set[T] {
	result := &Set[T]{m: make(map[T]struct{}, len(s.m))}
	for v := range s.m {
		result.m[v] = struct{}{}
	}
	return result
}

// Union returns a new set with every element that is in this set or in any of
// the other sets. The receiver and the arguments are not modified.
//
// Example usage:
//
//	s1 := set.New(1, 2, 3)
//	s2 := set.New(3, 4, 5)
//	s1.Union(s2) // 1, 2, 3, 4, 5
func (s *Set[T]) Union(others ...*Set[T]) *Set[T] {
	result := s.Copy()
	result.Append(others...)
	return result
}

// Intersection returns a new set with the elements common to this set and
// every one of the other sets. With no arguments it returns a copy of this
// set.
//
// To keep the work proportional to the smallest input, the result is seeded
// from whichever set is smaller at each step.
//
// Example usage:
//
//	s1 := set.New(1, 2, 3)
//	s2 := set.New(3, 4, 5)
//	s1.Intersection(s2) // 3
func (s *Set[T]) Intersection(others ...*Set[T]) *Set[T] {
	result := s.Copy()
	for _, other := range others {
		if other == nil {
			result.Clear()
			break
		}

		// Iterate the smaller side, probe the larger one.
		small, large := result, other
		if large.Len() < small.Len() {
			small, large = large, small
		}

		next := &Set[T]{m: make(map[T]struct{}, small.Len())}
		for v := range small.m {
			if _, ok := large.m[v]; ok {
				next.m[v] = struct{}{}
			}
		}
		result = next
	}
	return result
}

// Inter is an alias for Intersection.
func (s *Set[T]) Inter(others ...*Set[T]) *Set[T] {
	return s.Intersection(others...)
}

// Difference returns a new set with the elements that are in this set but in
// none of the other sets.
//
// Example usage:
//
//	s1 := set.New(1, 2, 3)
//	s2 := set.New(3, 4, 5)
//	s1.Difference(s2) // 1, 2
func (s *Set[T]) Difference(others ...*Set[T]) *Set[T] {
	result := &Set[T]{m: make(map[T]struct{}, len(s.m))}
	for v := range s.m {
		inOther := false
		for _, other := range others {
			if other == nil {
				continue
			}
			if _, ok := other.m[v]; ok {
				inOther = true
				break
			}
		}
		if !inOther {
			result.m[v] = struct{}{}
		}
	}
	return result
}

// Diff is an alias for Difference.
func (s *Set[T]) Diff(others ...*Set[T]) *Set[T] {
	return s.Difference(others...)
}

// SymmetricDifference returns a new set with the elements that appear in an
// odd number of the input sets (this set together with the others). For two
// sets this is the classic symmetric difference: elements in exactly one of
// the two. For more sets it generalises to elements whose total number of
// memberships across all sets is odd.
//
// Example usage:
//
//	s1 := set.New(1, 2, 3)
//	s2 := set.New(3, 4, 5)
//	s1.SymmetricDifference(s2) // 1, 2, 4, 5
func (s *Set[T]) SymmetricDifference(others ...*Set[T]) *Set[T] {
	result := s.Copy()
	for _, other := range others {
		if other == nil {
			continue
		}
		for v := range other.m {
			if _, ok := result.m[v]; ok {
				delete(result.m, v)
			} else {
				result.m[v] = struct{}{}
			}
		}
	}
	return result
}

// Sdiff is an alias for SymmetricDifference.
func (s *Set[T]) Sdiff(others ...*Set[T]) *Set[T] {
	return s.SymmetricDifference(others...)
}

// Equal reports whether this set and the other set contain exactly the same
// elements. A nil other is treated as the empty set.
//
// Example usage:
//
//	set.New(1, 2, 3).Equal(set.New(3, 2, 1)) // true
//	set.New(1, 2).Equal(set.New(1, 2, 3))    // false
func (s *Set[T]) Equal(other *Set[T]) bool {
	if other == nil {
		return len(s.m) == 0
	}
	if len(s.m) != len(other.m) {
		return false
	}
	for v := range s.m {
		if _, ok := other.m[v]; !ok {
			return false
		}
	}
	return true
}

// IsSubset reports whether every element of this set is also in the other set
// (this ⊆ other). A set is a subset of itself, so equal sets are subsets of
// each other. A nil other is treated as the empty set.
//
// Example usage:
//
//	set.New(1, 2).IsSubset(set.New(1, 2, 3)) // true
//	set.New(1, 2, 3).IsSubset(set.New(1, 2)) // false
//	set.New(1, 2).IsSubset(set.New(1, 2))    // true
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	if other == nil {
		return len(s.m) == 0
	}
	if len(s.m) > len(other.m) {
		return false
	}
	for v := range s.m {
		if _, ok := other.m[v]; !ok {
			return false
		}
	}
	return true
}

// IsSub is an alias for IsSubset.
func (s *Set[T]) IsSub(other *Set[T]) bool {
	return s.IsSubset(other)
}

// IsProperSubset reports whether this set is a subset of the other set and
// the two are not equal (this ⊊ other). A nil other is treated as the empty
// set.
//
// Example usage:
//
//	set.New(1, 2).IsProperSubset(set.New(1, 2, 3)) // true
//	set.New(1, 2).IsProperSubset(set.New(1, 2))    // false
func (s *Set[T]) IsProperSubset(other *Set[T]) bool {
	if other == nil {
		return false
	}
	if len(s.m) >= len(other.m) {
		return false
	}
	return s.IsSubset(other)
}

// IsSuperset reports whether this set contains every element of the other set
// (this ⊇ other). A set is a superset of itself. A nil other is treated as
// the empty set, of which every set is a superset.
//
// Example usage:
//
//	set.New(1, 2, 3).IsSuperset(set.New(1, 2)) // true
//	set.New(1, 2).IsSuperset(set.New(1, 2, 3)) // false
//	set.New(1, 2).IsSuperset(set.New(1, 2))    // true
func (s *Set[T]) IsSuperset(other *Set[T]) bool {
	if other == nil {
		return true
	}
	if len(other.m) > len(s.m) {
		return false
	}
	for v := range other.m {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// IsSup is an alias for IsSuperset.
func (s *Set[T]) IsSup(other *Set[T]) bool {
	return s.IsSuperset(other)
}

// IsProperSuperset reports whether this set is a superset of the other set
// and the two are not equal (this ⊋ other). A nil other is treated as the
// empty set.
//
// Example usage:
//
//	set.New(1, 2, 3).IsProperSuperset(set.New(1, 2)) // true
//	set.New(1, 2).IsProperSuperset(set.New(1, 2))    // false
func (s *Set[T]) IsProperSuperset(other *Set[T]) bool {
	if other == nil {
		return len(s.m) > 0
	}
	if len(s.m) <= len(other.m) {
		return false
	}
	return s.IsSuperset(other)
}

// IsDisjoint reports whether this set and the other set share no elements. A
// nil other is treated as the empty set, which is disjoint from everything.
//
// Example usage:
//
//	set.New(1, 2).IsDisjoint(set.New(3, 4)) // true
//	set.New(1, 2).IsDisjoint(set.New(2, 3)) // false
func (s *Set[T]) IsDisjoint(other *Set[T]) bool {
	if other == nil {
		return true
	}

	// Probe the larger set with the smaller one's elements.
	small, large := s, other
	if large.Len() < small.Len() {
		small, large = large, small
	}
	for v := range small.m {
		if _, ok := large.m[v]; ok {
			return false
		}
	}
	return true
}

// Filter returns a new set with the elements that satisfy the predicate fn.
//
// Example usage:
//
//	s := set.New(1, 2, 3, 4, 5)
//	s.Filter(func(v int) bool { return v > 3 }) // 4, 5
func (s *Set[T]) Filter(fn func(item T) bool) *Set[T] {
	result := &Set[T]{m: make(map[T]struct{})}
	for v := range s.m {
		if fn(v) {
			result.m[v] = struct{}{}
		}
	}
	return result
}

// Map returns a new set with the result of applying fn to every element. The
// result type is the same as the element type; to change the type use the
// package-level Map function, since Go methods cannot introduce a new type
// parameter.
//
// Note that mapping can shrink the set: if fn maps two distinct elements to
// the same value, the result holds that value once.
//
// Example usage:
//
//	s := set.New(1, 2, 3)
//	s.Map(func(v int) int { return v * 2 }) // 2, 4, 6
func (s *Set[T]) Map(fn func(item T) T) *Set[T] {
	result := &Set[T]{m: make(map[T]struct{}, len(s.m))}
	for v := range s.m {
		result.m[fn(v)] = struct{}{}
	}
	return result
}

// Reduce combines all elements into a single value of the element type by
// repeatedly applying fn, starting from the zero value of T. The order in
// which elements are visited is not specified, so for a well-defined result
// fn should be associative and commutative (for example sum or max).
//
// When you need an accumulator of a different type, or an explicit initial
// value, use the package-level Reduce or Fold functions.
//
// Example usage:
//
//	s := set.New(1, 2, 3)
//	s.Reduce(func(acc, v int) int { return acc + v }) // 6
func (s *Set[T]) Reduce(fn func(acc, item T) T) T {
	var acc T
	for v := range s.m {
		acc = fn(acc, v)
	}
	return acc
}

// Any reports whether at least one element satisfies the predicate fn. It
// returns false for an empty set. Iteration stops as soon as a match is
// found.
//
// Example usage:
//
//	s := set.New(1, 2, 3)
//	s.Any(func(v int) bool { return v > 2 }) // true
func (s *Set[T]) Any(fn func(item T) bool) bool {
	for v := range s.m {
		if fn(v) {
			return true
		}
	}
	return false
}

// All reports whether every element satisfies the predicate fn. It returns
// true for an empty set (the condition holds vacuously). Iteration stops as
// soon as an element fails the predicate.
//
// Example usage:
//
//	s := set.New(2, 4, 6)
//	s.All(func(v int) bool { return v%2 == 0 }) // true
//	set.New[int]().All(func(int) bool { return false }) // true (empty)
func (s *Set[T]) All(fn func(item T) bool) bool {
	for v := range s.m {
		if !fn(v) {
			return false
		}
	}
	return true
}

// MarshalJSON implements the json.Marshaler interface. The set is encoded as a
// JSON array of its elements; the order is not specified.
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Elements())
}

// UnmarshalJSON implements the json.Unmarshaler interface. It decodes a JSON
// array and replaces the contents of the set with its elements, collapsing
// duplicates.
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var elements []T
	if err := json.Unmarshal(data, &elements); err != nil {
		return fmt.Errorf("set: failed to unmarshal elements: %w", err)
	}

	if s.m == nil {
		s.m = make(map[T]struct{}, len(elements))
	} else {
		clear(s.m)
	}
	s.Add(elements...)
	return nil
}
