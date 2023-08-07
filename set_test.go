package set

import (
	"reflect"
	"sort"
	"testing"
)

// TestToHashSimple tests toHash function for simple types.
func TestToHashSimple(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "integer 1",
			input:    1,
			expected: "1",
		},
		{
			name:     "integer 0",
			input:    0,
			expected: "0",
		},
	}

	set := New[int]()
	for _, tc := range tests {
		result := set.toHash(nil, tc.input)
		if result != tc.expected {
			t.Errorf("%s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

// TestToStr tests toStr function.
func TestToStr(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name:  "Pointer",
			input: new(int),
			want:  "0",
		},
		{
			name:  "NilPointer",
			input: (*int)(nil),
			want:  "nil",
		},
		{
			name:  "Interface",
			input: (interface{})(new(int)),
			want:  "0",
		},
		{
			name:  "Func",
			input: func() {},
			want:  "func() Value",
		},
		{
			name:  "NilFunc",
			input: (func())(nil),
			want:  "func:nil",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := toStr(nil, reflect.ValueOf(test.input))
			if got != test.want {
				t.Errorf("toStr(%s) = %s, want %s", test.name, got, test.want)
			}
		})
	}
}

// TestToHashComplex tests toHash function for complex types.
func TestToHashComplex(t *testing.T) {
	tests := []struct {
		name     string
		input    complexType
		expected string
	}{
		{
			name:     "complex {1, \"one\"}",
			input:    complexType{1, "one"},
			expected: "{field1:1, field2:one}",
		},
		{
			name:     "complex {2, \"two\"}",
			input:    complexType{2, "two"},
			expected: "{field1:2, field2:two}",
		},
	}

	set := New[complexType]()
	for _, tc := range tests {
		result := set.toHash(nil, tc.input)
		if result != tc.expected {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

// TestIsSimple tests IsSimple function.
func TestIsSimple(t *testing.T) {
	t.Parallel()

	t.Run("simple types", func(t *testing.T) {
		t.Parallel()

		// Test with int type.
		if !New[int](1, 2, 3).IsSimple() {
			t.Error("Int: expected set of type int to be simple")
		}

		// Test with string type.
		if !New[string]("a", "b", "c").IsSimple() {
			t.Error("String: expected set of type string to be simple")
		}

		// Test with bool type.
		if !New[bool](true, false).IsSimple() {
			t.Error("Bool: expected set of type bool to be simple")
		}

		// Test with byte type.
		if !New[byte]('a', 'b', 'c').IsSimple() {
			t.Error("Byte: expected set of type byte to be simple")
		}

		// Test with rune type.
		if !New[rune]('a', 'b', 'c').IsSimple() {
			t.Error("Rune: expected set of type rune to be simple")
		}

		// Test with float32 type.
		if !New[float32](1.1, 2.2, 3.3).IsSimple() {
			t.Error("Float32: expected set of type float32 to be simple")
		}

		// Test with complex64 type.
		if !New[complex64](complex(1, 2), complex(3, 4)).IsSimple() {
			t.Error("Complex64: expected set of type complex64 to be simple")
		}

		// Test with complex128 type.
		if !New[complex128](complex(1, 2), complex(3, 4)).IsSimple() {
			t.Error("Complex128: expected set of type complex128 to be simple")
		}
	})

	t.Run("complex types", func(t *testing.T) {
		t.Parallel()

		// Test with struct type.
		if New[complexType](
			complexType{1, "one"},
			complexType{2, "two"}).IsSimple() {
			t.Error("Struct: expected set of type struct to be complex")
		}

		// Test with slice type.
		slice := []int{1, 2, 3}
		if New[[]int](slice).IsSimple() {
			t.Error("Slice: expected set of type slice to be complex")
		}

		// Test with map type.
		m := map[int]string{1: "one", 2: "two", 3: "three"}
		if New[map[int]string](m).IsSimple() {
			t.Error("Map: expected set of type map to be complex")
		}

		// Test with func type.
		if New[func()](func() {}).IsSimple() {
			t.Error("Func: expected set of type func to be complex")
		}

		// Test with chan type.
		ch := make(chan int)
		if New[chan int](ch).IsSimple() {
			t.Error("Chan: expected set of type chan to be complex")
		}

		// Test with array type.
		arr := [3]int{1, 2, 3}
		if New[[3]int](arr).IsSimple() {
			t.Error("Array: expected set of type array to be complex")
		}

		// Test with pointer type.
		ptr := new(int)
		if New[*int](ptr).IsSimple() {
			t.Error("Pointer: expected set of type pointer to be complex")
		}
	})
}

// TestIsComplex tests IsComplex function.
func TestIsComplex(t *testing.T) {
	t.Parallel()

	t.Run("simple types", func(t *testing.T) {
		t.Parallel()

		// Test with int type.
		if New[int](1, 2, 3).IsComplex() {
			t.Error("Int: expected set of type int to be simple")
		}

		// Test with string type.
		if New[string]("a", "b", "c").IsComplex() {
			t.Error("String: expected set of type string to be simple")
		}

		// Test with bool type.
		if New[bool](true, false).IsComplex() {
			t.Error("Bool: expected set of type bool to be simple")
		}

		// Test with byte type.
		if New[byte]('a', 'b', 'c').IsComplex() {
			t.Error("Byte: expected set of type byte to be simple")
		}

		// Test with rune type.
		if New[rune]('a', 'b', 'c').IsComplex() {
			t.Error("Rune: expected set of type rune to be simple")
		}

		// Test with float32 type.
		if New[float32](1.1, 2.2, 3.3).IsComplex() {
			t.Error("Float32: expected set of type float32 to be simple")
		}

		// Test with complex64 type.
		if New[complex64](complex(1, 2), complex(3, 4)).IsComplex() {
			t.Error("Complex64: expected set of type complex64 to be simple")
		}

		// Test with complex128 type.
		if New[complex128](complex(1, 2), complex(3, 4)).IsComplex() {
			t.Error("Complex128: expected set of type complex128 to be simple")
		}
	})

	t.Run("complex types", func(t *testing.T) {
		t.Parallel()

		// Test with struct type.
		if !New[complexType](
			complexType{1, "one"},
			complexType{2, "two"}).IsComplex() {
			t.Error("Struct: expected set of type struct to be complex")
		}

		// Test with slice type.
		slice := []int{1, 2, 3}
		if !New[[]int](slice).IsComplex() {
			t.Error("Slice: expected set of type slice to be complex")
		}

		// Test with map type.
		m := map[int]string{1: "one", 2: "two", 3: "three"}
		if !New[map[int]string](m).IsComplex() {
			t.Error("Map: expected set of type map to be complex")
		}

		// Test with func type.
		if !New[func()](func() {}).IsComplex() {
			t.Error("Func: expected set of type func to be complex")
		}

		// Test with chan type.
		ch := make(chan int)
		if !New[chan int](ch).IsComplex() {
			t.Error("Chan: expected set of type chan to be complex")
		}

		// Test with array type.
		arr := [3]int{1, 2, 3}
		if !New[[3]int](arr).IsComplex() {
			t.Error("Array: expected set of type array to be complex")
		}

		// Test with pointer type.
		ptr := new(int)
		if !New[*int](ptr).IsComplex() {
			t.Error("Pointer: expected set of type pointer to be complex")
		}
	})
}

// TestAdd tests Add function.
func TestAdd(t *testing.T) {
	s := New[int]()

	s.Add(1, 2, 3, 4)

	expected := &Set[int]{
		heap: map[string]int{
			"1": 1,
			"2": 2,
			"3": 3,
			"4": 4,
		},
		simple: 1,
	}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Add: expected %v, but got %v", expected, s)
	}
}

// TestDelete tests Delete function.
func TestDelete(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	s.Delete(1, 3)

	expected := &Set[int]{
		heap: map[string]int{
			"2": 2,
			"4": 4,
		},
		simple: 1,
	}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Delete: expected %v, but got %v", expected, s)
	}
}

// TestClear tests Clear function.
func TestContains(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	tests := []struct {
		input    int
		expected bool
	}{
		{
			input:    1,
			expected: true,
		},
		{
			input:    5,
			expected: false,
		},
	}

	for _, tc := range tests {
		result := s.Contains(tc.input)
		if result != tc.expected {
			t.Errorf("Contains(%d): expected %v, but got %v",
				tc.input, tc.expected, result)
		}
	}
}

// TestElements tests for the Elements method.
func TestElements(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	expected := []int{1, 2, 3, 4}
	result := s.Elements()

	// Since the order of elements is not guaranteed,
	// we need to sort the slices before comparing them.
	sort.Ints(result)
	sort.Ints(expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

// TestLen tests for the Len method.
func TestLen(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	expected := 4
	result := s.Len()

	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

// TestUnion tests for the Union method.
func TestUnion(t *testing.T) {
	s1 := New[int]()
	s1.Add(3)

	s2 := New[int](0, 5, 7)

	expected := New[int]()
	expected.Add(0, 3, 5, 7)

	result := s1.Union(s2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v",
			expected.Elements(), result.Elements())
	}
}

// TestIntersection tests for the Intersection method.
func TestIntersection(t *testing.T) {
	s1 := New[int]()
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(3, 4, 5)

	expected := New[int]()
	expected.Add(3)

	result := s1.Intersection(s2)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v",
			expected.Elements(), result.Elements())
	}
}

// TestDifference tests for the Difference method.
func TestDifference(t *testing.T) {
	tests := []struct {
		name     string
		set1     *Set[int]
		set2     *Set[int]
		expected *Set[int]
	}{
		{
			name:     "Test difference between two sets",
			set1:     New[int](1, 2, 3),
			set2:     New[int](3, 4, 5),
			expected: New[int](1, 2),
		},
		{
			name:     "Test difference with no common elements",
			set1:     New[int](1, 2, 3),
			set2:     New[int](4, 5, 6),
			expected: New[int](1, 2, 3),
		},
	}

	for _, tc := range tests {
		result := tc.set1.Diff(tc.set2)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected.Elements(), result.Elements())
		}
	}
}

