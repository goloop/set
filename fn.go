package set

import (
	"context"
	"runtime"
)

// New is a constructor function that creates a new Set[T] instance.
// It accepts an arbitrary number of items of a generic type 'T' which
// can be either simple types (e.g., int, string, bool) or complex types
// (e.g., struct, slice).
//
// This function first creates a new, empty set. It then determines whether
// the Set is simple or complex based on the type of the first item, and
// caches this information for efficient subsequent operations. Finally,
// it adds the provided items to the Set.
//
// Note: All items must be of the same type. If different types are provided,
// the behavior is undefined.
//
// Example usage:
//
//	// Creating a set of simple type (int)
//	emptySet := set.New[int]()       // empty set of int
//	simpleSet := set.New(1, 2, 3, 4) // set of int
//
//	// Creating a set of complex type (struct).
//	type ComplexType struct {
//	    field1 int
//	    field2 string
//	}
//	complexSet := set.New(
//	    ComplexType{1, "one"},
//	    ComplexType{2, "two"},
//	)
//
//	// Adding an item to the set.
//	simpleSet.Add(5)
//	complexSet.Add(ComplexType{3, "three"})
//
//	// Checking if an item exists in the set.
//	existsSimple := simpleSet.Contains(3)                       // returns true
//	existsComplex := complexSet.Contains(ComplexType{2, "two"}) // returns true
//
//	// Getting the size of the set.
//	size := simpleSet.Len() // returns 5
func New[T any](items ...T) *Set[T] {
	return NewWithContext[T](nil, items...)
}

// NewWithContext is a constructor function that creates a new Set[T] instance.
// It accepts a context.Context as the first argument, followed by an arbitrary
// number of items of a generic type 'T' which can be either simple types
// (e.g., int, string, bool) or complex types (e.g., struct, slice).
func NewWithContext[T any](ctx context.Context, items ...T) *Set[T] {
	set := Set[T]{
		heap:   make(map[string]T),
		simple: 0,
		ctx:    ctx,
	}
	set.IsSimple()    // cache the complexity of the object
	set.Add(items...) // add items to the set

	return &set
}

// Map returns a new set with the results of applying the provided function
// to each item in the set.
//
// Example usage:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	s := set.New[User]()
//	s.Add(User{"John", 20}, User{"Jane", 30})
//
//	names := sort.Map(s, func(item User) string {
//	    return item.Name
//	})
//
//	fmt.Println(names.Elements()) // "Jane", "John"
func Map[T any, R any](s *Set[T], fn func(item T) R) *Set[R] {
	r, _ := MapWithContext[T, R](nil, s, fn)
	return r
}

// MapWithContext returns a new set with the results of applying the provided
// function to each item in the set. The function is passed a context.Context
// as the first argument.
func MapWithContext[T any, R any](ctx context.Context,
	s *Set[T], fn func(item T) R) (*Set[R], error) {
	result := New[R]()

	if ctx == nil {
		ctx = context.Background()
	}

	for _, v := range s.heap {
		result.Add(fn(v))
		select {
		case <-ctx.Done():
			return New[R](), ctx.Err()
		default:
		}
	}

	return result, nil
}

// Reduce returns a single value by applying the provided function to each
// item in the set and passing the result of previous function call as the
// first argument in the next call.
//
// Example usage:
//
//	type User struct {
//			Name string
//			Age  int
//	}
//
//	 s := set.New[User]()
//	 s.Add(User{"John", 20}, User{"Jane", 30})
//
//	 sum := sort.Reduce(s, func(acc int, item User) int {
//			return acc + item.Age
//	 }) // sum is 50
func Reduce[T any, R any](s *Set[T], fn func(acc R, item T) R) R {
	r, _ := ReduceWithContext[T, R](nil, s, fn)
	return r
}

// ReduceWithContext returns a single value by applying the provided function
// to each item in the set and passing the result of previous function call as
// the first argument in the next call.
//
// The function is passed a context.Context as the first argument.
func ReduceWithContext[T any, R any](ctx context.Context,
	s *Set[T], fn func(acc R, item T) R) (R, error) {
	var acc R

	if ctx == nil {
		ctx = context.Background()
	}

	for _, v := range s.heap {
		acc = fn(acc, v)
		select {
		case <-ctx.Done():
			return acc, ctx.Err()
		default:
		}
	}

	return acc, nil
}

// Union returns a new set with all the items that are in either the set
// or in the other set.
//
// Example usage:
//
//	s1 := set.New[int](1, 2, 3)
//	s2 := set.New[int](3, 4, 5)
//	s3 := set.New[int](5, 6, 7)
//	s4 := set.New[int](7, 8, 9)
//
//	r := set.Union(s1, s2, s3, s4)
//	fmt.Println(r.Sorted()) // 1, 2, 3, 4, 5, 6, 7, 8, 9
func Union[T any](s *Set[T], others ...*Set[T]) *Set[T] {
	r, _ := UnionWithContext[T](nil, s, others...)
	return r
}

// UnionWithContext returns a new set with all the items that are in either
// the set or in the other set.
func UnionWithContext[T any](ctx context.Context,
	s *Set[T], others ...*Set[T]) (*Set[T], error) {
	result := New[T]()

	if ctx == nil {
		ctx = context.Background()
	}

	for _, v := range s.heap {
		result.Add(v)
	}

	for _, other := range others {
		for _, v := range other.heap {
			result.Add(v)

			select {
			case <-ctx.Done():
				return New[T](), ctx.Err()
			default:
			}
		}
	}

	return result, nil
}

