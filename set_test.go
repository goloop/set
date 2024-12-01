package set

import (
	"context"
	"reflect"
	"sort"
	"testing"
)

type jsonTestStruct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TestToHashMethodSimple tests toHash method for simple types.
func TestToHashMethodSimple(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected uint64
	}{
		{
			name:     "integer 1",
			input:    1,
			expected: 12638134423997487868,
		},
		{
			name:     "integer 0",
			input:    0,
			expected: 12638135523509116079,
		},
	}

	set := New[int]()
	for _, tc := range tests {
		result, err := set.toHash(nil, tc.input)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tc.name, err)
		}

		if result != tc.expected {
			t.Errorf("%s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

// TestToHashMethodComplex tests toHash method for complex types.
func TestToHashMethodComplex(t *testing.T) {
	tests := []struct {
		name     string
		input    complexType
		expected uint64
	}{
		{
			name:     "complex {1, \"one\"}",
			input:    complexType{1, "one"},
			expected: 2272318830438166496,
		},
		{
			name:     "complex {2, \"two\"}",
			input:    complexType{2, "two"},
			expected: 2243055450779406681,
		},
	}

	set := New[complexType]()
	for _, tc := range tests {
		result, err := set.toHash(nil, tc.input)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tc.name, err)
		}

		if result != tc.expected {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

// TestIsSimpleMethod tests IsSimple method.
func TestIsSimpleMethod(t *testing.T) {
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

// TestIsComplexMethod tests IsComplex method.
func TestIsComplexMethod(t *testing.T) {
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

// TestAddWithContextMethod tests AddWithContext method.
func TestAddWithContextMethod(t *testing.T) {
	s := New[int]()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 100; i++ {
		if i == 10 {
			cancel()
		}

		if err := s.addWithContext(ctx, i); err != nil {
			break
		}
	}

	if s.Len() != 10 {
		t.Errorf("AddWithContext: expected length 10, but got %d", s.Len())
	}
}

// TestAddMethod tests Add method.
func TestAddMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	expected := &Set[int]{
		heap: map[uint64]int{
			12638134423997487868: 1,
			12638137722532372501: 2,
			12638136623020744290: 3,
			12638131125462603235: 4,
		},
		simple: 1,
	}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Add: expected %v, but got %v", expected, s)
	}
}

// TestDeleteWithContextMethod tests DeleteWithContext method.
func TestDeleteWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 1; i < 20; i++ {
		if i == 10 {
			cancel()
		}

		if err := s.deleteWithContext(ctx, i); err != nil {
			break
		}
	}

	if s.Len() != 10 {
		t.Errorf("DeleteWithContext: expected length 3, but got %d", s.Len())
	}
}

// TestDeleteMethod tests Delete method.
func TestDeleteMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)
	s.Delete(1, 3)

	expected := &Set[int]{
		heap: map[uint64]int{
			12638137722532372501: 2,
			12638131125462603235: 4,
		},
		simple: 1,
	}

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Delete: expected %v, but got %v", expected, s)
	}
}

// TestContainsWithContextMethod tests ContainsWithContext method.
func TestContainsWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if ok, err := s.containsWithContext(ctx, 3); !ok || err != nil {
		t.Errorf("ContainsWithContext: expected (true, nil), but got (%v, %v)",
			ok, err)
	}

	cancel()
	if ok, _ := s.containsWithContext(ctx, 3); ok {
		t.Errorf("ContainsWithContext: expected false, but got %v", ok)
	}
}

// TestContainsMethod tests Contains method.
func TestContainsMethod(t *testing.T) {
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

// TestElementsWithContextMethod tests ElementsWithContext method.
func TestElementsWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expected := []int{1, 2, 3, 4}
	result, _ := s.elementsWithContext(ctx)

	// Since the order of elements is not guaranteed,
	// we need to sort the slices before comparing them.
	sort.Ints(result)
	sort.Ints(expected)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	cancel()
	if _, err := s.elementsWithContext(ctx); err == nil {
		t.Errorf("ElementsWithContext: expected error")
	}
}

