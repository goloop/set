package set

import (
	"context"
	"reflect"
	"runtime"
	"sync"
)

var (
	// The parallelTasks the number of parallel tasks.
	parallelTasks = 1

	// The maxParallelTasks is the maximum number of parallel tasks.
	maxParallelTasks = runtime.NumCPU() * 3

	// The minLoadPerGoroutine is the minimum slice size for processing
	// in an individual goroutine. Essentially, it delineates the threshold
	// at which it becomes worthwhile to divide the slice processing amongst
	// multiple goroutines. If each goroutine isn't handling a sufficiently
	// large subslice, the overhead of goroutine creation and management
	// may outweigh the benefits of concurrent processing. This variable
	// specifies the minimum number of iterations per goroutine to ensure
	// an efficient division of labor.
	minLoadPerGoroutine = 1024
)

// The init initializes the randomGenerator variable.
func init() {
	parallelTasks = runtime.NumCPU() * 2
}

// The logicFoundValue is a helper struct that holds a boolean value
// and a Mutex to protect it from concurrent access.
//
// They are used in the In function to detect the desired result
// in a separate goroutine.
type logicFoundValue struct {
	m     sync.Mutex
	value bool
	err   error
}

// SetValue sets a new value for the Found. It locks the Mutex before
// changing the value and unlocks it after the change is complete.
func (f *logicFoundValue) SetValue(value bool, err error) {
	f.m.Lock()
	defer f.m.Unlock()
	f.value = value
	f.err = err
}

// GetValue retrieves the current value of the Found. It locks the Mutex
// before reading the value and unlocks it after the read is complete.
func (f *logicFoundValue) GetValue() (bool, error) {
	f.m.Lock()
	defer f.m.Unlock()
	return f.value, f.err
}

// ParallelTasks returns the number of parallel tasks.
//
// If the function is called without parameters, it returns the
// current value of parallelTasks.
//
// A function can receive one or more values for parallelTasks,
// these values are added together to form the final result for
// parallelTasks. If the new value for parallelTasks is less than
// or equal to zero - it will be set to 1, if it is greater than
// maxParallelTasks - it will be set to maxParallelTasks.
func ParallelTasks(v ...int) int {
	if len(v) > 0 {
		n := 0
		for _, value := range v {
			n += value
		}

		if n <= 0 {
			parallelTasks = 1
		} else if n > maxParallelTasks {
			parallelTasks = maxParallelTasks
		} else {
			parallelTasks = n
		}
	}

	return parallelTasks
}

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
		heap:   make(map[uint64]T),
		simple: 0,
		ctx:    ctx,
	}
	set.IsSimple()    // cache the complexity of the object
	set.Add(items...) // add items to the set

	return &set
}

// AddWithContext adds the provided items to the set.
//
// The function takes a context as the first argument and
// can be interrupted externally.
func AddWithContext[T any](ctx context.Context, s *Set[T], items ...T) error {
	return s.addWithContext(ctx, items...)
}

// Add adds the provided items to the set.
//
// Example usage:
//
//	s := set.New[int]()
//	set.Add(s, 1, 2, 3) // s is 1, 2, 3
func Add[T any](s *Set[T], items ...T) {
	AddWithContext[T](nil, s, items...)
}

// DeleteWithContext deletes the provided items from the set.
//
// The function takes a context as the first argument and
// can be interrupted externally.
func DeleteWithContext[T any](
	ctx context.Context,
	s *Set[T],
	items ...T,
) error {
	return s.deleteWithContext(ctx, items...)
}

// Delete deletes the provided items from the set.
//
// Example usage:
//
//	s := set.New[int]()
//	set.Add(s, 1, 2, 3) // s is 1, 2, 3
//	set.Delete(s, 1, 2) // s is 3
func Delete[T any](s *Set[T], items ...T) {
	DeleteWithContext[T](nil, s, items...)
}

// ContainsWithContext returns true if the set contains all of the provided
// items, otherwise it returns false.
//
// The function takes a context as the first argument and
// can be interrupted externally.
func ContainsWithContext[T any](
	ctx context.Context,
	s *Set[T],
	item T,
) (bool, error) {
	return s.containsWithContext(ctx, item)
}

// Contains returns true if the set contains all of the provided
// items, otherwise it returns false.
//
// Example usage:
//
//	s := set.New[int]()
//	set.Add(s, 1, 2, 3) // s is 1, 2, 3
//	set.Contains(s, 1)  // returns true
//	set.Contains(s, 4)  // returns false
func Contains[T any](s *Set[T], item T) bool {
	return s.Contains(item)
}

