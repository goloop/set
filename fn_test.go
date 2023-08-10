package set

import (
	"context"
	"reflect"
	"sort"
	"testing"
)

// complexType is helper for testing complex sets.
type complexType struct {
	FieldA int
	FieldB string
}

// userType is an another helper for testing complex sets.
type userType struct {
	Name string
	Age  int
}

// TestParallelTasks tests ParallelTasks function.
func TestParallelTasks(t *testing.T) {
	// Testing the default value.
	if got := ParallelTasks(); got != parallelTasks {
		t.Errorf("ParallelTasks() = %d; want %d", got, parallelTasks)
	}

	// Testing addition of values.
	if got := ParallelTasks(3, 4); got != 7 {
		t.Errorf("ParallelTasks(3, 4) = %d; want 7", got)
	}

	// Testing that the value is set to maxParallelTasks
	// if it exceeds the maximum.
	if got := ParallelTasks(maxParallelTasks + 1); got != maxParallelTasks {
		t.Errorf("ParallelTasks(%d) = %d; want %d", maxParallelTasks+1,
			got, maxParallelTasks)
	}

	// Testing that the value is set to 1 if it is less than or equal to zero.
	if got := ParallelTasks(-3); got != 1 {
		t.Errorf("ParallelTasks(-3) = %d; want 1", got)
	}

	// Testing that the new value is applied and can be retrieved.
	if got, want := ParallelTasks(5), 5; got != want {
		t.Errorf("ParallelTasks(5) = %d; want %d", got, want)
	}
	if got := ParallelTasks(); got != 5 {
		t.Errorf("ParallelTasks() after setting = %d; want 5", got)
	}
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
				heap: map[uint64]int{
					12638134423997487868: 1,
					12638137722532372501: 2,
					12638136623020744290: 3,
					12638131125462603235: 4,
					12638130025950975024: 5,
				},
				simple: 1,
			},
		},
		{
			name:  "[]int{}",
			input: []int{},
			expected: &Set[int]{
				heap:   make(map[uint64]int),
				simple: 1,
			},
		},
	}

	for _, tc := range tests {
		result := New(tc.input...)
		if !reflect.DeepEqual(result.Sorted(), tc.expected.Sorted()) {
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
				heap: map[uint64]complexType{
					2272318830438166496: {1, "one"},
					2243055450779406681: {2, "two"},
				},
				simple: -1,
			},
		},
		{
			name:  "two",
			input: []complexType{},
			expected: &Set[complexType]{
				heap:   make(map[uint64]complexType),
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
				heap: map[uint64]complexType{
					2272318830438166496: {1, "one"},
					2243055450779406681: {2, "two"},
				},
				simple: -1,
			},
		},
		{
			name:  "two",
			input: []complexType{},
			expected: &Set[complexType]{
				heap:   make(map[uint64]complexType),
				simple: -1,
			},
		},
	}

	for _, tc := range tests {
		ctx := context.Background()
		result := NewWithContext(ctx, tc.input...)
		if !reflect.DeepEqual(result.Sorted(), tc.expected.Sorted()) {
			t.Errorf("Test %s (%v): expected %v, but got %v",
				tc.name, result, tc.expected.Sorted(), result.Sorted())
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
				heap: map[uint64]complexType{
					2272318830438166496: {1, "one"},
					2243055450779406681: {2, "two"},
				},
				simple: -1,
			},
		},
		{
			name:  "two",
			input: []complexType{},
			expected: &Set[complexType]{
				heap:   make(map[uint64]complexType),
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
			t.Errorf("Test %s (%v): expected %v, but got %v",
				tc.name, s, tc.expected.Sorted(), s.Sorted())
		}

		cancel()
		AddWithContext(ctx, s, complexType{3, "three"})
		if s.Len() != tc.expected.Len() {
			t.Errorf("Test %s (%v): expected %v, but got %v",
				tc.name, s, tc.expected.Sorted(), s.Sorted())
		}
	}
}

// TestAdd tests Add function.
func TestAdd(t *testing.T) {
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
				heap: map[uint64]complexType{
					2272318830438166496: {1, "one"},
					2243055450779406681: {2, "two"},
				},
				simple: -1,
			},
		},
		{
			name:  "two",
			input: []complexType{},
			expected: &Set[complexType]{
				heap:   make(map[uint64]complexType),
				simple: -1,
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		Add(s, tc.input...)
		if !reflect.DeepEqual(s.Sorted(), tc.expected.Sorted()) {
			t.Errorf("Test %s (%v): expected %v, but got %v",
				tc.name, s, tc.expected.Sorted(), s.Sorted())
		}
	}
}