// TestElementsMethod tests for the Elements method.
func TestElementsMethod(t *testing.T) {
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

// TestSortedWithContextMethod tests SortedWithContext method.
func TestSortedWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(3, 2, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sorted, _ := s.sortedWithContext(ctx, func(a, b int) bool {
		return a < b
	})

	// Check that the sorted slice is in ascending order.
	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("SortedWithContext() = %v, want %v", sorted, expected)
	}

	// Test with a comparison function.
	descending, _ := s.sortedWithContext(ctx, func(a, b int) bool {
		return a > b
	})
	expectedDesc := []int{3, 2, 1}
	if !reflect.DeepEqual(descending, expectedDesc) {
		t.Errorf("SortedWithContext() = %v, want %v", descending, expectedDesc)
	}

	cancel()
	if _, err := s.sortedWithContext(ctx); err == nil {
		t.Errorf("SortedWithContext: expected error")
	}
}

// TestSortedMethod tests for the Sorted method.
func TestSortedMethod(t *testing.T) {
	s := New[int]()
	s.Add(3, 2, 1)

	sorted := s.Sorted(func(a, b int) bool { return a < b })

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

// TestFilteredWithContextMethod tests FilteredWithContext method.
func TestFilteredWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4, 5)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filtered, _ := s.filteredWithContext(ctx, func(item int) bool {
		return item > 3
	})

	expected := []int{4, 5}
	sort.Ints(filtered)
	if !reflect.DeepEqual(filtered, expected) {
		t.Errorf("FilteredWithContext() failed, expected "+
			"elements = %v, got %v", expected, filtered)
	}

	cancel()
	if _, err := s.filteredWithContext(ctx, func(item int) bool {
		return item > 3
	}); err == nil {
		t.Errorf("FilteredWithContext: expected error")
	}
}

// TestFilteredMethod tests for the Filtered method.
func TestFilteredMethod(t *testing.T) {
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

// TestLenMethod tests for the Len method.
func TestLenMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	expected := 4
	result := s.Len()

	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

// TestUnionWithContextMethod tests UnionWithContext method.
func TestUnionWithContextMethod(t *testing.T) {
	s1 := New[int]()
	s1.Add(3)

	s2 := New[int](0, 5, 7)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expected := New[int]()
	expected.Add(0, 3, 5, 7)

	result, _ := s1.unionWithContext(ctx, s2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v",
			expected.Elements(), result.Elements())
	}

	cancel()
	if _, err := s1.unionWithContext(ctx, s2); err == nil {
		t.Errorf("UnionWithContext: expected error")
	}
}

// TestUnionMethod tests for the Union method.
func TestUnionMethod(t *testing.T) {
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

// TestIntersectionWithContextMethod tests IntersectionWithContext method.
func TestIntersectionWithContextMethod(t *testing.T) {
	s1 := New[int]()
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(3, 4, 5)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expected := New[int]()
	expected.Add(3)

	result, _ := s1.intersectionWithContext(ctx, s2)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v",
			expected.Elements(), result.Elements())
	}

	cancel()
	if _, err := s1.intersectionWithContext(ctx, s2); err == nil {
		t.Errorf("IntersectionWithContext: expected error")
	}
}

// TestIntersectionMethod tests for the Intersection method.
func TestIntersectionMethod(t *testing.T) {
	s1 := New[int]()
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(3, 4, 5)

	expected := New[int]()
	expected.Add(3)

	result := s1.Inter(s2)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v",
			expected.Elements(), result.Elements())
	}
}

