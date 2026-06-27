// Package set provides a generic Set: an unordered collection of unique
// elements of a comparable type, built directly on Go's built-in map.
//
// # Identity by ==
//
// An element's identity is the language's own equality. Two elements are the
// same if and only if they compare equal with ==, and the runtime map decides
// uniqueness. There is no hashing, no reflection and no custom equality
// contract to get wrong, so a Set can never silently lose an element to a
// collision and Len always reflects the true number of distinct elements.
//
// The comparable constraint admits the numeric kinds, string, bool, pointers,
// channels, interfaces, and any struct or array whose fields are themselves
// comparable. Slices, maps and functions are not comparable and cannot be
// elements directly; derive a comparable key from such values (for example a
// string, or a struct of comparable fields) and build a Set of that key.
//
// # Concurrency
//
// A Set is not safe for concurrent use by multiple goroutines, exactly like
// the built-in map. If a Set is shared and at least one goroutine mutates it,
// guard access with your own synchronization, for example:
//
//	type SafeSet[T comparable] struct {
//	    mu sync.RWMutex
//	    s  *set.Set[T]
//	}
//
// Keeping the core unsynchronized avoids per-operation locking overhead and
// lets the caller, who knows the application's concurrency model, choose the
// right strategy.
//
// # Basic operations
//
//   - New, NewWithCapacity: create a set
//   - Add, Delete, Overwrite, Append: mutate a set
//   - Contains, ContainsAll, ContainsAny: membership tests
//   - Len, IsEmpty: size queries
//   - Clear: empty a set
//   - Pop: remove and return an arbitrary element
//
// # Set algebra
//
//   - Union: elements in any of the sets
//   - Intersection (Inter): elements common to all sets
//   - Difference (Diff): elements in the first set but no other
//   - SymmetricDifference (Sdiff): elements with an odd membership count
//
// # Relations
//
//   - Equal: same elements
//   - IsSubset (IsSub), IsProperSubset
//   - IsSuperset (IsSup), IsProperSuperset
//   - IsDisjoint: no shared elements
//
// # Iteration and ordering
//
//   - Elements: all elements as a slice (unordered)
//   - Iter: an iter.Seq[T] for use with range
//   - AddSeq, Collect: build a set from an iter.Seq[T]
//   - Sorted method: order by a comparison function
//   - Sorted function: natural order for cmp.Ordered element types
//
// The zero value of a Set is an empty, ready-to-use set; the first insertion
// allocates its backing map.
//
// # Functional operations
//
//   - Filter, Filtered: select elements by a predicate
//   - Map method / Map function: transform elements (the function form may
//     change the element type)
//   - Reduce method / Reduce, Fold functions: aggregate into a single value
//   - Any, All: test a predicate over the set
//
// Because the iteration order of a set is unspecified, reducing functions
// should be associative and commutative for a deterministic result.
//
// # JSON
//
// A Set encodes as a JSON array of its elements and decodes from one,
// collapsing duplicates, via the standard encoding/json interfaces
// MarshalJSON and UnmarshalJSON.
//
// Example usage:
//
//	s1 := set.New(1, 2, 3)
//	s2 := set.New(3, 4, 5)
//
//	s1.Add(4)
//	s1.Delete(1)
//	exists := s1.Contains(2)
//
//	union := s1.Union(s2)
//	common := s1.Intersection(s2)
//	only1 := s1.Difference(s2)
//
//	even := s1.Filter(func(v int) bool { return v%2 == 0 })
//	doubled := set.Map(s1, func(v int) int { return v * 2 })
package set