// TestDeleteWithContext tests DeleteWithContext function.
func TestDeleteWithContext(t *testing.T) {
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
				heap: map[uint64]complexType{
					1: {2, "two"},
				},
				simple: -1,
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		DeleteWithContext(ctx, s, tc.input[0])
		if !reflect.DeepEqual(s.Sorted(), tc.expected.Sorted()) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected.Sorted(), s.Sorted())
		}

		cancel()
		DeleteWithContext(ctx, s, tc.input[1])
		if s.Len() != tc.expected.Len() {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected.Sorted(), s.Sorted())
		}
	}
}

// TestDelete tests Delete function.
func TestDelete(t *testing.T) {
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
				heap: map[uint64]complexType{
					1: {2, "two"},
				},
				simple: -1,
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		Delete(s, tc.input[0])
		if !reflect.DeepEqual(s.Sorted(), tc.expected.Sorted()) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected.Sorted(), s.Sorted())
		}
	}
}

// TestContainsWithContext tests ContainsWithContext function.
func TestContainsWithContext(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		expected bool
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		v, _ := ContainsWithContext(ctx, s, tc.input[0])
		if v != tc.expected {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}

		cancel()
		v, _ = ContainsWithContext(ctx, s, tc.input[1])
		if v != false {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, false, v)
		}
	}
}

// TestContains tests Contains function.
func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		expected bool
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		v := Contains(s, tc.input[0])
		if v != tc.expected {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}
	}
}

// TetsElementsWithContext tests ElementsWithContext function.
func TestElementsWithContext(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			expected: []complexType{
				{1, "one"},
				{2, "two"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		v, _ := ElementsWithContext(ctx, s)
		sort.Slice(v, func(i, j int) bool {
			return v[i].FieldA < v[j].FieldA
		})
		sort.Slice(tc.expected, func(i, j int) bool {
			return tc.expected[i].FieldA < tc.expected[j].FieldA
		})

		if !reflect.DeepEqual(v, tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}

		cancel()
		v, _ = ElementsWithContext(ctx, s)
		if len(v) != 0 {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, []complexType{}, v)
		}
	}
}

// TestElements tests Elements function.
func TestElements(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			expected: []complexType{
				{1, "one"},
				{2, "two"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		v := Elements(s)
		sort.Slice(v, func(i, j int) bool {
			return v[i].FieldA < v[j].FieldA
		})
		sort.Slice(tc.expected, func(i, j int) bool {
			return tc.expected[i].FieldA < tc.expected[j].FieldA
		})

		if !reflect.DeepEqual(v, tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}
	}
}

// TestSortedWithContext tests SortedWithContext function.
func TestSortedWithContext(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{2, "two"},
				{1, "one"},
			},
			expected: []complexType{
				{1, "one"},
				{2, "two"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		v, _ := SortedWithContext(ctx, s, func(a, b complexType) bool {
			return a.FieldA < b.FieldA
		})
		if !reflect.DeepEqual(v, tc.expected) {
			t.Errorf("Test %s (%v): expected %v, but got %v",
				tc.name, s, tc.expected, v)
		}

		cancel()
		v, _ = SortedWithContext(ctx, s)
		if len(v) != 0 {
			t.Errorf("Test %s (%v): expected %v, but got %v",
				tc.name, s, []complexType{}, v)
		}
	}
}

// TestSorted tests Sorted function.
func TestSorted(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{2, "two"},
				{1, "one"},
			},
			expected: []complexType{
				{1, "one"},
				{2, "two"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		v := Sorted(s, func(a, b complexType) bool {
			return a.FieldA < b.FieldA
		})
		if !reflect.DeepEqual(v, tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}
	}
}

// TestFilteredWithContext tests FilteredWithContext function.
func TestFilteredWithContext(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		filter   func(complexType) bool
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			filter: func(item complexType) bool {
				return item.FieldA == 1
			},
			expected: []complexType{
				{1, "one"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		v, _ := FilteredWithContext(ctx, s, tc.filter)
		sort.Slice(v, func(i, j int) bool {
			return v[i].FieldA < v[j].FieldA
		})
		sort.Slice(tc.expected, func(i, j int) bool {
			return tc.expected[i].FieldA < tc.expected[j].FieldA
		})

		if !reflect.DeepEqual(v, tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}

		cancel()
		v, _ = FilteredWithContext(ctx, s, tc.filter)
		if len(v) != 0 {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, []complexType{}, v)
		}
	}
}

// TestFiltered tests Filtered function.
func TestFiltered(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		filter   func(complexType) bool
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			filter: func(item complexType) bool {
				return item.FieldA == 1
			},
			expected: []complexType{
				{1, "one"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		v := Filtered(s, tc.filter)
		sort.Slice(v, func(i, j int) bool {
			return v[i].FieldA < v[j].FieldA
		})
		sort.Slice(tc.expected, func(i, j int) bool {
			return tc.expected[i].FieldA < tc.expected[j].FieldA
		})

		if !reflect.DeepEqual(v, tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}
	}
}

// TestLen tests Len function.
func TestLen(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		expected int
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			expected: 2,
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		if v := Len(s); v != tc.expected {
			t.Errorf("Test %s: expected %d, but got %d",
				tc.name, tc.expected, v)
		}
	}
}

// UnionWithContext tests UnionWithContext function.
func TestUnionWithContext(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		input2   []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			input2: []complexType{
				{1, "one"},
				{3, "three"},
			},
			expected: []complexType{
				{1, "one"},
				{2, "two"},
				{3, "three"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		s2 := New(tc.input2...)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		v, _ := UnionWithContext(ctx, s, s2)
		if !reflect.DeepEqual(v.Sorted(), New(tc.expected...).Sorted()) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, New(tc.expected...).Sorted(), v.Sorted())
		}

		cancel()
		v, _ = UnionWithContext(ctx, s, s2)
		if Len(v) != 0 {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, []complexType{}, v.Sorted())
		}
	}
}

