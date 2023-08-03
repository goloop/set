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
// Import the package
import (
	"github.com/goloop/set"
)

...

// Define a new set
s := set.New[int]()

// Add elements to the set.
s.Add(1, 2, 3, 4)

// Check if the set contains a certain element.
containsThree := s.Contains(3)  // returns true
containsFive := s.Contains(5)   // returns false

// Remove an element from the set.
s.Delete(3)

// Check the size of the set.
size := s.Len() // returns 3
```

Documentation
You can read more about the Set package and its functions on [Godoc](https://godoc.org/github.com/goloop/set).


