package set_test

import (
	"fmt"

	"github.com/goloop/set/v2"
)

func ExampleNew() {
	s := set.New(1, 2, 2, 3, 3, 3)
	fmt.Println(s.Len())
	fmt.Println(set.Sorted(s))
	// Output:
	// 3
	// [1 2 3]
}

func ExampleSet_Union() {
	a := set.New(1, 2, 3)
	b := set.New(3, 4, 5)
	fmt.Println(set.Sorted(a.Union(b)))
	// Output: [1 2 3 4 5]
}

func ExampleSet_Intersection() {
	a := set.New(1, 2, 3, 4)
	b := set.New(3, 4, 5, 6)
	fmt.Println(set.Sorted(a.Intersection(b)))
	// Output: [3 4]
}

func ExampleSet_Difference() {
	a := set.New(1, 2, 3)
	b := set.New(3, 4, 5)
	fmt.Println(set.Sorted(a.Difference(b)))
	// Output: [1 2]
}

func ExampleSet_SymmetricDifference() {
	a := set.New(1, 2, 3)
	b := set.New(3, 4, 5)
	fmt.Println(set.Sorted(a.SymmetricDifference(b)))
	// Output: [1 2 4 5]
}

func ExampleSet_IsSubset() {
	small := set.New(1, 2)
	big := set.New(1, 2, 3)
	fmt.Println(small.IsSubset(big))       // ⊆ holds
	fmt.Println(big.IsSubset(big))         // a set is a subset of itself
	fmt.Println(small.IsProperSubset(big)) // strict ⊊
	fmt.Println(big.IsProperSubset(big))   // not strictly a subset of itself
	// Output:
	// true
	// true
	// true
	// false
}

func ExampleSet_Filter() {
	s := set.New(1, 2, 3, 4, 5)
	even := s.Filter(func(v int) bool { return v%2 == 0 })
	fmt.Println(set.Sorted(even))
	// Output: [2 4]
}

func ExampleMap() {
	type user struct {
		Name string
		Age  int
	}
	s := set.New(user{"John", 20}, user{"Jane", 30})
	names := set.Map(s, func(u user) string { return u.Name })
	fmt.Println(set.Sorted(names))
	// Output: [Jane John]
}

func ExampleFold() {
	s := set.New(2, 3, 4)
	product := set.Fold(s, 1, func(acc, v int) int { return acc * v })
	fmt.Println(product)
	// Output: 24
}

func ExampleSet_Iter() {
	s := set.New(1, 2, 3)
	sum := 0
	for v := range s.Iter() {
		sum += v
	}
	fmt.Println(sum)
	// Output: 6
}