// TestUnion tests Union function.
func TestUnion(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		input2   []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			input2: []complexType{
				{1, "one"},
				{3, "three"},
			},
			expected: []complexType{
				{1, "one"},
				{2, "two"},
				{3, "three"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		s2 := New(tc.input2...)
		v := Union(s, s2)
		if !reflect.DeepEqual(v.Sorted(), New(tc.expected...).Sorted()) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, New(tc.expected...).Sorted(), v.Sorted())
		}
	}
}

// TestIntersectionWithContext tests IntersectionWithContext function.
func TestIntersectionWithContext(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		input2   []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			input2: []complexType{
				{1, "one"},
				{3, "three"},
			},
			expected: []complexType{
				{1, "one"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		s2 := New(tc.input2...)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		v, _ := InterWithContext(ctx, s, s2)
		if !reflect.DeepEqual(v.Sorted(), tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}

		cancel()
		v, _ = InterWithContext(ctx, s, s2)
		if Len(v) != 0 {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, []complexType{}, v)
		}
	}
}

// TestIntersection tests Intersection function.
func TestIntersection(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		input2   []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			input2: []complexType{
				{1, "one"},
				{3, "three"},
			},
			expected: []complexType{
				{1, "one"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		s2 := New(tc.input2...)
		v := Inter(s, s2)
		if !reflect.DeepEqual(v.Sorted(), tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}
	}
}

// TestDifferenceWithContext tests DifferenceWithContext function.
func TestDifferenceWithContext(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		input2   []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			input2: []complexType{
				{1, "one"},
				{3, "three"},
			},
			expected: []complexType{
				{2, "two"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		s2 := New(tc.input2...)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		v, _ := DiffWithContext(ctx, s, s2)
		if !reflect.DeepEqual(v.Sorted(), tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}

		cancel()
		v, _ = DiffWithContext(ctx, s, s2)
		if Len(v) != 0 {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, []complexType{}, v)
		}
	}
}

// TestDifference tests Difference function.
func TestDifference(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		input2   []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			input2: []complexType{
				{1, "one"},
				{3, "three"},
			},
			expected: []complexType{
				{2, "two"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		s2 := New(tc.input2...)
		v := Diff(s, s2)
		if !reflect.DeepEqual(v.Sorted(), tc.expected) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, tc.expected, v)
		}
	}
}

// TestSymmetricDifferenceWithContext tests SymmetricDifferenceWithContext fn.
func TestSymmetricDifferenceWithContext(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		input2   []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			input2: []complexType{
				{1, "one"},
				{3, "three"},
			},
			expected: []complexType{
				{2, "two"},
				{3, "three"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		s2 := New(tc.input2...)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		v, _ := SdiffWithContext(ctx, s, s2)
		if !reflect.DeepEqual(v.Sorted(), New(tc.expected...).Sorted()) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, New(tc.expected...).Sorted(), v.Sorted())
		}

		cancel()
		v, _ = SdiffWithContext(ctx, s, s2)
		if Len(v) != 0 {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, []complexType{}, v)
		}
	}
}

// TestSymmetricDifference tests SymmetricDifference function.
func TestSymmetricDifference(t *testing.T) {
	tests := []struct {
		name     string
		input    []complexType
		input2   []complexType
		expected []complexType
	}{
		{
			name: "one",
			input: []complexType{
				{1, "one"},
				{2, "two"},
			},
			input2: []complexType{
				{1, "one"},
				{3, "three"},
			},
			expected: []complexType{
				{2, "two"},
				{3, "three"},
			},
		},
	}

	for _, tc := range tests {
		s := New(tc.input...)
		s2 := New(tc.input2...)
		v := Sdiff(s, s2)
		if !reflect.DeepEqual(v.Sorted(), New(tc.expected...).Sorted()) {
			t.Errorf("Test %s: expected %v, but got %v",
				tc.name, New(tc.expected...).Sorted(), v.Sorted())
		}
	}
}