// TestDifferenceWithContextMethod tests DifferenceWithContext method.
func TestDifferenceWithContextMethod(t *testing.T) {
	s1 := New[int]()
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(3, 4, 5)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expected := New[int]()
	expected.Add(1, 2)

	result, _ := s1.differenceWithContext(ctx, s2)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v",
			expected.Elements(), result.Elements())
	}

	cancel()
	if _, err := s1.differenceWithContext(ctx, s2); err == nil {
		t.Errorf("DifferenceWithContext: expected error")
	}
}

// TestDifferenceMethod tests for the Difference method.
func TestDifferenceMethod(t *testing.T) {
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

// TestSymmetricDifferenceWithContextMethod tests
// SymmetricDifferenceWithContext method.
func TestSymmetricDifferenceWithContextMethod(t *testing.T) {
	s1 := New[int]()
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(3, 4, 5)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	expected := New[int]()
	expected.Add(1, 2, 4, 5)

	result, _ := s1.symmetricDifferenceWithContext(ctx, s2)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v",
			expected.Elements(), result.Elements())
	}

	cancel()
	if _, err := s1.symmetricDifferenceWithContext(ctx, s2); err == nil {
		t.Errorf("SymmetricDifferenceWithContext: expected error")
	}
}

// TestSymmetricDifferenceMethod tests for the SymmetricDifference method.
func TestSymmetricDifferenceMethod(t *testing.T) {
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

// TestMapWithContextMethod tests MapWithContext method.
func TestMapWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mapped, _ := s.mapWithContext(ctx, func(item int) int {
		return item * 2
	})

	expected := []int{2, 4, 6}
	if v := mapped.Sorted(func(a, b int) bool {
		return a < b
	}); !reflect.DeepEqual(v, expected) {
		t.Errorf("MapWithContext() failed, expected elements = %v, got %v",
			expected, v)
	}

	cancel()
	if _, err := s.mapWithContext(ctx, func(item int) int {
		return item * 2
	}); err == nil {
		t.Errorf("MapWithContext: expected error")
	}
}

// TestMapMethod tests for the Map method.
func TestMapMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	mapped := s.Map(func(item int) int {
		return item * 2
	})

	expected := []int{2, 4, 6}
	if v := mapped.Sorted(func(a, b int) bool {
		return a < b
	}); !reflect.DeepEqual(v, expected) {
		t.Errorf("Map() failed, expected elements = %v, got %v",
			expected, v)
	}
}

// TestReduceWithContextMethod tests ReduceWithContext method.
func TestReduceWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reduced, _ := s.reduceWithContext(ctx, func(acc, item int) int {
		return acc + item
	})

	if reduced != 6 {
		t.Errorf("ReduceWithContext() failed, expected %d, got %d",
			6, reduced)
	}

	cancel()
	if _, err := s.reduceWithContext(ctx, func(acc, item int) int {
		return acc + item
	}); err == nil {
		t.Errorf("ReduceWithContext: expected error")
	}
}

// TestReduceMethod tests for the Reduce method.
func TestReduceMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	reduced := s.Reduce(func(acc, item int) int {
		return acc + item
	})

	if reduced != 6 {
		t.Errorf("Reduce() failed, expected %d, got %d",
			6, reduced)
	}
}

// TestCopyWithContextMethod tests CopyWithContext method.
func TestCopyWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	copied, _ := s.copyWithContext(ctx)

	if !reflect.DeepEqual(copied, s) {
		t.Errorf("CopyWithContext() failed, expected %v, got %v",
			s, copied)
	}

	cancel()
	if _, err := s.copyWithContext(ctx); err == nil {
		t.Errorf("CopyWithContext: expected error")
	}
}

// TestCopyMethod tests for the Copy method.
func TestCopyMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	copied := s.Copy()

	if !reflect.DeepEqual(copied, s) {
		t.Errorf("Copy() failed, expected %v, got %v",
			s, copied)
	}
}

