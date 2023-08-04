[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/set?v1)](https://goreportcard.com/report/github.com/goloop/set) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/set/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/set) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# Set

Set is a Go package that provides a parameterized Set data structure.

A Set can contain any type of object, including simple and complex types. It provides basic set operations, such as `Add`, `Delete`, `Contains`, and `Len`. In addition, it also provides complex set operations, such as `Union`, `Intersection`, `Difference`, `SymmetricDifference`, `IsSubset`, and `IsSuperset`.

## Installation

You can download and install the package by running:

```shell
go get github.com/goloop/set
```

## Usage

Here's a simple example of how to use the Set package.

### Basic functions

```go
package main

import (
	"fmt"

	"github.com/goloop/set"
)

// Address is user's address.
type Address struct {
	City string
}

// User is the user type.
type User struct {
	Name    string   // a simple field of string type
	Age     int      // a simple field of int type
	Address *Address // nested fields
}

func main() {
	// Define a new sets.
	// Empty set, the type of set must be specified.
	iSet := set.New[int]()

	// A type can be defined from the elements of a set.
	sSet := set.New("one", "two", "three")

	// Support for complex types such as slice, arrays, maps, structures.
	cSet := set.New[User](
		User{
			Name: "John",
			Age:  21,
			Address: &Address{
				City: "Kyiv",
			},
		},
		User{
			Name: "Bob",
			Age:  22,
			Address: &Address{
				City: "Chernihiv",
			},
		},
	)

	// Add elements to the set.
	iSet.Add(1, 2, 3, 4)
	sSet.Add("three", "four", "five")
	cSet.Add(
		User{
			Name: "John",
			Age:  21,
			Address: &Address{
				City: "Kyiv",
			},
		}, // duplicate
		User{
			Name: "Victoria",
			Age:  23,
			Address: &Address{
				City: "Chernihiv",
			},
		},
	)

	// Check the size of the set.
	fmt.Println("\nLength:")
	fmt.Println("iSet: ", iSet.Len())
	fmt.Println("sSet: ", sSet.Len())
	fmt.Println("cSet: ", cSet.Len())

	// Elements.
	fmt.Println("\nElements:")
	fmt.Println("iSet: ", iSet.Elements())
	fmt.Println("sSet: ", sSet.Elements())
	fmt.Println("cSet: ", cSet.Elements())

	// Check if the set contains a certain element.
	fmt.Println("\nContains:")
	fmt.Println("iSet: ", iSet.Contains(3), iSet.Contains(10))
	fmt.Println("sSet: ", sSet.Contains("five"), sSet.Contains("nine"))
	fmt.Println("cSet: ",
		cSet.Contains(
			User{
				Name: "John",
				Age:  21,
				Address: &Address{
					City: "Kyiv",
				},
			}, // [!] new object here
		),
		cSet.Contains(
			User{
				Name: "Adam",
				Age:  23,
				Address: &Address{
					City: "Chernihiv",
				},
			},
		),
	)

	// Delete.
	iSet.Delete(1, 2, 4)
	sSet.Delete("four")
	cSet.Delete(
		User{
			Name: "John",
			Age:  21,
			Address: &Address{
				City: "Kyiv",
			},
		},
	)

	fmt.Println("\nElements after deletion:")
	fmt.Println("iSet: ", iSet.Elements())
	fmt.Println("sSet: ", sSet.Elements())
	fmt.Println("cSet: ", cSet.Elements())
}


// Output:
// Length:
// iSet:  4
// sSet:  5
// cSet:  3

// Elements:
// iSet:  [1 2 3 4]
// sSet:  [two three four five one]
// cSet:  [{John 21 0xc033} {Bob 22 0xc055} {Victoria 23 0xc077}]

// Contains:
// iSet:  true false
// sSet:  true false
// cSet:  true false

// Elements after deletion:
// iSet:  [3]
// sSet:  [one two three five]
// cSet:  [{Bob 22 0xc055} {Victoria 23 0xc077}]
```

See example [here](https://go.dev/play/p/1yjv4imgOiD).


### Operations on set

```go
package main

import (
	"fmt"

	"github.com/goloop/set"
)

func main() {
	a := set.New(1, 3, 5, 7)
	b := set.New(0, 2, 4, 7)

	// Union.
	c := a.Union(b)
	fmt.Println("Union:")
	fmt.Println("a: ", a.Elements())
	fmt.Println("b: ", b.Elements())
	fmt.Println("c: ", c.Elements())

	// Intersection.
	c = a.Intersection(b)
	fmt.Println("\nIntersection:")
	fmt.Println("a: ", a.Elements())
	fmt.Println("b: ", b.Elements())
	fmt.Println("c: ", c.Elements())

	// Difference.
	c = a.Difference(b)
	fmt.Println("\nDifference:")
	fmt.Println("a: ", a.Elements())
	fmt.Println("b: ", b.Elements())
	fmt.Println("c: ", c.Elements())

	// SymmetricDifference.
	c = a.SymmetricDifference(b)
	fmt.Println("\nSymmetricDifference:")
	fmt.Println("a: ", a.Elements())
	fmt.Println("b: ", b.Elements())
	fmt.Println("c: ", c.Elements())

	// IsSubset.
	a = set.New(1, 2, 3)
	b = set.New(1, 2, 3, 4, 5, 6)
	fmt.Println("\nSubset:")
	fmt.Println("a: ", a.Elements())
	fmt.Println("b: ", b.Elements())
	fmt.Println("result: ", a.IsSubset(b))

	// IsSuperset.
	a = set.New(1, 2, 3)
	b = set.New(1, 2, 3, 4, 5, 6)
	fmt.Println("\nIsSuperset:")
	fmt.Println("a: ", a.Elements())
	fmt.Println("b: ", b.Elements())
	fmt.Println("result: ", a.IsSuperset(b))

}

// Output:
// Union:
// a:  [3 5 7 1]
// b:  [0 2 4 7]
// c:  [7 2 4 0 1 3 5]

// Intersection:
// a:  [1 3 5 7]
// b:  [7 0 2 4]
// c:  [7]

// Difference:
// a:  [7 1 3 5]
// b:  [0 2 4 7]
// c:  [1 3 5]

// SymmetricDifference:
// a:  [1 3 5 7]
// b:  [0 2 4 7]
// c:  [4 3 5 1 0 2]

// Subset:
// a:  [3 1 2]
// b:  [1 2 3 4 5 6]
// result:  true

// IsSuperset:
// a:  [1 2 3]
// b:  [4 5 6 1 2 3]
// result:  false
```

See example [here](https://go.dev/play/p/pwv3PALIroT).

### Other functions

**Example 1.**

```go
package main

import (
	"fmt"

	"github.com/goloop/set"
)

// User object.
type User struct {
	Name string
	Age  int
}

func main() {
	// Simple set values.
	a := set.New(1, 3, 5, 7)
	b := set.New(0, 2, 4, 7)
	c := set.New(3, 4, 5, 6, 7, 8)

	// Update instance.
	a.Append(b, c) // a.Extend([]*set.Set[int]{b, c})
	fmt.Println("Append/Extend:", a.Elements())

	// Return sorted elements.
	fmt.Println("Sorted:", a.Sorted())
	fmt.Println("Reverse:", a.Sorted(func(a, b int) bool {
		return a > b
	}))

	// Clear.
	a.Clear()
	fmt.Println("Cleared:", a.Elements())

	// Complex set values.
	d := set.New(
		User{"John", 21},
		User{"Bob", 27},
		User{"Maya", 25},
	)

	fmt.Println("Sorted by Age:", d.Sorted(func(a, b User) bool {
		return a.Age < b.Age
	}))
	fmt.Println("Sorted by Name:", d.Sorted(func(a, b User) bool {
		return a.Name < b.Name
	}))
}

// Output:
// Append/Extend: [7 4 8 6 1 3 5 2 0]
// Sorted: [0 1 2 3 4 5 6 7 8]
// Reverse: [8 7 6 5 4 3 2 1 0]
// Cleared: []
// Sorted by Age: [{John 21} {Maya 25} {Bob 27}]
// Sorted by Name: [{Bob 27} {John 21} {Maya 25}]
```

See example [here](https://go.dev/play/p/WTiW6_viwrO).

**Example 2.**

```go
package main

import (
	"fmt"

	"github.com/goloop/set"
)

// User object.
type User struct {
	Name   string
	Age    int
	Gender string
}

func main() {
	s := set.New(
		User{"Alice", 20, "f"},
		User{"Bob", 30, "m"},
		User{"Charlie", 40, "m"},
		User{"Dave", 50, "m"},
		User{"Eve", 16, "f"},
	)

	// Filtered.
	fmt.Println("\nFiltered:")
	fmt.Println("Women:", s.Filtered(func(item User) bool {
		return item.Gender == "f"
	}))
	fmt.Println("Adults:", s.Filtered(func(item User) bool {
		return item.Age > 18
	}))

	// Filter.
	f := s.Filter(func(item User) bool {
		return item.Gender == "f" && item.Age > 18
	})
	fmt.Println("\nFilter:")
	fmt.Println("Adults women:", f.Elements())

	// Map.
	// Methods cannot support generics, so we need to use the set.Map
	// function to change the types of generated values.
	//
	// Better to use the Map method for simple types only, like:
	// int, uint, bool, etc.
	// s := set.New[int](1, 2, 3, 4)
	// m := s.Map(func(item int) int {
	//     return item * 2
	// }) // returns a new set with the values {2, 4, 6, 8}
	names := set.Map(s, func(item User) string {
		return item.Name
	})
	fmt.Println("\nMap:")
	fmt.Println("Names:", names.Elements())

	// Reduce.
	// Methods cannot support generics, so we need to use the set.Reduce
	// function to change the types of generated values.
	//
	// We can use the Reduce method for simple types only, like:
	// int, uint, bool, etc.
	// s := set.New[int](1, 2, 3, 4)
	// sum := s.Reduce(func(acc int, item int) int {
	//     return acc + item
	// }) // returns 10
	sum := set.Reduce(s, func(acc int, item User) int {
		return acc + item.Age
	})
	fmt.Println("\nReduce:")
	fmt.Println("Total age:", sum)

	// Any.
	fmt.Println("\nAny:")
	fmt.Println("Any adult:", s.Any(func(item User) bool {
		return item.Age > 18
	}))

	// All.
	fmt.Println("\nAll:")
	fmt.Println("All adults:", s.All(func(item User) bool {
		return item.Age > 18
	}))
}

// Output:
// Filtered:
// Women: [{Eve 16 f} {Alice 20 f}]
// Adults: [{Alice 20 f} {Bob 30 m} {Charlie 40 m} {Dave 50 m}]

// Filter:
// Adults women: [{Alice 20 f}]

// Map:
// Names: [Alice Bob Charlie Dave Eve]

// Reduce:
// Total age: 156

// Any:
// Any adult: true

// All:
// All adults: false
```

See example [here](https://go.dev/play/p/nfvIji29YhN).

## Documentation
You can read more about the Set package and its functions on [Godoc](https://godoc.org/github.com/goloop/set).