// TestMapWithContext tests MapWithContext function.
func TestMapWithContext(t *testing.T) {
	s := New[userType]()
	s.Add(userType{"John", 20}, userType{"Jane", 30})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	names, _ := MapWithContext(ctx, s, func(item userType) string {
		return item.Name
	})

	expected := []string{"Jane", "John"}
	if v := names.Sorted(); !reflect.DeepEqual(v, expected) {
		t.Errorf("Map() failed, expected names = %v, got %v",
			expected, v)
	}

	cancel()
	names, _ = MapWithContext(ctx, s, func(item userType) string {
		return item.Name
	})

	if Len(names) != 0 {
		t.Errorf("Map() failed, expected names = %v, got %v",
			[]string{}, names)
	}
}

// TestMap tests Map function.
func TestMap(t *testing.T) {
	s := New[userType]()
	s.Add(userType{"John", 20}, userType{"Jane", 30})

	names := Map(s, func(item userType) string {
		return item.Name
	})

	expected := []string{"Jane", "John"}
	if v := names.Sorted(); !reflect.DeepEqual(v, expected) {
		t.Errorf("Map() failed, expected names = %v, got %v",
			expected, v)
	}
}

// TestReduceWithContext tests ReduceWithContext function.
func TestReduceWithContext(t *testing.T) {
	s := New[userType]()
	s.Add(userType{"John", 20}, userType{"Jane", 30})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sum, _ := ReduceWithContext(ctx, s, func(acc int, item userType) int {
		return acc + item.Age
	})

	if sum != 50 {
		t.Errorf("Reduce() failed, expected sum = %d, got %d",
			50, sum)
	}

	cancel()
	sum, _ = ReduceWithContext(ctx, s, func(acc int, item userType) int {
		return acc + item.Age
	})

	if sum != 0 {
		t.Errorf("Reduce() failed, expected sum = %d, got %d",
			0, sum)
	}
}

// TestReduce tests Reduce function.
func TestReduce(t *testing.T) {
	s := New[userType]()
	s.Add(userType{"John", 20}, userType{"Jane", 30})

	sum := Reduce(s, func(acc int, item userType) int {
		return acc + item.Age
	})

	if sum != 50 {
		t.Errorf("Reduce() failed, expected sum = %d, got %d",
			50, sum)
	}
}

// TestCopyWithContext tests CopyWithContext function.
func TestCopyWithContext(t *testing.T) {
	s := New[userType]()
	s.Add(userType{"John", 20}, userType{"Jane", 30})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s2, _ := CopyWithContext(ctx, s)

	if !reflect.DeepEqual(s, s2) {
		t.Errorf("Copy() failed, expected s = %v, got %v",
			s, s2)
	}

	cancel()
	s2, _ = CopyWithContext(ctx, s)

	if Len(s2) != 0 {
		t.Errorf("Copy() failed, expected s = %v, got %v",
			New[userType](), s2)
	}
}

// TestCopy tests Copy function.
func TestCopy(t *testing.T) {
	s := New[userType]()
	s.Add(userType{"John", 20}, userType{"Jane", 30})

	s2 := Copy(s)

	if !reflect.DeepEqual(s, s2) {
		t.Errorf("Copy() failed, expected s = %v, got %v",
			s, s2)
	}
}

// TestFilterWithContext tests FilterWithContext function.
func TestFilterWithContext(t *testing.T) {
	s := New[userType]()
	s.Add(userType{"John", 20}, userType{"Jane", 30})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s2, _ := FilterWithContext(ctx, s, func(item userType) bool {
		return item.Age > 20
	})

	expected := New[userType]()
	expected.Add(userType{"Jane", 30})

	if !reflect.DeepEqual(s2, expected) {
		t.Errorf("Filter() failed, expected s = %v, got %v",
			expected, s2)
	}

	cancel()
	s2, _ = FilterWithContext(ctx, s, func(item userType) bool {
		return item.Age > 20
	})

	if Len(s2) != 0 {
		t.Errorf("Filter() failed, expected s = %v, got %v",
			New[userType](), s2)
	}
}

// TestFilter tests Filter function.
func TestFilter(t *testing.T) {
	s := New[userType]()
	s.Add(userType{"John", 20}, userType{"Jane", 30})

	s2 := Filter(s, func(item userType) bool {
		return item.Age > 20
	})

	expected := New[userType]()
	expected.Add(userType{"Jane", 30})

	if !reflect.DeepEqual(s2, expected) {
		t.Errorf("Filter() failed, expected s = %v, got %v",
			expected, s2)
	}
}
