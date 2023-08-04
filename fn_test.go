package set

import (
	"reflect"
	"testing"
)

// Define a complex type for testing.
type ComplexType struct {
	field1 int
	field2 string
}

// Define a simple type for testing.
type UserType struct {
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

// TestMapFn tests Map function.
func TestMapFn(t *testing.T) {
	s := New[UserType]()
	s.Add(UserType{"John", 20}, UserType{"Jane", 30})

	names := Map(s, func(item UserType) string {
		return item.Name
	})

	expected := []string{"Jane", "John"}
	if v := names.Sorted(); !reflect.DeepEqual(v, expected) {
		t.Errorf("Map() failed, expected names = %v, got %v",
			expected, v)
	}
}

// TestReduceFn tests Reduce function.
func TestReduceFn(t *testing.T) {
	s := New[UserType]()
	s.Add(UserType{"John", 20}, UserType{"Jane", 30})

	sum := Reduce(s, func(acc int, item UserType) int {
		return acc + item.Age
	})

	if sum != 50 {
		t.Errorf("Reduce() failed, expected sum = %v, got %v", 50, sum)
	}
}