// TestAppendWithContextMethod tests AppendWithContext method.
func TestAppendWithContextMethod(t *testing.T) {
	s1 := New[int]()
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(4, 5, 6)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s1.appendWithContext(ctx, s2)

	expected := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(s1.Sorted(func(a, b int) bool {
		return a < b
	}), expected) {
		t.Errorf("AppendWithContext() = %v, want %v", s1.Sorted(), expected)
	}

	cancel()
	if err := s1.appendWithContext(ctx, s2); err == nil {
		t.Errorf("AppendWithContext: expected error")
	}
}

// TestAppendMethod tests for the Append method.
func TestAppendMethod(t *testing.T) {
	s1 := New[int]()
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(4, 5, 6)

	s1.Append(s2)

	expected := []int{1, 2, 3, 4, 5, 6}
	if !reflect.DeepEqual(s1.Sorted(func(a, b int) bool {
		return a < b
	}), expected) {
		t.Errorf("Append() = %v, want %v", s1.Sorted(), expected)
	}
}

// TestExtendWithContextMethod tests ExtendWithContext method.
func TestExtendWithContextMethod(t *testing.T) {
	// Initialize two sets
	s1 := New[int]()
	s1.Add(1, 2, 3)

	s2 := New[int]()
	s2.Add(4, 5, 6)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Extend s1 with s2
	s1.extendWithContext(ctx, []*Set[int]{s2})

	// Test that the extended set has the correct length
	if s1.Len() != 6 {
		t.Errorf("ExtendWithContext() failed, expected "+
			"length = %v, got %v", 6, s1.Len())
	}

	// Test that the extended set contains the correct items
	expected := []int{1, 2, 3, 4, 5, 6}
	sort.Ints(expected)
	sort.Ints(s1.Elements())
	if !reflect.DeepEqual(s1.Sorted(func(a, b int) bool {
		return a < b
	}), expected) {
		t.Errorf("ExtendWithContext() failed, expected elements = %v, got %v",
			expected, s1.Sorted())
	}

	cancel()
	if err := s1.extendWithContext(ctx, []*Set[int]{s2}); err == nil {
		t.Errorf("ExtendWithContext: expected error")
	}
}

// TestExtendMethod tests for the Extend method.
func TestExtendMethod(t *testing.T) {
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
	if !reflect.DeepEqual(s1.Sorted(func(a, b int) bool {
		return a < b
	}), expected) {
		t.Errorf("Extend() failed, expected elements = %v, got %v",
			expected, s1.Elements())
	}
}

// TestOverwriteWithContextMethod tests OverwriteWithContext method.
func TestOverwriteWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.overwriteWithContext(ctx, 5, 6, 7)

	expected := []int{5, 6, 7}
	if !reflect.DeepEqual(s.Sorted(func(a, b int) bool {
		return a < b
	}), expected) {
		t.Errorf("OverwriteWithContext() = %v, want %v", s.Sorted(), expected)
	}

	cancel()
	if err := s.overwriteWithContext(ctx, 1, 2, 3); err == nil {
		t.Errorf("OverwriteWithContext: expected error")
	}
}

// TestOverwriteMethod tests for the Overwrite method.
func TestOverwriteMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	s.Overwrite(5, 6, 7)

	expected := []int{5, 6, 7}
	if !reflect.DeepEqual(s.Sorted(func(a, b int) bool {
		return a < b
	}), expected) {
		t.Errorf("Overwrite() = %v, want %v", s.Sorted(), expected)
	}
}

// TestIsSubsetWithContextMethod tests IsSubsetWithContext method.
func TestIsSubsetWithContextMethod(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			got, err := tc.set1.isSubsetWithContext(ctx, tc.set2)
			if err != nil {
				t.Errorf("IsSubsetWithContext() = %v, want %v", err, nil)
			}

			if got != tc.expected {
				t.Errorf("IsSubsetWithContext() = %v, want %v",
					got, tc.expected)
			}

			cancel()
			ok, _ := tc.set1.isSubsetWithContext(ctx, tc.set2)
			if ok {
				t.Errorf("IsSubsetWithContext: expected false")
			}
		})
	}
}