// ElementsWithContext returns a slice of the elements of the set using the
// provided context.
//
// The function takes a context as the first argument and
// can be interrupted externally.
func ElementsWithContext[T any](ctx context.Context, s *Set[T]) ([]T, error) {
	return s.elementsWithContext(ctx)
}

// Elements returns a slice of the elements of the set.
func Elements[T any](s *Set[T]) []T {
	return s.Elements()
}

// SortedWithContext returns a slice of the sorted elements of the set
// using the provided context.
//
// The function takes a context as the first argument and
// can be interrupted externally.
func SortedWithContext[T any](ctx context.Context, s *Set[T],
	fns ...func(a, b T) bool) ([]T, error) {
	return s.sortedWithContext(ctx, fns...)
}

// Sorted returns a slice of the sorted elements of the set.
func Sorted[T any](s *Set[T], fns ...func(a, b T) bool) []T {
	return s.Sorted(fns...)
}

// FilteredWithContext returns a slice of the elements of the set that
// satisfy the provided filter function using the provided context.
//
// The function takes a context as the first argument and
// can be interrupted externally.
func FilteredWithContext[T any](
	ctx context.Context,
	s *Set[T],
	fn func(item T) bool,
) ([]T, error) {
	return s.filteredWithContext(ctx, fn)
}

// Filtered returns a slice of the elements of the set that
// satisfy the provided filter function.
func Filtered[T any](s *Set[T], fn func(item T) bool) []T {
	return s.Filtered(fn)
}

// Len returns the number of items in the set.
func Len[T any](s *Set[T]) int {
	return s.Len()
}