// TestSymmetricDifference tests for the SymmetricDifference method.
func TestSymmetricDifference(t *testing.T) {
	tests := []struct {
		name     string
		set1     *Set[int]
		set2     *Set[int]
		expected *Set[int]
	}{
		{
			name:     "Test symmetric difference between two sets",
			set1:     New[int](1, 2, 3),
			set2:     New[int](3, 4, 5),
			expected: New[int](1, 2, 4, 5),
		},
		{
			name:     "Test symmetric difference with no common elements",
			set1:     New[int](1, 2, 3),
			set2:     New[int](4, 5, 6),
			expected: New[int](1, 2, 3, 4, 5, 6),
		},
	}

	for _, tc := range tests {
		result := tc.set1.Sdiff(tc.set2)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected.Elements(), result.Elements())
		}
	}
}

// TestIsSubset tests for the IsSubset method.
func TestIsSubset(t *testing.T) {
	tests := []struct {
		name     string
		set1     *Set[int]
		set2     *Set[int]
		expected bool
	}{
		{
			name:     "Test when set1 is a subset of set2",
			set1:     New[int](1, 2, 3),
			set2:     New[int](1, 2, 3, 4, 5),
			expected: true,
		},
		{
			name:     "Test when set1 is not a subset of set2",
			set1:     New[int](1, 2, 3, 4, 5),
			set2:     New[int](1, 2, 3),
			expected: false,
		},
	}

	for _, tc := range tests {
		result := tc.set1.IsSubset(tc.set2)
		if result != tc.expected {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

// TestIsSuperset tests for the IsSuperset method.
func TestIsSuperset(t *testing.T) {
	tests := []struct {
		name     string
		set1     *Set[int]
		set2     *Set[int]
		expected bool
	}{
		{
			name:     "Test when set1 is a superset of set2",
			set1:     New[int](1, 2, 3, 4, 5),
			set2:     New[int](1, 2, 3),
			expected: true,
		},
		{
			name:     "Test when set1 is not a superset of set2",
			set1:     New[int](1, 2, 3),
			set2:     New[int](1, 2, 3, 4, 5),
			expected: false,
		},
	}

	for _, tc := range tests {
		result := tc.set1.IsSuperset(tc.set2)
		if result != tc.expected {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

// TestSorted tests for the Sorted method.
func TestSorted(t *testing.T) {
	s := New[int]()
	s.Add(3, 2, 1)

	sorted := s.Sorted()

	// Check that the sorted slice is in ascending order.
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Sorted() = %v, want %v", sorted, expected)
	}

	// Test with a comparison function.
	descending := s.Sorted(func(a, b int) bool { return a > b })
	expectedDesc := []int{3, 2, 1}
	if !reflect.DeepEqual(descending, expectedDesc) {
		t.Errorf("Sorted() = %v, want %v", descending, expectedDesc)
	}
}

// TestAppend tests for the Append method.
func TestAppend(t *testing.T) {
	s1 := New[int]()
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(4, 5, 6)

	s1.Append(s2)

	expected := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(s1.Sorted(), expected) {
		t.Errorf("Append() = %v, want %v", s1.Sorted(), expected)
	}
}

// TestExtend tests for the Extend method.
func TestExtend(t *testing.T) {
	// Initialize two sets
	s1 := New[int]()
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(4, 5, 6)

	// Extend s1 with s2
	s1.Extend([]*Set[int]{s2})

	// Test that the extended set has the correct length
	if s1.Len() != 6 {
		t.Errorf("Extend() failed, expected length = %v, got %v", 6, s1.Len())
	}

	// Test that the extended set contains the correct items
	expected := []int{1, 2, 3, 4, 5, 6}
	sort.Ints(expected)
	sort.Ints(s1.Elements())
	if !reflect.DeepEqual(s1.Sorted(), expected) {
		t.Errorf("Extend() failed, expected elements = %v, got %v",
			expected, s1.Elements())
	}
}

// TestCopy tests for the Copy method.
func TestCopy(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	copied := s.Copy()

	// Check that the copied set contains the same elements
	// as the original set.
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(copied.Sorted(), expected) {
		t.Errorf("Copy() = %v, want %v", copied.Sorted(), expected)
	}

	// Check that modifying the original set does not affect the copied set.
	s.Add(4)
	if copied.Contains(4) {
		t.Errorf("Copy() did not create a deep copy")
	}
}

// TestClear tests for the Clear method.
func TestClear(t *testing.T) {
	// Initialize a new set
	s := New[int]()
	s.Add(1, 2, 3)

	// Clear the set
	s.Clear()

	// Test that the set is empty after clearing
	if s.Len() != 0 {
		t.Errorf("Clear() failed, expected length = %v, got %v", 0, s.Len())
	}
}

// TestOverwrite tests for the Overwrite method.
func TestOverwrite(t *testing.T) {
	// Initialize a new set
	s := New[int]()
	s.Add(1, 2, 3)

	// Overwrite the set
	s.Overwrite(4, 5, 6)

	// Test that the set has the correct length after overwriting
	if s.Len() != 3 {
		t.Errorf("Overwrite() failed, expected length = %v, got %v",
			3, s.Len())
	}

	// Test that the set contains the correct items after overwriting
	expected := []int{4, 5, 6}
	sort.Ints(expected)
	sort.Ints(s.Elements())
	if v := s.Sorted(); !reflect.DeepEqual(v, expected) {
		t.Errorf("Overwrite() failed, expected elements = %v, got %v",
			expected, v)
	}
}

// TestFiltered tests for the Filtered method.
func TestFiltered(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4, 5)

	filtered := s.Filtered(func(item int) bool {
		return item > 3
	})

	expected := []int{4, 5}
	sort.Ints(filtered)
	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("Filtered() failed, expected elements = %v, got %v",
			expected, filtered)
	}
}

// TestFilter tests for the Filter method.
func TestFilter(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4, 5)

	filtered := s.Filter(func(item int) bool {
		return item > 3
	})

	expected := []int{4, 5}
	sort.Ints(filtered.Elements())
	if v := filtered.Sorted(); !reflect.DeepEqual(v, expected) {
		t.Errorf("Filter() failed, expected elements = %v, got %v",
			expected, v)
	}
}

// TestMap tests for the Map method.
func TestMap(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	mapped := s.Map(func(item int) int {
		return item * 2
	})

	expected := []int{2, 4, 6}
	if v := mapped.Sorted(); !reflect.DeepEqual(v, expected) {
		t.Errorf("Map() failed, expected elements = %v, got %v",
			expected, v)
	}
}

// TestReduce tests for the Reduce method.
func TestReduce(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	sum := s.Reduce(func(acc, item int) int {
		return acc + item
	})

	if sum != 6 {
		t.Errorf("Reduce() failed, expected value = %v, got %v", 6, sum)
	}
}

// TestAny tests for the Any method.
func TestAny(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	any := s.Any(func(item int) bool {
		return item > 2
	})

	if !any {
		t.Errorf("Any() failed, expected value = %v, got %v", true, any)
	}
}

// TestAll tests for the All method.
func TestAll(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	all := s.All(func(item int) bool {
		return item > 2
	})

	if all {
		t.Errorf("All() failed, expected value = %v, got %v", false, all)
	}
}
