[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/set/v2)](https://goreportcard.com/report/github.com/goloop/set/v2) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/set/blob/master/LICENSE) [![Go Reference](https://pkg.go.dev/badge/github.com/goloop/set/v2.svg)](https://pkg.go.dev/github.com/goloop/set/v2) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)

# set

`set` is a small, fast, generic Set for Go: an unordered collection of unique
elements of a `comparable` type, built directly on the built-in map.

An element's identity is the language's own equality (`==`). Two elements are
the same if and only if they compare equal, and the runtime map decides
uniqueness — there is no hashing, no reflection and no custom equality contract.
As a result a Set can never silently drop an element to a collision, and `Len`
always reflects the true number of distinct elements.

## Features

- Generic over any `comparable` element type; exact `==` identity.
- Full set algebra (`Union`, `Intersection`, `Difference`,
  `SymmetricDifference`) and relations (`Equal`, `IsSubset`, `IsSuperset`,
  `IsProperSubset`, `IsProperSuperset`, `IsDisjoint`).
- Functional helpers: `Map`, `Filter`, `Reduce`, `Fold`, `Any`, `All`.
- `iter.Seq[T]` iteration for `range`, plus `AddSeq` / `Collect`.
- Usable zero value: `var s set.Set[int]` is an empty, ready-to-use set.
- JSON serialization through the standard `encoding/json` interfaces.
- Zero dependencies.

## Installation

```shell
go get github.com/goloop/set/v2
```

```go
import "github.com/goloop/set/v2"
```

Requires Go 1.24 or newer. The package has no third-party dependencies.

## Quick start

```go
package main

import (
    "fmt"

    "github.com/goloop/set/v2"
)

func main() {
    ints := set.New[int]()          // empty; element type explicit
    words := set.New("one", "two")  // or inferred from the elements

    ints.Add(1, 2, 3, 4)
    words.Add("two", "three")       // "two" is already present

    fmt.Println(ints.Len())         // 4
    fmt.Println(set.Sorted(ints))   // [1 2 3 4]
    fmt.Println(ints.Contains(3))   // true

    // Set algebra (each accepts several sets at once).
    a := set.New(1, 3, 5, 7)
    b := set.New(0, 2, 4, 7)
    fmt.Println(set.Sorted(a.Intersection(b))) // [7]
    fmt.Println(a.IsDisjoint(set.New(8, 9)))   // true
}
```

> A `Set` is **not** safe for concurrent use, like the built-in map it is built
> on — synchronize externally if a shared set is mutated from several goroutines.

## Documentation

- Full reference and recipes: [DOC.md](DOC.md) · [DOC.UK.md](DOC.UK.md)
- Package API: [pkg.go.dev/github.com/goloop/set/v2](https://pkg.go.dev/github.com/goloop/set/v2)
- Changes between versions: [CHANGELOG.md](CHANGELOG.md)

## Contributing

Contributions are welcome. Please run `go test ./...`, `go vet ./...` and
`gofmt -l .` before submitting a pull request.

## License

`set` is released under the MIT License. See [LICENSE](LICENSE).