// UnionWithContext returns a new set with all the items that are in either
// the set or in the other set.
//
// The function takes a context as the first argument and
// can be interrupted externally.
func UnionWithContext[T any](
	ctx context.Context,
	s *Set[T],
	others ...*Set[T],
) (*Set[T], error) {
	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	// Create a new set and add all the items from the current set.
	result := New[T]()
	for _, v := range s.heap {
		if err := result.addWithContext(ctx, v); err != nil {
			return New[T](), err
		}
	}

	// Add all the items from the other sets.
	for _, other := range others {
		for _, v := range other.heap {
			if err := result.addWithContext(ctx, v); err != nil {
				return New[T](), err
			}
		}
	}

	return result, nil
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

// IntersectionWithContext returns a new set with all the items
// that are in both the set and in the other set.
//
// The function takes a context as the first argument and
// can be interrupted externally.
func IntersectionWithContext[T any](
	ctx context.Context,
	s *Set[T],
	others ...*Set[T],
) (*Set[T], error) {
	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	// Create a new set.
	result := New[T]()
	for _, v := range s.heap {
		found := true
		for _, other := range others {
			ok, err := other.containsWithContext(ctx, v)
			if !ok && err == nil {
				found = false
				break
			} else if err != nil {
				return New[T](), err
			}
		}

		// If the item is in all the other sets, add it to the result.
		if found {
			if err := result.addWithContext(ctx, v); err != nil {
				return New[T](), err
			}
		}
	}

	return result, nil
}

// InterWithContext is a shortcut for IntersectionWithContext.
func InterWithContext[T any](
	ctx context.Context,
	s *Set[T],
	others ...*Set[T],
) (*Set[T], error) {
	return IntersectionWithContext(ctx, s, others...)
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

// DifferenceWithContext returns a new set with all the items that are in
// the set but not in the other set.
func DifferenceWithContext[T any](
	ctx context.Context,
	s *Set[T],
	others ...*Set[T],
) (*Set[T], error) {

	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	// Create a new set and add all the items from the current set.
	result := New[T]()
	for _, v := range s.heap {
		if err := result.addWithContext(ctx, v); err != nil {
			return New[T](), err
		}
	}

	// Remove all the items from the other sets.
	for _, other := range others {
		for _, v := range other.heap {
			ok, err := result.containsWithContext(ctx, v)
			if ok && err == nil {
				result.Delete(v)
			} else if err != nil {
				return New[T](), err
			}
		}
	}

	return result, nil
}

// DiffWithContext is an alias for DifferenceWithContext.
func DiffWithContext[T any](
	ctx context.Context,
	s *Set[T],
	others ...*Set[T],
) (*Set[T], error) {
	return DifferenceWithContext(ctx, s, others...)
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

// SymmetricDifferenceWithContext returns a new set with all the items that
// are in the set or in the other set but not in both.
func SymmetricDifferenceWithContext[T any](
	ctx context.Context,
	s *Set[T],
	others ...*Set[T],
) (*Set[T], error) {
	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	// Add all the items from the set.
	result := New[T]()
	for _, v := range s.heap {
		if err := result.addWithContext(ctx, v); err != nil {
			return New[T](), err
		}
	}

	// Fiilter out the items that are in both sets.
	runtime.Gosched()
	for _, other := range others {
		for _, v := range other.heap {
			ok, err := result.containsWithContext(ctx, v)
			if ok && err == nil {
				result.Delete(v)
			} else if !ok && err == nil {
				result.Add(v)
			} else if err != nil {
				return New[T](), err
			}
		}
	}

	return result, nil
}

// SdiffWithContext is an alias for SymmetricDifferenceWithContext.
func SdiffWithContext[T any](
	ctx context.Context,
	s *Set[T],
	others ...*Set[T],
) (*Set[T], error) {
	return SymmetricDifferenceWithContext(ctx, s, others...)
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

// MapWithContext returns a new set with the results of applying the provided
// function to each item in the set. The function is passed a context.Context
// as the first argument.
func MapWithContext[T any, R any](
	ctx context.Context,
	s *Set[T],
	fn func(item T) R,
) (*Set[R], error) {
	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	// Add all the items from the set.
	result := New[R]()
	for _, v := range s.heap {
		if err := result.addWithContext(ctx, fn(v)); err != nil {
			return New[R](), err
		}
	}

	return result, nil
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

// ReduceWithContext returns a single value by applying the provided function
// to each item in the set and passing the result of previous function call as
// the first argument in the next call.
//
// The function is passed a context.Context as the first argument.
func ReduceWithContext[T any, R any](
	ctx context.Context,
	s *Set[T],
	fn func(acc R, item T) R,
) (R, error) {
	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	var acc R
	for _, v := range s.heap {
		acc = fn(acc, v)
		select {
		case <-ctx.Done():
			z := reflect.Zero(reflect.TypeOf((*R)(nil)).Elem()).Interface().(R)
			return z, ctx.Err()
		default:
		}
	}

	return acc, nil
}

// Reduce returns a single value by applying the provided function to each
// item in the set and passing the result of previous function call as the
// first argument in the next call.
//
// Example usage:
//
//		type User struct {
//				Name string
//				Age  int
//		}
//
//	 s := set.New[User]()
//	 s.Add(User{"John", 20}, User{"Jane", 30})
//
//	 sum := sort.Reduce(s, func(acc int, item User) int {
//	     return acc + item.Age
//	 }) // sum is 50
func Reduce[T any, R any](s *Set[T], fn func(acc R, item T) R) R {
	r, _ := ReduceWithContext[T, R](nil, s, fn)
	return r
}

// CopyWithContext returns a new set with all the items from the set.
// The function is passed a context.Context as the first argument.
func CopyWithContext[T any](
	ctx context.Context,
	s *Set[T],
) (*Set[T], error) {
	return s.copyWithContext(ctx)
}

// Copy returns a new set with all the items from the set.
//
// Example usage:
//
//	s1 := set.New[int](1, 2, 3)
//	s2 := set.Copy(s1)
//	fmt.Println(s2.Sorted()) // 1, 2, 3
func Copy[T any](s *Set[T]) *Set[T] {
	r, _ := CopyWithContext[T](nil, s)
	return r
}

// FilterWithContext returns a new set with all the items from the set that
// pass the test implemented by the provided function.
// The function is passed a context.Context as the first argument.
func FilterWithContext[T any](
	ctx context.Context,
	s *Set[T],
	fn func(item T) bool,
) (*Set[T], error) {
	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	// Add all the items from the set.
	result := New[T]()
	for _, v := range s.heap {
		if fn(v) {
			if err := result.addWithContext(ctx, v); err != nil {
				return New[T](), err
			}
		}
	}

	return result, nil
}

// Filter returns a new set with all the items from the set that pass the
// test implemented by the provided function.
//
// Example usage:
//
//	s := set.New[int](1, 2, 3, 4, 5)
//	r := set.Filter(s, func(item int) bool {
//	    return item%2 == 0
//	})
//	fmt.Println(r.Sorted()) // 2, 4
func Filter[T any](s *Set[T], fn func(item T) bool) *Set[T] {
	r, _ := FilterWithContext[T](nil, s, fn)
	return r
}
