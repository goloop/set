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
	"fmt"
	"reflect"
	"sort"
	"strings"
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
//	emptySet := New[int]()       // empty set of int
//	simpleSet := New(1, 2, 3, 4) // set of int
//
//	// Creating a set of complex type (struct).
//	type ComplexType struct {
//	    field1 int
//	    field2 string
//	}
//	complexSet := New(
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
	set := Set[T]{
		heap:   make(map[string]T),
		simple: 0,
	}
	set.isSimple()    // cache the complexity of the object
	set.Add(items...) // add items to the set

	return &set
}

// sortMarker is a helper struct that is used to sort the set.
type sortMarker[T any] struct {
	name  string
	value T
}

// Set is a set of any objects. The set can contain both simple and complex
// types. It is important to note that the set can only contain either simple
// or complex types, not both. This information is stored in the 'simple' field
// where -1 denotes complex objects, 0 denotes that the type hasn't been set
// yet, and 1 denotes simple objects. The actual elements are stored in a map
// called 'heap' where the keys are hashed string representations of the
// objects, and the values are the objects themselves.
type Set[T any] struct {
	heap   map[string]T // collection of objects
	simple int          // -1 - complex object, 0 - not set, 1 - simple object
}

// toHash converts the given object to a string. If the set contains simple
// objects, this function uses the built-in Sprintf function to create the
// string representation. If the set contains complex objects, this function
// uses the 'valueToString' function to create a string representation of the
// object. This function is mainly used as a helper function to create unique
// keys for the 'heap' map in the Set.
func (s *Set[T]) toHash(obj T) string {
	// I think there is no point in hashing the result string or doing
	// something like strip - it's just additional resources for string
	// conversion.
	if s.isSimple() {
		return fmt.Sprintf("%v", obj)
	}

	return toStr(reflect.ValueOf(obj))
}

// toStr is a helper function that takes a reflect.Value and creates a
// string representation of it. This function uses a switch statement to handle
// different kinds of complex types like Struct, Array, Slice, Map, Ptr,
// Interface, and Func. For each kind, it recursively builds a string
// representation and joins them together. If the kind doesn't fall into one of
// these categories, it uses the built-in Sprintf function to create a string.
// This function is mainly used by 'toHash' function to create unique keys for
// complex objects in the Set.
func toStr(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Struct:
		var r []string
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			name := t.Field(i).Name
			value := toStr(v.Field(i))
			r = append(r, fmt.Sprintf("%s:%s", name, value))
		}
		return "{" + strings.Join(r, ", ") + "}"
	case reflect.Array, reflect.Slice:
		var elements []string
		for i := 0; i < v.Len(); i++ {
			elements = append(elements, toStr(v.Index(i)))
		}
		return "[" + strings.Join(elements, ", ") + "]"
	case reflect.Map:
		var r []string
		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			r = append(r, fmt.Sprintf("%s:%s", toStr(key), toStr(value)))
		}
		return "{" + strings.Join(r, ", ") + "}"
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return "nil"
		}
		return toStr(v.Elem())
	case reflect.Func:
		if v.IsNil() {
			return "func:nil"
		}
		return v.Type().String() + " Value"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// isSimple determines the complexity of the objects in the set, i.e.,
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
func (s *Set[T]) isSimple() bool {
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

// Add adds the given items to the set.
//
// Example usage:
//
//	// Define a new set
//	s := New[int]()
//
//	// Add elements to the set
//	s.Add(1, 2, 3, 4)
//
//	// Now, the set contains the elements 1, 2, 3, and 4
func (s *Set[T]) Add(items ...T) {
	for _, v := range items {
		s.heap[s.toHash(v)] = v
	}
}

// Delete removes the given items from the set.
//
// Example usage:
//
//	// Define a new set and add some elements
//	s := New[int]()
//	s.Add(1, 2, 3, 4)
//
//	// Remove elements from the set
//	s.Delete(1, 3)
//
//	// Now, the set contains the elements 2 and 4
func (s *Set[T]) Delete(items ...T) {
	for _, v := range items {
		delete(s.heap, s.toHash(v))
	}
}

// Contains returns true if the set contains the given item.
//
// Example usage:
//
//	// Define a new set and add some elements
//	s := New[int]()
//	s.Add(1, 2, 3, 4)
//
//	// Check if the set contains certain elements
//	containsOne := s.Contains(1)  // returns true
//	containsFive := s.Contains(5) // returns false
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.heap[s.toHash(item)]
	return ok
}

// Elements returns all items in the set.
// This is useful when you need to iterate over the set,
// or when you need to convert the set to a slice ([]T).
// Note that the order of items is not guaranteed.
//
// Example usage:
//
//	s := New[int]()
//	s.Add(1, 2, 3, 4)
//	elements := s.Elements()  // elements is []int{1, 2, 3, 4}
func (s *Set[T]) Elements() []T {
	var items []T
	for _, v := range s.heap {
		items = append(items, v)
	}
	return items
}

// Len returns the number of items in the set.
// This is useful when you need to know how many items are in the set.
//
// Example usage:
//
//	s := New[int]()
//	s.Add(1, 2, 3, 4)
//	length := s.Len()  // length is 4
func (s *Set[T]) Len() int {
	return len(s.heap)
}

