// Package set provides a thread-safe, generic Set data structure
// implementation for Go, supporting both simple and complex data types
// with rich functionality for set operations and concurrent processing
// capabilities.
//
// Core Features:
//   - Generic type support for any comparable type
//   - Thread-safe operations through sync.RWMutex
//   - Context-aware methods for cancellation support
//   - Efficient parallel processing for large datasets
//   - JSON serialization support
//   - Functional programming operations (Map, Filter, Reduce)
//
// Type System:
// The Set implementation distinguishes between simple and complex types:
//   - Simple types: basic Go types (int, string, bool, etc.)
//   - Complex types: structs, slices, maps, etc.
//
// Note: A single Set instance can only contain either simple or complex types,
// not both.
//
// Basic Operations:
//   - New[T](...T): Create a new Set with optional initial elements
//   - Add(...T): Add elements to the Set
//   - Delete(...T): Remove elements from the Set
//   - Contains(T): Check if an element exists in the Set
//   - Len(): Get the number of elements in the Set
//   - Clear(): Remove all elements from the Set
//
// Set Operations:
//   - Union: Combine elements from multiple sets
//   - Intersection: Find common elements between sets
//   - Difference: Find elements in one set but not in others
//   - SymmetricDifference: Find elements unique to each set
//   - IsSubset: Check if one set is contained within another
//   - IsSuperset: Check if one set contains another set
//
// Functional Operations:
//   - Map: Transform elements using a mapping function
//   - Filter: Select elements based on a predicate
//   - Reduce: Aggregate elements into a single value
//   - Any: Check if any element satisfies a condition
//   - All: Check if all elements satisfy a condition
//
// Concurrent Processing:
// The package automatically handles parallel processing for large datasets:
//   - Configurable number of parallel tasks (default: 2 * NumCPU)
//   - Automatic task distribution based on data size
//   - Minimum load threshold for parallel processing
//
// Context Support:
// Most operations have context-aware variants for cancellation support:
//   - AddWithContext
//   - DeleteWithContext
//   - UnionWithContext
//   - IntersectionWithContext
//   - etc.
//
// JSON Support:
// Sets can be serialized to and from JSON format:
//   - MarshalJSON(): Convert Set to JSON
//   - UnmarshalJSON(data []byte): Create Set from JSON
//
// Example usage:
//
//	// Creating sets
//	s1 := set.New[int](1, 2, 3)
//	s2 := set.New[int](3, 4, 5)
//
//	// Basic operations
//	s1.Add(4)
//	s1.Delete(1)
//	exists := s1.Contains(2)
//
//	// Set operations
//	union := s1.Union(s2)
//	intersection := s1.Intersection(s2)
//	difference := s1.Difference(s2)
//
//	// Using with context
//	ctx := context.Background()
//	if err := set.AddWithContext(ctx, s1, 6, 7); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Functional operations
//	even := s1.Filter(func(i int) bool {
//	    return i%2 == 0
//	})
//
//	doubled := set.Map(s1, func(i int) int {
//	    return i * 2
//	})
//
//	// Using with complex types
//	type User struct {
//	    ID   int
//	    Name string
//	}
//
//	users := set.New[User](
//	    User{1, "Alice"},
//	    User{2, "Bob"},
//	)
//
//	// JSON serialization
//	data, _ := users.MarshalJSON()
//	newUsers := set.New[User]()
//	newUsers.UnmarshalJSON(data)
//
// Performance Considerations:
//   - Parallel processing activates for datasets larger than minLoadPerGoroutine
//   - Thread-safety adds minimal overhead for normal operations
//   - Complex type operations may be slower due to reflection-based hashing
//   - Memory usage is optimized for the specific type being stored
//
// Thread Safety:
// All operations are thread-safe by default. The Set uses sync.RWMutex
// internally to ensure safe concurrent access. For bulk operations,
// consider using dedicated methods instead of multiple single operations.
//
// Error Handling:
// Context-aware methods return errors for:
//   - Context cancellation
//   - Invalid JSON format
//   - Type mismatches during unmarshaling
package set