// Intersection returns a new set with all the items that are in both the
// set and in the other set.
//
// Example usage:
//
//	s1 := set.New[int](1, 2, 3)
//	s2 := set.New[int](3, 4, 5)
//	s3 := set.New[int](5, 6, 7)
//	s4 := set.New[int](7, 8, 9)
//
//	r := set.Intersection(s1, s2, s3, s4)
//	fmt.Println(r.Sorted()) // 7
func Intersection[T any](s *Set[T], others ...*Set[T]) *Set[T] {
	r, _ := IntersectionWithContext[T](nil, s, others...)
	return r
}

// Inter is a shortcut for Intersection.
func Inter[T any](s *Set[T], others ...*Set[T]) *Set[T] {
	return Intersection(s, others...)
}

// IntersectionWithContext returns a new set with all the items
// that are in both the set and in the other set.
func IntersectionWithContext[T any](ctx context.Context,
	s *Set[T], others ...*Set[T]) (*Set[T], error) {
	result := New[T]()

	if ctx == nil {
		ctx = context.Background()
	}

	for _, v := range s.heap {
		isInAll := true
		for _, other := range others {
			if !other.Contains(v) {
				isInAll = false
				break
			}
		}

		if isInAll {
			result.Add(v)
		}

		select {
		case <-ctx.Done():
			return New[T](), ctx.Err()
		default:
		}
	}

	return result, nil
}

// InterWithContext is a shortcut for IntersectionWithContext.
func InterWithContext[T any](ctx context.Context,
	s *Set[T], others ...*Set[T]) (*Set[T], error) {
	return IntersectionWithContext(ctx, s, others...)
}

// Difference returns a new set with all the items that are in the set but
// not in the other set.
//
// Example usage:
//
//	s1 := set.New[int](1, 2, 3)
//	s2 := set.New[int](3, 4, 5)
//	s3 := set.New[int](5, 6, 7)
//	s4 := set.New[int](7, 8, 9)
//
//	r := set.Difference(s1, s2, s3, s4)
//	fmt.Println(r.Sorted()) // 1, 2
func Difference[T any](s *Set[T], others ...*Set[T]) *Set[T] {
	r, _ := DifferenceWithContext[T](nil, s, others...)
	return r
}

// Diff is an alias for Difference.
func Diff[T any](s *Set[T], others ...*Set[T]) *Set[T] {
	return Difference(s, others...)
}

// DifferenceWithContext returns a new set with all the items that are in
// the set but not in the other set.
func DifferenceWithContext[T any](ctx context.Context,
	s *Set[T], others ...*Set[T]) (*Set[T], error) {
	result := New[T]()

	if ctx == nil {
		ctx = context.Background()
	}

	for _, v := range s.heap {
		result.Add(v)
	}

	for _, other := range others {
		for _, v := range other.heap {
			if result.Contains(v) {
				result.Delete(v)
			}

			select {
			case <-ctx.Done():
				return New[T](), ctx.Err()
			default:
			}
		}
	}

	return result, nil
}

// DiffWithContext is an alias for DifferenceWithContext.
func DiffWithContext[T any](ctx context.Context,
	s *Set[T], others ...*Set[T]) (*Set[T], error) {
	return DifferenceWithContext(ctx, s, others...)
}

// SymmetricDifference returns a new set with all the items that are in the
// set or in the other set but not in both.
//
// Example usage:
//
//	s1 := set.New[int](1, 2, 3)
//	s2 := set.New[int](3, 4, 5)
//	s3 := set.New[int](5, 6, 7)
//	s4 := set.New[int](7, 8, 9)
//
//	r := set.SymmetricDifference(s1, s2, s3, s4)
//	fmt.Println(r.Sorted()) // 1, 2, 4, 6, 8, 9
func SymmetricDifference[T any](s *Set[T], others ...*Set[T]) *Set[T] {
	r, _ := SymmetricDifferenceWithContext[T](nil, s, others...)
	return r
}

// Sdiff is an alias for SymmetricDifference.
func Sdiff[T any](s *Set[T], others ...*Set[T]) *Set[T] {
	return SymmetricDifference(s, others...)
}

// SymmetricDifferenceWithContext returns a new set with all the items that
// are in the set or in the other set but not in both.
func SymmetricDifferenceWithContext[T any](ctx context.Context,
	s *Set[T], others ...*Set[T]) (*Set[T], error) {
	result := New[T]()

	if ctx == nil {
		ctx = context.Background()
	}

	// Add all the items from the set.
	for _, v := range s.heap {
		result.Add(v)
	}

	// Fiilter out the items that are in both sets.
	runtime.Gosched()
	for _, other := range others {
		for _, v := range other.heap {
			if result.Contains(v) {
				result.Delete(v)
			} else {
				result.Add(v)
			}

			select {
			case <-ctx.Done():
				return New[T](), ctx.Err()
			default:
			}
		}
	}

	return result, nil
}

// SdiffWithContext is an alias for SymmetricDifferenceWithContext.
func SdiffWithContext[T any](ctx context.Context,
	s *Set[T], others ...*Set[T]) (*Set[T], error) {
	return SymmetricDifferenceWithContext(ctx, s, others...)
}

// Elements returns a slice of the elements of the set.
func Elements[T any](s *Set[T]) []T {
	return s.Elements()
}

// ElementsWithContext returns a slice of the elements of the set using the
// provided context.
func ElementsWithContext[T any](ctx context.Context, s *Set[T]) ([]T, error) {
	return s.elementsWithContext(ctx)
}

// Sorted returns a slice of the sorted elements of the set.
func Sorted[T any](s *Set[T], fns ...func(a, b T) bool) []T {
	return s.Sorted(fns...)
}

// SortedWithContext returns a slice of the sorted elements of the set
// using the provided context.
func SortedWithContext[T any](ctx context.Context, s *Set[T],
	fns ...func(a, b T) bool) ([]T, error) {
	return s.sortedWithContext(ctx, fns...)
}
