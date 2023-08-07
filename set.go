// Package set provides a parameterized Set data structure
// for Go.
//
// A Set can contain any type of object, including
// both simple and complex types. However, it is important
// to note that a Set can only contain either simple or complex
// types, not both.
//
// This package provides basic set operations, such as Add, Delete,
// Contains, and Len. In addition, it also provides complex set
// operations, such as Union, Intersection, Difference, SymmetricDifference,
// IsSubset, and IsSuperset.
package set

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

// sortingElement is a helper struct that is used to sort the set.
type sortingElement[T any] struct {
	key   string
	value T
}

// Set is a set of any objects. The set can contain both simple and complex
// types. It is important to note that the set can only one specific type.
// This information is stored in the 'simple' field where -1 denotes complex
// objects, 0 denotes that the type hasn't been set yet, and 1 denotes simple
// objects. The actual elements are stored in a map called 'heap' where the
// keys are hashed string representations of the objects, and the values are
// the objects themselves.
type Set[T any] struct {
	heap   map[string]T // collection of objects
	simple int          // -1 - complex object, 0 - not set, 1 - simple object
	ctx    context.Context
}

// toHash converts the given object to a string. If the set contains simple
// objects, this function uses the built-in Sprintf function to create the
// string representation. If the set contains complex objects, this function
// uses the 'valueToString' function to create a string representation of the
// object. This function is mainly used as a helper function to create unique
// keys for the 'heap' map in the Set.
func (s *Set[T]) toHash(ctx context.Context, obj T) string {
	// I think there is no point in hashing the result string or doing
	// something like strip - it's just additional resources for string
	// conversion.
	if s.IsSimple() {
		return fmt.Sprintf("%v", obj)
	}

	return toStr(ctx, reflect.ValueOf(obj))
}

// toStr is a helper function that takes a reflect.Value and creates a
// string representation of it. This function uses a switch statement to
// handle different kinds of complex types like Struct, Array, Slice, Map,
// Ptr, Interface, and Func. For each kind, it recursively builds a string
// representation and joins them together. If the kind doesn't fall into one of
// these categories, it uses the built-in Sprintf function to create a string.
// This function is mainly used by 'toHash' function to create unique keys for
// complex objects in the Set.
func toStr(ctx context.Context, v reflect.Value) string {
	// If the context is nil, create a new one.
	if ctx == nil {
		ctx = context.Background()
	}

	// If the context is done, return an empty string.
	select {
	case <-ctx.Done():
		return ""
	default:
	}

	// Create a string representation of the given reflect.Value.
	// This procedure performs a recursive call toStr.
	switch v.Kind() {
	case reflect.Struct:
		var r []string
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			name := t.Field(i).Name
			value := toStr(ctx, v.Field(i))
			r = append(r, fmt.Sprintf("%s:%s", name, value))
		}
		return "{" + strings.Join(r, ", ") + "}"
	case reflect.Array, reflect.Slice:
		var elements []string
		for i := 0; i < v.Len(); i++ {
			elements = append(elements, toStr(ctx, v.Index(i)))
		}
		return "[" + strings.Join(elements, ", ") + "]"
	case reflect.Map:
		var r []string
		for _, k := range v.MapKeys() {
			v := v.MapIndex(k)
			r = append(r, fmt.Sprintf("%s:%s", toStr(ctx, k), toStr(ctx, v)))
		}
		return "{" + strings.Join(r, ", ") + "}"
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return "nil"
		}
		return toStr(ctx, v.Elem())
	case reflect.Func:
		if v.IsNil() {
			return "func:nil"
		}
		return v.Type().String() + " Value"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// IsSimple determines the complexity of the objects in the set, i.e.,
// whether the objects are simple or complex.
//
// This method sets the field 'simple' based on the type of the object.
// If the set contains simple types such as byte, chan, bool, string, rune,
// int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
// uintptr, float32, float64, complex64, or complex128, the 'simple'
// field is set to 1.
//
// If the set contains complex types such as struct, array, slice,
// map, func, etc., the 'simple' field is set to -1.
//
// This method is invoked upon the creation of a set, and the complexity
// information  is cached for efficient subsequent operations.
// It returns true if the objects in the set are simple, and false otherwise.
func (s *Set[T]) IsSimple() bool {
	// If the complexity of the object is already defined.
	if s.simple != 0 {
		return s.simple == 1
	}

	// Determine the complexity of the object.
	// All simple types like: byte, chan, bool, string, rune, int,
	// int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64,
	// uintptr, float32, float64, complex64, complex128.
	// So set s.simple = 1.
	//
	// Other types of data, such as struct, array, slice, map, func, etc. -
	// are complex types. So set s.simple = -1.
	s.simple = 1
	k := reflect.TypeOf(s.heap).Elem().Kind()
	if k != reflect.String && k >= reflect.Array && k <= reflect.Struct {
		s.simple = -1
	}

	return s.simple == 1
}

