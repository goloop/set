package set

import (
	"context"
	"reflect"
	"testing"
)

// complexType is helper for testing complex sets.
type complexType struct {
	field1 int
	field2 string
}

// userType is an another helper for testing complex sets.
type userType struct {
	Name string
	Age  int
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
		input    []complexType
		expected *Set[complexType]
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			expected: &Set[complexType]{
				heap: map[string]complexType{
					"{field1:1, field2:one}": {1, "one"},
					"{field1:2, field2:two}": {2, "two"},
				},
				simple: -1,
			},
		},
		{
			name:  "two",
			input: []complexType{},
			expected: &Set[complexType]{
				heap:   make(map[string]complexType),
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

// NewWithContext tests NewWithContext function.
func TestNewWithContext(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		expected *Set[complexType]
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			expected: &Set[complexType]{
				heap: map[string]complexType{
					"{field1:1, field2:one}": {1, "one"},
					"{field1:2, field2:two}": {2, "two"},
				},
				simple: -1,
			},
		},
		{
			name:  "two",
			input: []complexType{},
			expected: &Set[complexType]{
				heap:   make(map[string]complexType),
				simple: -1,
			},
		},
	}

	for _, tc := range tests {
		ctx := context.Background()
		result := NewWithContext(ctx, tc.input...)
		if !reflect.DeepEqual(result.Sorted(), tc.expected.Sorted()) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected.Sorted(), result.Sorted())
		}
	}
}

// AddWithContext tests AddWithContext function.
func TestAddWithContext(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		expected *Set[complexType]
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			expected: &Set[complexType]{
				heap: map[string]complexType{
					"{field1:1, field2:one}": {1, "one"},
					"{field1:2, field2:two}": {2, "two"},
				},
				simple: -1,
			},
		},
		{
			name:  "two",
			input: []complexType{},
			expected: &Set[complexType]{
				heap:   make(map[string]complexType),
				simple: -1,
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		AddWithContext(ctx, s, tc.input...)
		if !reflect.DeepEqual(s.Sorted(), tc.expected.Sorted()) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected.Sorted(), s.Sorted())
		}

		cancel()
		err := AddWithContext(ctx, s, tc.input...)
		if err == nil {
			t.Errorf("Test %s: expected error, but got nil", tc.name)
		}
	}
}

// // TestMapFn tests Map function.
// func TestMapFn(t *testing.T) {
// 	s := New[userType]()
// 	s.Add(userType{"John", 20}, userType{"Jane", 30})

// 	names := Map(s, func(item userType) string {
// 		return item.Name
// 	})

// 	expected := []string{"Jane", "John"}
// 	if v := names.Sorted(); !reflect.DeepEqual(v, expected) {
// 		t.Errorf("Map() failed, expected names = %v, got %v",
// 			expected, v)
// 	}
// }

// // TestReduceFn tests Reduce function.
// func TestReduceFn(t *testing.T) {
// 	s := New[userType]()
// 	s.Add(userType{"John", 20}, userType{"Jane", 30})

// 	sum := Reduce(s, func(acc int, item userType) int {
// 		return acc + item.Age
// 	})

// 	if sum != 50 {
// 		t.Errorf("Reduce() failed, expected sum = %v, got %v", 50, sum)
// 	}
// }

// // TestUnionFn tests Union function.
// func TestUnionFn(t *testing.T) {
// 	s1 := New[int](1, 2, 3)
// 	s2 := New[int](3, 4, 5)
// 	s3 := New[int](5, 6, 7)
// 	s4 := New[int](7, 8, 9)

// 	r := Union(s1, s2, s3, s4)

// 	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
// 	actual := r.Elements()
// 	sort.Ints(actual)

// 	if !reflect.DeepEqual(expected, actual) {
// 		t.Errorf("Expected %v, but got %v", expected, actual)
// 	}
// }

// // TestIntersectionFn tests Intersection function.
// func TestIntersectionFn(t *testing.T) {
// 	s1 := New[int](1, 2, 3)
// 	s2 := New[int](3, 4, 5)
// 	s3 := New[int](3, 6, 7)
// 	s4 := New[int](3, 8, 9)

// 	r := Intersection(s1, s2, s3, s4)

// 	expected := []int{3}
// 	actual := r.Elements()
// 	sort.Ints(actual)

// 	if !reflect.DeepEqual(expected, actual) {
// 		t.Errorf("Expected %v, but got %v", expected, actual)
// 	}
// }

// // TestDiffFn tests Diff function.
// func TestDiffFn(t *testing.T) {
// 	s1 := New[int](1, 2, 3)
// 	s2 := New[int](3, 4, 5)
// 	s3 := New[int](5, 6, 7)
// 	s4 := New[int](7, 8, 9)

// 	r := Diff(s1, s2, s3, s4)

// 	expected := []int{1, 2}
// 	actual := r.Elements()
// 	sort.Ints(actual)

// 	if !reflect.DeepEqual(expected, actual) {
// 		t.Errorf("Expected %v, but got %v", expected, actual)
// 	}
// }

// // TestSdiffFn tests Sdiff function.
// func TestSdiffFn(t *testing.T) {
// 	s1 := New[int](1, 2, 3)
// 	s2 := New[int](3, 4, 5)
// 	s3 := New[int](5, 6, 7)
// 	s4 := New[int](7, 8, 9)

// 	r := Sdiff(s1, s2, s3, s4)

// 	expected := []int{1, 2, 4, 6, 8, 9}
// 	actual := r.Elements()
// 	sort.Ints(actual)

// 	if !reflect.DeepEqual(expected, actual) {
// 		t.Errorf("Expected %v, but got %v", expected, actual)
// 	}
// }
