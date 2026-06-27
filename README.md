[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/set)](https://goreportcard.com/report/github.com/goloop/set) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/set/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/set) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# set

A small, fast, generic Set for Go: an unordered collection of unique elements
of a comparable type, built directly on the built-in map.

An element's identity is the language's own equality (`==`). Two elements are
the same if and only if they compare equal, and the runtime map decides
uniqueness — there is no hashing, no reflection and no custom equality
contract. As a result a Set can never silently drop an element to a collision,
and `Len` always reflects the true number of distinct elements.

It provides the basic operations (`Add`, `Delete`, `Contains`, `Len`), the set
algebra (`Union`, `Intersection`, `Difference`, `SymmetricDifference`), the
relations (`Equal`, `IsSubset`, `IsSuperset`, `IsProperSubset`,
`IsProperSuperset`, `IsDisjoint`), and functional helpers (`Map`, `Filter`,
`Reduce`, `Fold`, `Any`, `All`), plus iteration via `iter.Seq` and JSON
serialization.

## Features

- Generic over any `comparable` element type
- Exact `==` identity — no hashing, no lost elements
- Full set algebra and relation predicates
- Functional helpers: `Map`, `Filter`, `Reduce`, `Fold`, `Any`, `All`
- `iter.Seq[T]` iteration for `range`, plus `AddSeq` / `Collect`
- Usable zero value: `var s set.Set[int]` is an empty, ready-to-use set
- JSON serialization through the standard `encoding/json` interfaces
- Zero dependencies

## Element types

The `comparable` constraint admits the numeric kinds, `string`, `bool`,
pointers, channels, interfaces, and any struct or array whose fields are
themselves comparable.

Identity is `==`, so structs are compared field by field by value. A struct
that holds a pointer is compared by that pointer, not by the value it points
to: two structs with different pointers are different elements even if the
pointed-to values are equal. When you want value identity, store values rather
than pointers.

Slices, maps and functions are not comparable and cannot be elements directly.
To deduplicate such values, derive a comparable key from them (for example a
`string`, or a struct of comparable fields) and build a `Set` of that key.

## Concurrency

A `Set` is **not** safe for concurrent use by multiple goroutines, exactly
like the built-in map it is built on. If a `Set` is shared and at least one
goroutine mutates it, guard access with your own synchronization:

```go
type SafeSet[T comparable] struct {
    mu sync.RWMutex
    s  *set.Set[T]
}

func (c *SafeSet[T]) Add(items ...T) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.s.Add(items...)
}
```

Keeping the core unsynchronized avoids per-operation locking and lets the
caller, who knows the application's concurrency model, choose the right
strategy.

## Installation

```shell
go get github.com/goloop/set/v2
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/goloop/set/v2"
)

func main() {
	// Empty set: the element type must be given explicitly.
	ints := set.New[int]()

	// Or infer the type from the initial elements.
	words := set.New("one", "two", "three")

	ints.Add(1, 2, 3, 4)
	words.Add("three", "four") // "three" is already present

	fmt.Println(ints.Len())        // 4
	fmt.Println(set.Sorted(ints))  // [1 2 3 4]
	fmt.Println(set.Sorted(words)) // [four one three two]

	fmt.Println(ints.Contains(3), ints.Contains(10)) // true false

	ints.Delete(1, 2)
	fmt.Println(set.Sorted(ints)) // [3 4]
}
```

### Complex (value) types

```go
package main

import (
	"fmt"

	"github.com/goloop/set/v2"
)

// All fields are comparable values, so equal users collapse to one element.
type Address struct {
	City string
}

type User struct {
	Name    string
	Age     int
	Address Address
}

func main() {
	users := set.New(
		User{"John", 21, Address{"Kyiv"}},
		User{"Bob", 22, Address{"Chernihiv"}},
		User{"John", 21, Address{"Kyiv"}}, // duplicate -> collapses
	)

	fmt.Println(users.Len()) // 2
	fmt.Println(users.Contains(User{"John", 21, Address{"Kyiv"}})) // true
}
```

### Set algebra

```go
a := set.New(1, 3, 5, 7)
b := set.New(0, 2, 4, 7)

set.Sorted(a.Union(b))               // [0 1 2 3 4 5 7]
set.Sorted(a.Intersection(b))        // [7]
set.Sorted(a.Difference(b))          // [1 3 5]
set.Sorted(a.SymmetricDifference(b)) // [0 1 2 3 4 5]
```

All four accept several sets at once, e.g. `a.Union(b, c, d)`.

### Relations

```go
a := set.New(1, 2, 3)
b := set.New(1, 2, 3, 4, 5)

a.IsSubset(b)        // true  (a ⊆ b)
a.IsProperSubset(b)  // true  (a ⊊ b)
b.IsSuperset(a)      // true  (b ⊇ a)
a.Equal(set.New(3, 2, 1)) // true (order does not matter)
a.IsDisjoint(set.New(8, 9)) // true
```

A set is a (non-proper) subset and superset of itself, matching the standard
mathematical definitions.

### Iteration and ordering

The iteration order of a set is unspecified. Use `Iter` for a `range` loop,
`Elements` for an unordered slice, or `Sorted` when you need a stable order.

```go
s := set.New(3, 1, 2)

for v := range s.Iter() {
	_ = v // visited in unspecified order
}

set.Sorted(s)                                 // [1 2 3]  (natural order)
s.Sorted(func(a, b int) int { return b - a }) // [3 2 1]  (custom order)
```

Build a set from any `iter.Seq[T]` with `Collect`, or feed one into an
existing set with `AddSeq`:

```go
keys := set.Collect(maps.Keys(m)) // set of the map's keys
s.AddSeq(slices.Values(items))    // add all values from a slice
```

The package-level `set.Sorted` works for element types that satisfy
`cmp.Ordered`; the `Sorted` method takes a comparison function (the same
contract as `cmp.Compare`) for any other order.

### Functional helpers

```go
s := set.New(1, 2, 3, 4, 5)

even := s.Filter(func(v int) bool { return v%2 == 0 }) // {2, 4}

// The Map method keeps the element type.
doubled := s.Map(func(v int) int { return v * 2 }) // {2, 4, 6, 8, 10}

// The Map function may change the element type.
labels := set.Map(s, func(v int) string {
	if v%2 == 0 {
		return "even"
	}
	return "odd"
}) // {"odd", "even"}

// Reduce starts from the zero value; Fold takes an explicit start.
sum := s.Reduce(func(acc, v int) int { return acc + v })          // 15
product := set.Fold(s, 1, func(acc, v int) int { return acc * v }) // 120

s.Any(func(v int) bool { return v > 4 }) // true
s.All(func(v int) bool { return v > 0 }) // true
```

### JSON

```go
s := set.New(1, 2, 3)

data, _ := json.Marshal(s) // e.g. [1,2,3] (order unspecified)

var back set.Set[int]
_ = json.Unmarshal(data, &back)
back.Equal(s) // true
```

## Migrating from v1

Version 2 is a deliberate clean break:

- **Import path** is now `github.com/goloop/set/v2`.
- **Elements must be `comparable`** (`Set[T comparable]`). The previous
  reflection-based hashing is gone, together with its silent-collision data
  loss. Slice/map elements are no longer supported directly — key them by a
  comparable value instead.
- **Not thread-safe.** The internal mutex and the "thread-safe" guarantee were
  removed; synchronize externally if needed.
- **No context API.** The `*WithContext` variants and the package-function
  duplicates were removed; operations on an in-memory set do not block.
- **No global `ParallelTasks`.** `Any`/`All` are simple linear scans.
- `IsSubset`/`IsSuperset` are now the non-strict relations (a set is a subset
  of itself); use `IsProperSubset`/`IsProperSuperset` for the strict ones.
- `Sorted` takes a single comparison function returning `int`
  (`func(a, b T) int`); the package-level `set.Sorted` sorts `cmp.Ordered`
  types with no argument.
- New: `Equal`, `IsProperSubset`, `IsProperSuperset`, `IsDisjoint`, `IsEmpty`,
  `ContainsAll`, `ContainsAny`, `Pop`, `Iter`, `AddSeq`, and the `Fold` and
  `Collect` functions.
- The zero value `var s set.Set[T]` is now usable directly (the first
  insertion allocates the backing map); `New` is still preferred when the size
  is known.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
