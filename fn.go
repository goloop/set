package set

// New is a constructor function that creates a new Set[T] instance.
// It accepts an arbitrary number of items of a generic type 'T' which
// can be either simple types (e.g., int, string, bool) or complex types
// (e.g., struct, slice).
//
// This function first creates a new, empty set. It then determines whether
// the Set is simple or complex based on the type of the first item, and
// caches this information for efficient subsequent operations. Finally,
// it adds the provided items to the Set.
//
// Note: All items must be of the same type. If different types are provided,
// the behavior is undefined.
//
// Example usage:
//
//	// Creating a set of simple type (int)
//	emptySet := set.New[int]()       // empty set of int
//	simpleSet := set.New(1, 2, 3, 4) // set of int
//
//	// Creating a set of complex type (struct).
//	type ComplexType struct {
//	    field1 int
//	    field2 string
//	}
//	complexSet := set.New(
//	    ComplexType{1, "one"},
//	    ComplexType{2, "two"},
//	)
//
//	// Adding an item to the set.
//	simpleSet.Add(5)
//	complexSet.Add(ComplexType{3, "three"})
//
//	// Checking if an item exists in the set.
//	existsSimple := simpleSet.Contains(3)                       // returns true
//	existsComplex := complexSet.Contains(ComplexType{2, "two"}) // returns true
//
//	// Getting the size of the set.
//	size := simpleSet.Len() // returns 5
func New[T any](items ...T) *Set[T] {
	set := Set[T]{
		heap:   make(map[string]T),
		simple: 0,
	}
	set.IsSimple()    // cache the complexity of the object
	set.Add(items...) // add items to the set

	return &set
}

// Map returns a new set with the results of applying the provided function
// to each item in the set.
//
// Example usage:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	s := set.New[User]()
//	s.Add(User{"John", 20}, User{"Jane", 30})
//
//	names := sort.Map(s, func(item User) string {
//	    return item.Name
//	})
//
//	fmt.Println(names.Elements()) // "Jane", "John"
func Map[T any, R any](s *Set[T], fn func(item T) R) *Set[R] {
	resulet := New[R]()
	for _, v := range s.heap {
		resulet.Add(fn(v))
	}

	return resulet
}

// Reduce returns a single value by applying the provided function to each
// item in the set and passing the result of previous function call as the
// first argument in the next call.
//
// Example usage:
//
//	type User struct {
//			Name string
//			Age  int
//	}
//
//	 s := set.New[User]()
//	 s.Add(User{"John", 20}, User{"Jane", 30})
//
//	 sum := sort.Reduce(s, func(acc int, item User) int {
//			return acc + item.Age
//	 }) // sum is 50
func Reduce[T any, R any](s *Set[T], fn func(acc R, item T) R) R {
	var acc R
	for _, v := range s.heap {
		acc = fn(acc, v)
	}

	return acc
}