// TestIsSubsetMethod tests for the IsSubset method.
func TestIsSubsetMethod(t *testing.T) {
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
		result := tc.set1.IsSub(tc.set2)
		if result != tc.expected {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
	}
}

// TestIsSupersetWithContextMethod tests IsSupersetWithContext method.
func TestIsSupersetWithContextMethod(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			got, err := tc.set1.isSupersetWithContext(ctx, tc.set2)
			if err != nil {
				t.Errorf("IsSupersetWithContext() = %v, want %v", err, nil)
			}

			if got != tc.expected {
				t.Errorf("IsSupersetWithContext() = %v, want %v",
					got, tc.expected)
			}

			// Test with cancelled context.
			cancel()
			ok, _ := tc.set1.isSupersetWithContext(ctx, tc.set2)
			if ok {
				t.Errorf("IsSupersetWithContext: expected false: "+
					"set1: %v, set2: %v", tc.set1.Sorted(), tc.set2.Sorted())
			}
		})
	}
}

// TestIsSupersetMethod tests for the IsSuperset method.
func TestIsSupersetMethod(t *testing.T) {
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
		result := tc.set1.IsSup(tc.set2)
		if result != tc.expected {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, result)
		}
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

// TestFilterWithContextMethod tests for the FilterWithContext method.
func TestFilterWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4, 5)

	filtered, err := s.filterWithContext(nil, func(item int) bool {
		return item > 3
	})
	if err != nil {
		t.Errorf("FilterWithContext() = %v, want %v", err, nil)
	}

	expected := []int{4, 5}
	sort.Ints(filtered.Elements())
	if v := filtered.Sorted(); !reflect.DeepEqual(v, expected) {
		t.Errorf("Filter() failed, expected elements = %v, got %v",
			expected, v)
	}

	// Test with cancelled context.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = s.filterWithContext(ctx, func(item int) bool {
		return item > 3
	})
	if err == nil {
		t.Errorf("FilterWithContext() = %v, want %v", err, nil)
	}
}

// TestFilterMethod tests for the Filter method.
func TestFilterMethod(t *testing.T) {
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

// TestAnyWithContextMethod tests for the AnyWithContext method.
func TestAnyWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	any, err := s.anyWithContext(nil, func(item int) bool {
		return item > 2
	})
	if err != nil {
		t.Errorf("AnyWithContext() = %v, want %v", err, nil)
	}

	if !any {
		t.Errorf("Any() failed, expected value = %v, got %v", true, any)
	}

	// Test with cancelled context.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = s.anyWithContext(ctx, func(item int) bool {
		return item > 2
	})

	if err == nil {
		t.Errorf("AnyWithContext() = %v, want %v", err, nil)
	}
}

// TestAnyMethod tests for the Any method.
func TestAnyMethod(t *testing.T) {
	s := New[int]()
	// Empty.
	any := s.Any(func(item int) bool {
		return item > 2
	})

	if any {
		t.Errorf("Any() failed, expected value = %v, got %v", false, any)
	}

	// Not empty.
	s.Add(1, 2, 3)
	any = s.Any(func(item int) bool {
		return item > 2
	})

	if !any {
		t.Errorf("Any() failed, expected value = %v, got %v", true, any)
	}
}

// TestAnyParallelMethod tests for the Any method.
func TestAnyParallelMethod(t *testing.T) {
	// Small processing block size.
	minLoadPerGoroutine = 5

	// Initialize of a large set.
	s := New[int]()
	for i := 0; i < 1000; i++ {
		if i == 800 {
			s.Add(-1)
		} else {
			s.Add(i)
		}
	}

	// The call will be in goroutines.
	any := s.Any(func(item int) bool {
		return item == -1
	})
	if !any {
		t.Errorf("Any() failed, expected value = %v, got %v", true, any)
	}

	// Cancel goroutine outside.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := s.anyWithContext(ctx, func(item int) bool {
		return item >= 0
	})

	if err == nil {
		t.Errorf("anyWithContext() = %v, want %v", err, nil)
	}
}

