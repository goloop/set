package set

import (
	"cmp"
	"iter"
	"slices"
)

// Collect builds a new set from all values produced by the iterator seq,
// collapsing duplicates. It mirrors slices.Collect and maps.Collect, and pairs
// with the Set.Iter method.
//
// Example usage:
//
//	m := map[string]int{"a": 1, "b": 2}
//	keys := set.Collect(maps.Keys(m)) // a set of "a" and "b"
func Collect[T comparable](seq iter.Seq[T]) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{})}
	for v := range seq {
		s.m[v] = struct{}{}
	}
	return s
}

// Map returns a new set whose elements are the results of applying fn to each
// element of s. Unlike the Map method, this package-level function can change
// the element type, because a Go method cannot introduce its own type
// parameter.
//
// Mapping can shrink the set: if fn maps two distinct elements of s to the
// same value, that value appears once in the result.
//
// Example usage:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	s := set.New(User{"John", 20}, User{"Jane", 30})
//	names := set.Map(s, func(u User) string { return u.Name })
//	// names holds "John" and "Jane"
func Map[T, R comparable](s *Set[T], fn func(item T) R) *Set[R] {
	result := &Set[R]{m: make(map[R]struct{}, len(s.m))}
	for v := range s.m {
		result.m[fn(v)] = struct{}{}
	}
	return result
}

// Reduce combines all elements of s into a single value of type R by
// repeatedly applying fn, starting from the zero value of R. The order in
// which elements are visited is not specified, so for a deterministic result
// fn should be associative and commutative.
//
// When you need a non-zero starting value (for example a product, which must
// start from 1, or a min/max), use Fold instead.
//
// Example usage:
//
//	s := set.New(User{"John", 20}, User{"Jane", 30})
//	total := set.Reduce(s, func(acc int, u User) int { return acc + u.Age })
//	// total is 50
func Reduce[T comparable, R any](s *Set[T], fn func(acc R, item T) R) R {
	var acc R
	for v := range s.m {
		acc = fn(acc, v)
	}
	return acc
}

// Fold is like Reduce but starts the accumulation from the given initial
// value instead of the zero value of R. This makes it suitable for operations
// that have no neutral zero, such as products or minima.
//
// As with Reduce, the iteration order is unspecified, so fn should be
// associative and commutative for a deterministic result.
//
// Example usage:
//
//	s := set.New(2, 3, 4)
//	product := set.Fold(s, 1, func(acc, v int) int { return acc * v })
//	// product is 24
func Fold[T comparable, R any](s *Set[T], initial R, fn func(acc R, item T) R) R {
	acc := initial
	for v := range s.m {
		acc = fn(acc, v)
	}
	return acc
}

// Sorted returns all elements of s as a slice in ascending natural order. It
// is available for element types that satisfy cmp.Ordered (the integer,
// floating-point and string kinds). For other comparable types, or for a
// custom order, use the Sorted method with a comparison function.
//
// Example usage:
//
//	s := set.New(3, 1, 2)
//	set.Sorted(s) // 1, 2, 3
func Sorted[T cmp.Ordered](s *Set[T]) []T {
	result := s.Elements()
	slices.Sort(result)
	return result
}
