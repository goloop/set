package set

import (
	"reflect"
	"sort"
	"testing"
)

// Define a complex type for testing.
type ComplexType struct {
	field1 int
	field2 string
}

// TestNewSimple tests New function for simple type.
func TestNewSimple(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected *Set[int]
	}{
		{
			name:  "[]int{1, 2, 3, 4, 5}",
			input: []int{1, 2, 3, 4, 5},
			expected: &Set[int]{
				heap: map[string]int{
					"1": 1,
					"2": 2,
					"3": 3,
					"4": 4,
					"5": 5,
				},
				simple: 1,
			},
		},
		{
			name:  "[]int{}",
			input: []int{},
			expected: &Set[int]{
				heap:   make(map[string]int),
				simple: 1,
			},
		},
	}

	for _, tc := range tests {
		result := New(tc.input...)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("%s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

// TestNewComplex tests New function for complex type.
func TestNewComplex(t *testing.T) {
	tests := []struct {
		name     string
		input    []ComplexType
		expected *Set[ComplexType]
	}{
		{
			name: "one",
			input: []ComplexType{
				{1, "one"},
				{2, "two"},
			},
			expected: &Set[ComplexType]{
				heap: map[string]ComplexType{
					"{field1:1, field2:one}": {1, "one"},
					"{field1:2, field2:two}": {2, "two"},
				},
				simple: -1,
			},
		},
		{
			name:  "two",
			input: []ComplexType{},
			expected: &Set[ComplexType]{
				heap:   make(map[string]ComplexType),
				simple: -1,
			},
		},
	}

	for _, tc := range tests {
		result := New(tc.input...)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

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
		result := set.toHash(tc.input)
		if result != tc.expected {
			t.Errorf("%s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

// TestToHashComplex tests toHash function for complex types.
func TestToHashComplex(t *testing.T) {
	tests := []struct {
		name     string
		input    ComplexType
		expected string
	}{
		{
			name:     "complex {1, \"one\"}",
			input:    ComplexType{1, "one"},
			expected: "{field1:1, field2:one}",
		},
		{
			name:     "complex {2, \"two\"}",
			input:    ComplexType{2, "two"},
			expected: "{field1:2, field2:two}",
		},
	}

	set := New[ComplexType]()
	for _, tc := range tests {
		result := set.toHash(tc.input)
		if result != tc.expected {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

// TestIsSimple tests isSimple function.
func TestIsSimple(t *testing.T) {
	t.Parallel()

	t.Run("simple types", func(t *testing.T) {
		t.Parallel()

		// Test with int type.
		if !New[int](1, 2, 3).isSimple() {
			t.Error("Int: expected set of type int to be simple")
		}

		// Test with string type.
		if !New[string]("a", "b", "c").isSimple() {
			t.Error("String: expected set of type string to be simple")
		}

		// Test with bool type.
		if !New[bool](true, false).isSimple() {
			t.Error("Bool: expected set of type bool to be simple")
		}

		// Test with byte type.
		if !New[byte]('a', 'b', 'c').isSimple() {
			t.Error("Byte: expected set of type byte to be simple")
		}

		// Test with rune type.
		if !New[rune]('a', 'b', 'c').isSimple() {
			t.Error("Rune: expected set of type rune to be simple")
		}

		// Test with float32 type.
		if !New[float32](1.1, 2.2, 3.3).isSimple() {
			t.Error("Float32: expected set of type float32 to be simple")
		}

		// Test with complex64 type.
		if !New[complex64](complex(1, 2), complex(3, 4)).isSimple() {
			t.Error("Complex64: expected set of type complex64 to be simple")
		}

		// Test with complex128 type.
		if !New[complex128](complex(1, 2), complex(3, 4)).isSimple() {
			t.Error("Complex128: expected set of type complex128 to be simple")
		}
	})

	t.Run("complex types", func(t *testing.T) {
		t.Parallel()

		// Test with struct type.
		if New[ComplexType](ComplexType{1, "one"}, ComplexType{2, "two"}).isSimple() {
			t.Error("Struct: expected set of type struct to be complex")
		}

		// Test with slice type.
		slice := []int{1, 2, 3}
		if New[[]int](slice).isSimple() {
			t.Error("Slice: expected set of type slice to be complex")
		}

		// Test with map type.
		m := map[int]string{1: "one", 2: "two", 3: "three"}
		if New[map[int]string](m).isSimple() {
			t.Error("Map: expected set of type map to be complex")
		}

		// Test with func type.
		if New[func()](func() {}).isSimple() {
			t.Error("Func: expected set of type func to be complex")
		}

		// Test with chan type.
		ch := make(chan int)
		if New[chan int](ch).isSimple() {
			t.Error("Chan: expected set of type chan to be complex")
		}

		// Test with array type.
		arr := [3]int{1, 2, 3}
		if New[[3]int](arr).isSimple() {
			t.Error("Array: expected set of type array to be complex")
		}

		// Test with pointer type.
		ptr := new(int)
		if New[*int](ptr).isSimple() {
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
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(3, 4, 5)

	expected := New[int]()
	expected.Add(1, 2, 3, 4, 5)

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
		result := tc.set1.Difference(tc.set2)
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
		result := tc.set1.SymmetricDifference(tc.set2)
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