// IsComplex returns true if the objects in the set are complex,
// and false otherwise.
func (s *Set[T]) IsComplex() bool {
	return !s.IsSimple()
}

// addWithContext adds the given items to the set.
func (s *Set[T]) addWithContext(ctx context.Context, items ...T) error {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	// Add the items to the set.
	for _, v := range items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			s.heap[s.toHash(s.ctx, v)] = v
		}
	}

	return nil
}

// Add adds the given items to the set.
//
// Example usage:
//
//	// Define a new set.
//	s := set.New[int]()
//
//	// Add elements to the set.
//	s.Add(1, 2, 3, 4) // s is 1, 2, 3, and 4
func (s *Set[T]) Add(items ...T) {
	s.addWithContext(s.ctx, items...)
}

// deleteWithContext removes the given items from the set.
func (s *Set[T]) deleteWithContext(ctx context.Context, items ...T) error {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	// Remove the items from the set.
	for _, v := range items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			delete(s.heap, s.toHash(s.ctx, v))
		}
	}

	return nil
}

// Delete removes the given items from the set.
//
// Example usage:
//
//	// Define a new set and add some elements
//	s := set.New[int]()
//	s.Add(1, 2, 3, 4)
//
//	// Remove elements from the set
//	s.Delete(1, 3) // s is 2 and 4
func (s *Set[T]) Delete(items ...T) {
	s.deleteWithContext(s.ctx, items...)
}

// containsWithContext returns true if the set contains the given item.
func (s *Set[T]) containsWithContext(
	ctx context.Context,
	item T,
) (bool, error) {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	_, ok := s.heap[s.toHash(s.ctx, item)]
	return ok, nil
}

// Contains returns true if the set contains the given item.
//
// Example usage:
//
//	// Define a new set and add some elements.
//	s := set.New[int]()
//	s.Add(1, 2, 3, 4)
//
//	// Check if the set contains certain elements.
//	containsOne := s.Contains(1)  // returns true
//	containsFive := s.Contains(5) // returns false
func (s *Set[T]) Contains(item T) bool {
	r, _ := s.containsWithContext(s.ctx, item)
	return r
}

// elementsWithContext returns all items in the set.
func (s *Set[T]) elementsWithContext(ctx context.Context) ([]T, error) {
	var items []T

	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	// Select all items from the set.
	for _, v := range s.heap {
		select {
		case <-ctx.Done():
			return []T{}, ctx.Err()
		default:
			items = append(items, v)
		}
	}

	return items, nil
}

// Elements returns all items in the set.
// This is useful when you need to iterate over the set,
// or when you need to convert the set to a slice ([]T).
// Note that the order of items is not guaranteed.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3, 4)
//	elements := s.Elements()  // elements is []int{1, 2, 3, 4}
func (s *Set[T]) Elements() []T {
	r, _ := s.elementsWithContext(s.ctx)
	return r
}

// sortedWithContext returns a slice of the sorted elements of the set
// using the provided context.
func (s *Set[T]) sortedWithContext(
	ctx context.Context,
	fns ...func(a, b T) bool,
) ([]T, error) {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	// Create a temporary slice of sortMarker[T] to hold
	// the data and sort it.
	tmp := make([]sortingElement[T], 0, len(s.heap)) // here is the change
	for k, v := range s.heap {
		tmp = append(tmp, sortingElement[T]{key: k, value: v})
		select {
		case <-ctx.Done():
			return []T{}, ctx.Err()
		default:
		}
	}

	// Sort the temporary slice.
	runtime.Gosched()
	if len(fns) == 0 {
		sort.Slice(tmp, func(i, j int) bool {
			return tmp[i].key < tmp[j].key
		})
	} else {
		for _, fn := range fns {
			sort.Slice(tmp, func(i, j int) bool {
				return fn(tmp[i].value, tmp[j].value)
			})
		}
	}

	// Create a new slice of T and copy the values over.
	var result = make([]T, len(tmp))
	runtime.Gosched()
	for i, v := range tmp {
		result[i] = v.value
		select {
		case <-ctx.Done():
			return []T{}, ctx.Err()
		default:
		}
	}

	return result, nil
}

