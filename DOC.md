# set — reference

The full reference for the `set` package: the mental model, element identity,
construction, the full method surface, the package-level generics, iteration,
JSON and practical recipes.

Ukrainian version: **[DOC.UK.md](DOC.UK.md)**.

## Contents

- [Mental model](#mental-model)
- [Element identity](#element-identity)
- [Construction](#construction)
- [Basic operations](#basic-operations)
- [Set algebra](#set-algebra)
- [Relations](#relations)
- [Iteration and ordering](#iteration-and-ordering)
- [Functional helpers](#functional-helpers)
- [JSON](#json)
- [Concurrency](#concurrency)
- [Recipes and tips](#recipes-and-tips)

## Mental model

`set` is a small, fast, generic `Set` for Go: an unordered collection of unique
elements of a `comparable` type, built directly on the built-in map.

The defining choice is that **an element's identity is the language's own
equality (`==`)**. Two elements are the same if and only if they compare equal,
and the runtime map decides uniqueness — there is no hashing, no reflection and
no custom equality contract. As a result:

- A `Set` can never silently drop an element to a hash collision.
- `Len` always reflects the true number of distinct elements.
- The behaviour matches your intuition about `==` for the element type.

```go
import "github.com/goloop/set/v2"
```

## Element identity

The `comparable` constraint admits the numeric kinds, `string`, `bool`,
pointers, channels, interfaces, and any struct or array whose fields are
themselves comparable.

Identity is `==`, so structs are compared field by field, by value:

```go
type Address struct{ City string }
type User struct{ Name string; Age int; Address Address }

users := set.New(
    User{"John", 21, Address{"Kyiv"}},
    User{"Bob", 22, Address{"Chernihiv"}},
    User{"John", 21, Address{"Kyiv"}}, // duplicate -> collapses
)
users.Len() // 2
```

A struct that holds a **pointer** is compared by that pointer, not by the value
it points to: two structs with different pointers are different elements even if
the pointed-to values are equal. Store values rather than pointers when you want
value identity.

Slices, maps and functions are not comparable and cannot be elements directly.
To deduplicate such values, derive a comparable key (a `string`, or a struct of
comparable fields) and build a `Set` of that key.

## Construction

```go
func New[T comparable](items ...T) *Set[T]
func NewWithCapacity[T comparable](capacity int, items ...T) *Set[T]
func Collect[T comparable](seq iter.Seq[T]) *Set[T]
```

`New` infers the element type from its arguments, or takes it explicitly when
empty (`set.New[int]()`). `NewWithCapacity` preallocates the backing map when
the size is known. `Collect` builds a set from any `iter.Seq[T]`.

The zero value is usable directly — `var s set.Set[int]` is an empty,
ready-to-use set (the first insertion allocates the backing map). `New` is still
preferred when the size is known.

```go
ints  := set.New[int]()
words := set.New("one", "two", "three")
keys  := set.Collect(maps.Keys(m))
var s set.Set[int] // also valid
```

## Basic operations

```go
func (s *Set[T]) Add(items ...T)
func (s *Set[T]) AddSeq(seq iter.Seq[T])
func (s *Set[T]) Delete(items ...T)
func (s *Set[T]) Contains(item T) bool
func (s *Set[T]) ContainsAll(items ...T) bool
func (s *Set[T]) ContainsAny(items ...T) bool
func (s *Set[T]) Len() int
func (s *Set[T]) IsEmpty() bool
func (s *Set[T]) Clear()
func (s *Set[T]) Copy() *Set[T]
func (s *Set[T]) Pop() (T, bool)
func (s *Set[T]) Elements() []T
func (s *Set[T]) Append(others ...*Set[T])
func (s *Set[T]) Overwrite(items ...T)
```

`Add`/`Delete` are variadic. `AddSeq` adds all values from an `iter.Seq[T]`.
`Pop` removes and returns an arbitrary element (`ok=false` when empty).
`Elements` returns the members as an unordered slice. `Append` merges other sets
into this one in place; `Overwrite` replaces the contents with the given items.

```go
ints.Add(1, 2, 3, 4)
ints.Delete(1, 2)
ints.Contains(3)          // true
ints.ContainsAll(3, 4)    // true
s.AddSeq(slices.Values(items))
```

## Set algebra

```go
func (s *Set[T]) Union(others ...*Set[T]) *Set[T]
func (s *Set[T]) Intersection(others ...*Set[T]) *Set[T] // alias: Inter
func (s *Set[T]) Difference(others ...*Set[T]) *Set[T]   // alias: Diff
func (s *Set[T]) SymmetricDifference(others ...*Set[T]) *Set[T] // alias: Sdiff
```

Each returns a new set and accepts several operands at once
(`a.Union(b, c, d)`):

```go
a := set.New(1, 3, 5, 7)
b := set.New(0, 2, 4, 7)

set.Sorted(a.Union(b))               // [0 1 2 3 4 5 7]
set.Sorted(a.Intersection(b))        // [7]
set.Sorted(a.Difference(b))          // [1 3 5]
set.Sorted(a.SymmetricDifference(b)) // [0 1 2 3 4 5]
```

## Relations

```go
func (s *Set[T]) Equal(other *Set[T]) bool
func (s *Set[T]) IsSubset(other *Set[T]) bool         // alias: IsSub
func (s *Set[T]) IsSuperset(other *Set[T]) bool       // alias: IsSup
func (s *Set[T]) IsProperSubset(other *Set[T]) bool
func (s *Set[T]) IsProperSuperset(other *Set[T]) bool
func (s *Set[T]) IsDisjoint(other *Set[T]) bool
```

`IsSubset`/`IsSuperset` are the **non-strict** relations — a set is a subset and
superset of itself, matching the standard mathematical definitions. Use the
`Proper` variants for the strict relations.

```go
a := set.New(1, 2, 3)
b := set.New(1, 2, 3, 4, 5)

a.IsSubset(b)               // true  (a ⊆ b)
a.IsProperSubset(b)         // true  (a ⊊ b)
b.IsSuperset(a)             // true  (b ⊇ a)
a.Equal(set.New(3, 2, 1))   // true  (order does not matter)
a.IsDisjoint(set.New(8, 9)) // true
```

## Iteration and ordering

The iteration order of a set is **unspecified**.

```go
func (s *Set[T]) Iter() iter.Seq[T]
func (s *Set[T]) Elements() []T
func (s *Set[T]) Sorted(cmp func(a, b T) int) []T
func Sorted[T cmp.Ordered](s *Set[T]) []T
```

Use `Iter` for a `range` loop, `Elements` for an unordered slice, or `Sorted`
when you need a stable order:

```go
for v := range s.Iter() {
    _ = v // unspecified order
}

set.Sorted(s)                                 // [1 2 3]  (natural order)
s.Sorted(func(a, b int) int { return b - a }) // [3 2 1]  (custom order)
```

The package-level `set.Sorted` works for `cmp.Ordered` element types with no
argument; the `Sorted` method takes a comparison function (the same contract as
`cmp.Compare`) for any other order.

## Functional helpers

Methods that keep the element type, plus package-level generics that may change
it:

```go
func (s *Set[T]) Map(fn func(item T) T) *Set[T]
func (s *Set[T]) Filter(fn func(item T) bool) *Set[T]
func (s *Set[T]) Filtered(fn func(item T) bool) []T
func (s *Set[T]) Reduce(fn func(acc, item T) T) T
func (s *Set[T]) Any(fn func(item T) bool) bool
func (s *Set[T]) All(fn func(item T) bool) bool

func Map[T, R comparable](s *Set[T], fn func(item T) R) *Set[R]
func Reduce[T comparable, R any](s *Set[T], fn func(acc R, item T) R) R
func Fold[T comparable, R any](s *Set[T], initial R, fn func(acc R, item T) R) R
```

```go
s := set.New(1, 2, 3, 4, 5)

even := s.Filter(func(v int) bool { return v%2 == 0 }) // {2, 4}
doubled := s.Map(func(v int) int { return v * 2 })     // same element type

// The Map function may change the element type.
labels := set.Map(s, func(v int) string {
    if v%2 == 0 { return "even" }
    return "odd"
}) // {"odd", "even"}

sum := s.Reduce(func(acc, v int) int { return acc + v })           // 15
product := set.Fold(s, 1, func(acc, v int) int { return acc * v }) // 120
```

`Reduce` (method) starts from the zero value; `Fold` takes an explicit start and
may accumulate into a different type. `Any`/`All` are simple linear scans.

## JSON

`Set` implements the standard `encoding/json` interfaces:

```go
s := set.New(1, 2, 3)

data, _ := json.Marshal(s) // e.g. [1,2,3] (order unspecified)

var back set.Set[int]
_ = json.Unmarshal(data, &back)
back.Equal(s) // true
```

A set marshals to a JSON array and unmarshals from one, deduplicating on the
way in.

## Concurrency

A `Set` is **not** safe for concurrent use by multiple goroutines, exactly like
the built-in map it is built on. If a set is shared and at least one goroutine
mutates it, guard access with your own synchronization:

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
caller, who knows the application's concurrency model, choose the right strategy.

## Recipes and tips

**Deduplicate a slice.** `set.Collect(slices.Values(xs)).Elements()` (or
`set.Sorted(...)`) returns the distinct values.

**Key non-comparable values.** For slices/maps, derive a comparable key
(`fmt.Sprint`, a hash string, or a struct of comparable fields) and store a
`Set` of that key.

**Value identity for structs.** Store struct values, not pointers, so equal
records collapse to one element.

**Stable output.** The iteration order is unspecified — pass results through
`set.Sorted` (or the `Sorted` method) whenever the order is observable, e.g. in
tests or serialized output.

**Preallocate when you know the size.** `NewWithCapacity(n, …)` avoids repeated
map growth when adding many elements.
