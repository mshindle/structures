# structures

A library of common data structures meant to prevent copy/paste amongst multiple projects. This collection provides reliable, thread-safe implementations of frequently used structures.

## Set

The `Set` structure is a thread-safe collection of unique elements. It ensures that no duplicate values are stored and allows for concurrent access from multiple goroutines using an internal `sync.RWMutex`.

### Usage

```go
package main

import (
	"fmt"
	"github.com/mshindle/structures"
)

func main() {
	// Initialize a new set of integers
	s := structures.NewSet[int](0)

	// Add elements
	s.Add(1)
	s.Add(2)
	s.Add(1) // Duplicate, won't be added

	// Check for existence
	if s.Has(1) {
		fmt.Println("Set contains 1")
	}

	// Add if unique
	added := s.AddIfUnique(3)
	fmt.Printf("Added 3: %v\n", added)

	added = s.AddIfUnique(2)
	fmt.Printf("Added 2: %v\n", added) // Should be false

	// Iterate over all elements
	fmt.Println("All elements:")
	for v := range s.All() {
		fmt.Println(v)
	}
}
```