// TestAllWithContextMethod tests for the AllWithContext method.
func TestAllWithContextMethod(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3)

	all, err := s.allWithContext(nil, func(item int) bool {
		return item > 2
	})
	if err != nil {
		t.Errorf("allWithContext() = %v, want %v", err, nil)
	}

	if all {
		t.Errorf("All() failed, expected value = %v, got %v", false, all)
	}

	// Test with cancelled context.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ok, _ := s.allWithContext(ctx, func(item int) bool {
		return item > 2
	})

	if ok {
		t.Errorf("AllWithContext() = %v, want %v", false, ok)
	}
}

// TestAllMethod tests for the All method.
func TestAllMethod(t *testing.T) {
	s := New[int]()

	// Empty.
	all := s.All(func(item int) bool {
		return item > 2
	})
	if all {
		t.Errorf("All() failed, expected value = %v, got %v", false, all)
	}

	// Not empty.
	s.Add(1, 2, 3)
	all = s.All(func(item int) bool {
		return item > 2
	})

	if all {
		t.Errorf("All() failed, expected value = %v, got %v", false, all)
	}
}

// TestAllParallelMethod tests for the All method.
func TestAllParallelMethod(t *testing.T) {
	// Small processing block size.
	minLoadPerGoroutine = 5

	// Initialize of a large set.
	s := New[int]()
	for i := 0; i < 1000; i++ {
		if i == 800 {
			s.Add(-1)
		} else {
			s.Add(i)
		}
	}

	// The call will be in goroutines.
	all := s.All(func(item int) bool {
		return item >= 0
	})
	if all {
		t.Errorf("All() failed, expected value = %v, got %v", false, all)
	}

	// Cancel goroutine outside.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := s.allWithContext(ctx, func(item int) bool {
		return item >= 0
	})

	if err == nil {
		t.Errorf("allWithContext() = %v, want %v", err, nil)
	}
}

// TestSetJSON tests Marshal/Unmarshal.
func TestSetJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    *Set[int]
		wantErr  bool
		validate func(*testing.T, *Set[int])
	}{
		{
			name:    "simple integers",
			input:   New(1, 2, 3, 4, 5),
			wantErr: false,
			validate: func(t *testing.T, s *Set[int]) {
				if s.Len() != 5 {
					t.Errorf("expected length 5, got %d", s.Len())
				}
				for i := 1; i <= 5; i++ {
					if !s.Contains(i) {
						t.Errorf("missing element %d", i)
					}
				}
			},
		},
		{
			name:    "empty set",
			input:   New[int](),
			wantErr: false,
			validate: func(t *testing.T, s *Set[int]) {
				if s.Len() != 0 {
					t.Errorf("expected empty set, got length %d", s.Len())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			data, err := tt.input.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Unmarshal into new set
			newSet := New[int]()
			err = newSet.UnmarshalJSON(data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Validate
			tt.validate(t, newSet)
		})
	}
}

// TestSetJSONWithStruct tests Marshal/Unmarshal with struct.
func TestSetJSONWithStruct(t *testing.T) {
	original := New[jsonTestStruct]()
	original.Add(
		jsonTestStruct{ID: 1, Name: "One"},
		jsonTestStruct{ID: 2, Name: "Two"},
		jsonTestStruct{ID: 3, Name: "Three"},
	)

	// Marshal
	data, err := original.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	// Unmarshal
	newSet := New[jsonTestStruct]()
	err = newSet.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	// Validate
	if newSet.Len() != original.Len() {
		t.Errorf("expected length %d, got %d", original.Len(), newSet.Len())
	}

	for _, item := range original.Elements() {
		if !newSet.Contains(item) {
			t.Errorf("missing element %v", item)
		}
	}
}