// Sorted returns a slice of the sorted elements of the set.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(3, 1, 2)
//
//	sorted := s.Sorted() // sorted contains 1, 2, 3
func (s *Set[T]) Sorted(fns ...func(a, b T) bool) []T {
	r, _ := s.sortedWithContext(s.ctx, fns...)
	return r
}

// filteredWithContext returns a slice of items that satisfy the
// provided predicate.
func (s *Set[T]) filteredWithContext(
	ctx context.Context,
	fn func(item T) bool,
) ([]T, error) {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	var result = make([]T, 0, len(s.heap))
	for _, v := range s.heap {
		if fn(v) {
			result = append(result, v)
		}

		select {
		case <-ctx.Done():
			return []T{}, ctx.Err()
		default:
		}
	}

	return result, nil
}

// Filtered returns slice of items that satisfy the provided predicate.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3, 4, 5)
//
//	filtered := s.Filtered(func(item int) bool {
//		return item > 3
//	}) // filtered contains 4, 5
func (s *Set[T]) Filtered(fn func(item T) bool) []T {
	r, _ := s.filteredWithContext(s.ctx, fn)
	return r
}

// Len returns the number of items in the set.
// This is useful when you need to know how many items are in the set.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3, 4)
//	length := s.Len()  // length is 4
func (s *Set[T]) Len() int {
	return len(s.heap)
}

// uniunWithContext returns a new set with all the items in both sets.
func (s *Set[T]) unionWithContext(
	ctx context.Context,
	set *Set[T],
) (*Set[T], error) {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	// Elements of the base set.
	e, err := s.elementsWithContext(ctx)
	if err != nil {
		return New[T](), err
	}
	result := New[T](e...)

	// Elements of the other set.
	e, err = set.elementsWithContext(ctx)
	if err != nil {
		return New[T](), err
	}
	result.Add(e...)

	return result, nil
}

// Union returns a new set with all the items in both sets.
// This is useful when you want to merge two sets into a new one.
// Note that the result set will not have any duplicate items, even
// if the input sets do.
//
// Example usage:
//
//	s1 := set.New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := set.New[int]()
//	s2.Add(3, 4, 5)
//
//	union := s1.Union(s2)  // union contains 1, 2, 3, 4, 5
func (s *Set[T]) Union(set *Set[T]) *Set[T] {
	r, _ := s.unionWithContext(s.ctx, set)
	return r
}

// intersectionWithContext returns a new set with items that exist
// only in both sets.
func (s *Set[T]) intersectionWithContext(
	ctx context.Context,
	set *Set[T],
) (*Set[T], error) {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	result := New[T]()
	for _, v := range s.heap {
		ok, err := set.containsWithContext(ctx, v)
		if ok {
			err = result.addWithContext(ctx, v)
		}

		if err != nil {
			return New[T](), err
		}
	}

	return result, nil
}

// Intersection returns a new set with items that exist only in both sets.
//
// Example usage:
//
//	s1 := set.New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := set.New[int]()
//	s2.Add(3, 4, 5)
//
//	intersection := s1.Intersection(s2)  // intersection contains 3
func (s *Set[T]) Intersection(set *Set[T]) *Set[T] {
	r, _ := s.intersectionWithContext(s.ctx, set)
	return r
}

// Inter is an alias for Intersection.
func (s *Set[T]) Inter(set *Set[T]) *Set[T] {
	return s.Intersection(set)
}

