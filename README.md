[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/set)](https://goreportcard.com/report/github.com/goloop/set) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/set/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/set) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# Set

Set is a Go package that provides a parameterized Set data structure.

A Set can contain any type of object, including simple and complex types. It provides basic set operations, such as `Add`, `Delete`, `Contains`, and `Len`. In addition, it also provides complex set operations, such as `Union`, `Intersection`, `Difference`, `SymmetricDifference`, `IsSubset`, and `IsSuperset`.

## Installation

You can download and install the package by running:

```shell
go get github.com/goloop/set
```

##Usage

Here's a simple example of how to use the Set package:

```go
// You can edit this code!
// Click here and start typing.
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
	fmt.Println("cSet: ", cSet.Contains(
		User{
			Name: "Natalia",
			Age:  23,
			Address: &Address{
				City: "Chernihiv",
			},
		}, // [!] new object here
	))

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
// sSet:  [one two three four five]
// cSet:  [{John 21 0xc033} {Bob 22 0xc055} {Victoria 23 0xc077}]

// Contains:
// iSet:  true false
// sSet:  true false
// cSet:  false

// Elements after deletion:
// iSet:  [3]
// sSet:  [three five one two]
// cSet:  [{Bob 22 0xc055} {Victoria 23 0xc077}]
```

See example [here](https://go.dev/play/p/YXrjFLcbxOo).


## Documentation
You can read more about the Set package and its functions on [Godoc](https://godoc.org/github.com/goloop/set).