// Union returns a new set with all the items in both sets.
// This is useful when you want to merge two sets into a new one.
// Note that the result set will not have any duplicate items, even
// if the input sets do.
//
// Example usage:
//
//	s1 := New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := New[int]()
//	s2.Add(3, 4, 5)
//
//	union := s1.Union(s2)  // union contains 1, 2, 3, 4, 5
func (s *Set[T]) Union(set *Set[T]) *Set[T] {
	result := New[T](s.Elements()...)
	result.Add(set.Elements()...)
	return result
}

// Intersection returns a new set with items that exist only in both sets.
// This is useful when you want to find common items between two sets.
//
// Example usage:
//
//	s1 := New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := New[int]()
//	s2.Add(3, 4, 5)
//
//	intersection := s1.Intersection(s2)  // intersection contains 3
func (s *Set[T]) Intersection(set *Set[T]) *Set[T] {
	result := New[T]()
	for _, v := range s.heap {
		if set.Contains(v) {
			result.Add(v)
		}
	}
	return result
}

// Difference returns a new set with items in the first set but
// not in the second. This is useful when you want to find items
// that are unique to the first set.
//
// Example usage:
//
//	s1 := New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := New[int]()
//	s2.Add(3, 4, 5)
//
//	difference := s1.Difference(s2)  // difference contains 1, 2
func (s *Set[T]) Difference(set *Set[T]) *Set[T] {
	result := New[T]()
	for _, v := range s.heap {
		if !set.Contains(v) {
			result.Add(v)
		}
	}

	return result
}

// Diff is an alias for Difference.
func (s *Set[T]) Diff(set *Set[T]) *Set[T] {
	return s.Difference(set)
}

// SymmetricDifference returns a new set with items in either
// the first or second set but not both. This is useful when you want to find
// items that are unique to each set.
//
// Example usage:
//
//	s1 := New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := New[int]()
//	s2.Add(3, 4, 5)
//
//	symmetricDifference := s1.SymmetricDifference(s2)  // 1, 2, 4, 5
func (s *Set[T]) SymmetricDifference(set *Set[T]) *Set[T] {
	result := New[T]()
	for _, v := range s.heap {
		if !set.Contains(v) {
			result.Add(v)
		}
	}
	for _, v := range set.heap {
		if !s.Contains(v) {
			result.Add(v)
		}
	}

	return result
}

// Sdiff is an alias for SymmetricDifference.
func (s *Set[T]) Sdiff(set *Set[T]) *Set[T] {
	return s.SymmetricDifference(set)
}

// IsSubset returns true if all items in the first set exist in the second.
// This is useful when you want to check if all items of one set
// belong to another set.
//
// Example usage:
//
//	s1 := New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := New[int]()
//	s2.Add(1, 2, 3, 4, 5)
//
//	isSubset := s1.IsSubset(s2)  // isSubset is true
func (s *Set[T]) IsSubset(set *Set[T]) bool {
	for _, v := range s.heap {
		if !set.Contains(v) {
			return false
		}
	}

	return true
}

// IsSuperset returns true if all items in the second set exist in the first.
// This is useful when you want to check if one set contains all items
// of another set.
//
// Example usage:
//
//	s1 := New[int]()
//	s1.Add(1, 2, 3, 4, 5)
//
//	s2 := New[int]()
//	s2.Add(1, 2, 3)
//
//	isSuperset := s1.IsSuperset(s2)  // isSuperset is true
func (s *Set[T]) IsSuperset(set *Set[T]) bool {
	return set.IsSubset(s)
}

// Sorted returns a new set with items sorted in ascending order.
// This is useful when you want to sort the items in the set.
//
// Example usage:
//
//	s := New[int]()
//	s.Add(3, 2, 1)
//
//	sorted := s.Sorted() // sorted contains 1, 2, 3
func (s *Set[T]) Sorted(fns ...func(a, b T) bool) []T {
	// Create a temporary slice of sortMarker[T] to hold
	// the data and sort it.
	tmp := make([]sortMarker[T], 0, len(s.heap)) // here is the change
	for k, v := range s.heap {
		tmp = append(tmp, sortMarker[T]{name: k, value: v})
	}

	if len(fns) == 0 {
		sort.Slice(tmp, func(i, j int) bool {
			return tmp[i].name < tmp[j].name
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
	for i, v := range tmp {
		result[i] = v.value
	}

	return result
}

// Append adds all elements from the provided sets to the current set.
//
// Example usage:
//
//	s1 := New[int]()
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

// Extend is an alias for Append. It adds all elements from the provided sets
// to the current set.
//
// Example usage:
//
//	s1 := New[int]()
//	s1.Add(1, 2, 3)
//
//	s2 := New[int]()
//	s2.Add(4, 5, 6)
//
//	s1.Extend(s2)  // s1 now contains 1, 2, 3, 4, 5, 6
func (s *Set[T]) Extend(sets ...*Set[T]) {
	s.Append(sets...)
}

// Copy returns a new set with a copy of items in the set.
// This is useful when you want to copy the set.
//
// Example usage:
//
//	s := New[int]()
//	s.Add(1, 2, 3)
//
//	copied := s.Copy() // copied contains 1, 2, 3
func (s *Set[T]) Copy() *Set[T] {
	result := New[T]()
	for _, v := range s.heap {
		result.Add(v)
	}

	return result
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
//		s := New[int]()
//		s.Add(1, 2, 3)
//		s.Elements() // returns []int{1, 2, 3}
//
//	 s.Overwrite(5, 6, 7) // as s.Clear() and s.Add(5, 6, 7)
//		s.Elements() // returns []int{5, 6, 7}
func (s *Set[T]) Overwrite(items ...T) {
	s.Clear()
	s.Add(items...)
}
