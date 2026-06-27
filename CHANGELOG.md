# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0]

A complete redesign. The element model, the concurrency contract and the API
surface all changed; version 2 is not source-compatible with version 1. The
import path is now `github.com/goloop/set/v2`.

### Changed
- `Set[T any]` is now `Set[T comparable]`. Element identity is the language's
  own equality (`==`); the runtime map decides uniqueness.
- `IsSubset`/`IsSuperset` are the non-strict relations: a set is a subset and
  a superset of itself.
- `Sorted` (method) takes a single comparison function `func(a, b T) int`
  following the `cmp.Compare` contract; the new package-level `Sorted` sorts
  `cmp.Ordered` element types with no argument.
- `Union`, `Intersection`, `Difference` and `SymmetricDifference` are variadic
  methods returning a new set.

### Added
- `Equal`, `IsProperSubset`, `IsProperSuperset`, `IsDisjoint`.
- `IsEmpty`, `ContainsAll`, `ContainsAny`, `Pop`.
- `Iter` returning `iter.Seq[T]` for range iteration.
- `Fold` (package function) with an explicit initial accumulator value.
- `NewWithCapacity` for pre-sizing.

### Removed
- The reflection-based hashing of elements, and with it the possibility of
  silently losing distinct elements to a hash collision.
- Direct support for non-comparable element types (slices, maps, functions);
  key such values by a comparable type instead.
- The internal mutex and the "thread-safe" guarantee. A `Set` is not safe for
  concurrent use; synchronize externally if needed.
- The `*WithContext` methods and the duplicated package-level wrappers.
- The global `ParallelTasks` setting and the parallel execution of `Any`/`All`.

### Fixed
- Distinct composite values are no longer collapsed (former collisions of
  field-swapped structs and equal-length slices).
- `All` over an empty set now returns `true` (vacuous truth).
- `IsSubset` returns `true` for equal sets, consistent with `IsSuperset`.
- Data races on the element map and on the former global parallelism setting
  are gone by construction.