// differenceWithContext returns a new set with items in the first set but
// not in the second.
func (s *Set[T]) differenceWithContext(
	ctx context.Context,
	set *Set[T],
) (*Set[T], error) {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	result := New[T]()
	for _, v := range s.heap {
		if !set.Contains(v) {
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

// Difference returns a new set with items in the first set but
// not in the second. This is useful when you want to find items
// that are unique to the first set.
//
// Example usage:
//
//	s1 := set.New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := set.New[int]()
//	s2.Add(3, 4, 5)
//
//	difference := s1.Difference(s2)  // difference contains 1, 2
func (s *Set[T]) Difference(set *Set[T]) *Set[T] {
	r, _ := s.differenceWithContext(s.ctx, set)
	return r
}

// Diff is an alias for Difference.
func (s *Set[T]) Diff(set *Set[T]) *Set[T] {
	return s.Difference(set)
}

// symmetricDifferenceWithContext returns a new set with items in either
// the first or second set but not both.
func (s *Set[T]) symmetricDifferenceWithContext(
	ctx context.Context,
	set *Set[T],
) (*Set[T], error) {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	// Elements of the base set.
	result := New[T]()
	for _, v := range s.heap {
		if !set.Contains(v) {
			result.Add(v)
		}

		select {
		case <-ctx.Done():
			return New[T](), ctx.Err()
		default:
		}
	}

	// Elements of the other set.
	runtime.Gosched()
	for _, v := range set.heap {
		if !s.Contains(v) {
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

// SymmetricDifference returns a new set with items in either
// the first or second set but not both. This is useful when you want to find
// items that are unique to each set.
//
// Example usage:
//
//	s1 := set.New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := set.New[int]()
//	s2.Add(3, 4, 5)
//
//	symmetricDifference := s1.SymmetricDifference(s2)  // 1, 2, 4, 5
func (s *Set[T]) SymmetricDifference(set *Set[T]) *Set[T] {
	r, _ := s.symmetricDifferenceWithContext(s.ctx, set)
	return r
}

// Sdiff is an alias for SymmetricDifference.
func (s *Set[T]) Sdiff(set *Set[T]) *Set[T] {
	return s.SymmetricDifference(set)
}

// mapWithContext returns a new set with the results of applying the
// provided function to each item in the set using the provided context.
func (s *Set[T]) mapWithContext(
	ctx context.Context,
	fn func(item T) T,
) (*Set[T], error) {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	// Create a new set to store the results.
	result := New[T]()
	for _, v := range s.heap {
		result.Add(fn(v))
		select {
		case <-ctx.Done():
			return New[T](), ctx.Err()
		default:
		}
	}

	return result, nil
}

// Map returns a new set with the results of applying the provided function
// to each item in the set.
//
// The result can only be of the same type as the elements of the set.
// For more flexibility, pay attention to the set.Reduce function.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3)
//
//	mapped := s.Map(func(item int) int {
//		return item * 2
//	}) // mapped contains 2, 4, 6
//
// Due to the fact that methods in Go don't support generics to change
// the result type we have to use the set.Map function.
func (s *Set[T]) Map(fn func(item T) T) *Set[T] {
	r, _ := s.mapWithContext(s.ctx, fn)
	return r
}

// reduceWithContext returns a single value by applying the provided function
// to each item in the set and passing the result of previous function call
// as the first argument in the next call.
func (s *Set[T]) reduceWithContext(
	ctx context.Context,
	fn func(acc, item T) T,
) (T, error) {
	// If context is nil, create default context.
	if ctx == nil {
		ctx = context.Background()
	}

	// Calculate.
	var acc T
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

// Reduce returns a single value by applying the provided function to each
// item in the set and passing the result of previous function call as the
// first argument in the next call.
//
// The result can only be of the same type as the elements of the set.
// For more flexibility, pay attention to the set.Reduce function.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3)
//
//	sum := s.Reduce(func(acc, item int) int) T {
//		return acc + item
//	}) // sum is 6
func (s *Set[T]) Reduce(fn func(acc, item T) T) T {
	acc, _ := s.reduceWithContext(nil, fn)
	return acc
}

// isSubsetWithContext returns true if all items in the first
// set exist in the second.
func (s *Set[T]) isSubsetWithContext(
	ctx context.Context,
	set *Set[T],
) (bool, error) {
	// If context is nil, create default context.
	if ctx == nil {
		ctx = context.Background()
	}

	// Elements of the set.
	for _, v := range s.heap {
		if !set.Contains(v) {
			return false, nil
		}

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}
	}

	return true, nil
}

// IsSubset returns true if all items in the first set exist in the second.
// This is useful when you want to check if all items of one set
// belong to another set.
//
// Example usage:
//
//	s1 := set.New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := set.New[int]()
//	s2.Add(1, 2, 3, 4, 5)
//
//	isSubset := s1.IsSubset(s2)  // isSubset is true
func (s *Set[T]) IsSubset(set *Set[T]) bool {
	r, _ := s.isSubsetWithContext(s.ctx, set)
	return r
}

// IsSub is an alias for IsSubset.
func (s *Set[T]) IsSub(set *Set[T]) bool {
	return s.IsSubset(set)
}

// isSupersetWithContext returns true if all items in the second
// set exist in the first.
func (s *Set[T]) isSupersetWithContext(
	ctx context.Context,
	set *Set[T],
) (bool, error) {
	// If the context is nil, create a new default context.
	if ctx == nil {
		ctx = context.Background()
	}

	// Elements of the other set.
	for _, v := range set.heap {
		ok, err := s.containsWithContext(ctx, v)
		if !ok && err == nil {
			return false, nil
		} else if err != nil {
			return false, err
		}
	}

	return true, nil
}

// IsSuperset returns true if all items in the second set exist in the first.
// This is useful when you want to check if one set contains all items
// of another set.
//
// Example usage:
//
//	s1 := set.New[int]()
//	s1.Add(1, 2, 3, 4, 5)
//
//	s2 := set.New[int]()
//	s2.Add(1, 2, 3)
//
//	isSuperset := s1.IsSuperset(s2)  // isSuperset is true
func (s *Set[T]) IsSuperset(set *Set[T]) bool {
	r, _ := s.isSupersetWithContext(s.ctx, set)
	return r
}

// IsSuper is an alias for IsSuperset.
func (s *Set[T]) IsSuper(set *Set[T]) bool {
	return s.IsSuperset(set)
}

// copyWithContext returns a new set with a copy of items in the set
// using the provided context.
func (s *Set[T]) copyWithContext(ctx context.Context) (*Set[T], error) {
	if ctx == nil {
		ctx = context.Background()
	}

	result := New[T]()
	for _, v := range s.heap {
		result.Add(v)
		select {
		case <-ctx.Done():
			return New[T](), ctx.Err()
		default:
		}
	}

	return result, nil
}

// Copy returns a new set with a copy of items in the set.
// This is useful when you want to copy the set.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3)
//
//	copied := s.Copy() // copied contains 1, 2, 3
func (s *Set[T]) Copy() *Set[T] {
	r, _ := s.copyWithContext(s.ctx)
	return r
}

// Append adds all elements from the provided sets to the current set.
//
// Example usage:
//
//	s1 := set.New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := New[int]()
//	s2.Add(4, 5, 6)
//
//	s1.Append(s2)  // s1 now contains 1, 2, 3, 4, 5, 6
func (s *Set[T]) Append(sets ...*Set[T]) {
	for _, set := range sets {
		for _, v := range set.heap {
			s.Add(v)
		}
	}
}

// Extend is an alias for Append. It adds all elements from
// the provided sets to the current set.
//
// Example usage:
//
//	s1 := set.New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := set.New[int]()
//	s2.Add(4, 5, 6)
//
//	s1.Extend(s2)  // s1 now contains 1, 2, 3, 4, 5, 6
func (s *Set[T]) Extend(sets []*Set[T]) {
	s.Append(sets...)
}

// Clear removes all items from the set.
//
// Example usage:
//
//	s := New[int]()
//	s.Add(1, 2, 3)
//
//	s.Clear() // s is now empty
func (s *Set[T]) Clear() {
	s.heap = make(map[string]T)
}

// Overwrite removes all items from the set and adds the provided items.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3)
//	s.Elements() // returns []int{1, 2, 3}
//
//	s.Overwrite(5, 6, 7) // as s.Clear() and s.Add(5, 6, 7)
//	s.Elements() // returns []int{5, 6, 7}
func (s *Set[T]) Overwrite(items ...T) {
	s.Clear()
	s.Add(items...)
}

// Filter returns a new set with items that satisfy the provided predicate.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3, 4, 5)
//
//	filtered := s.Filter(func(item int) bool {
//		return item > 3
//	}) // filtered contains 4, 5
func (s *Set[T]) Filter(fn func(item T) bool) *Set[T] {
	result := New[T]()
	for _, v := range s.heap {
		if fn(v) {
			result.Add(v)
		}
	}

	return result
}

// Any returns true if any of the items in the set satisfy
// the provided predicate.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3)
//
//	any := s.Any(func(item int) bool {
//		return item > 2
//	}) // any is true
func (s *Set[T]) Any(fn func(item T) bool) bool {
	for _, v := range s.heap {
		if fn(v) {
			return true
		}
	}

	return false
}

// All returns true if all of the items in the set satisfy
// the provided predicate.
//
// Example usage:
//
//	s := set.New[int]()
//	s.Add(1, 2, 3)
//
//	all := s.All(func(item int) bool {
//		return item > 2
//	}) // all is false
func (s *Set[T]) All(fn func(item T) bool) bool {
	for _, v := range s.heap {
		if !fn(v) {
			return false
		}
	}

	return true
}
