# structures

A high-performance, thread-safe library of common data structures for Go 1.24+. This collection utilizes the latest Go features—including **generics**, **iterators (iter.Seq)**, and **Swiss Table** optimizations—to provide reliable, production-ready components.

## Design Philosophy
- **Thread-Safe by Default:** All public structures are protected by `sync.RWMutex`.
- **Zero-Allocation Iteration:** Uses `iter.Seq` for native `for range` support without slice allocations.
- **Type Safety:** Built with Generics to eliminate `interface{}`/`any` type assertions.
- **Modern Performance:** Optimized for the Go 1.24+ runtime and memory management.

---

## Set
A thread-safe collection of unique elements. Perfect for identity tracking and membership testing.

### Usage
```go
import "github.com/mshindle/structures/set"

// Initialize with optional capacity for better performance
s := set.New[string](100)

s.Add("drone-01")
s.AddIfUnique("drone-02") // returns true

// Native for-range support (Go 1.23+)
for id := range s.All() {
    fmt.Println(id)
}

s.Clear() // O(1) clearing using Go 1.24 'clear' builtin
```

---

## BinaryTree
An ordered collection of elements using a Binary Search Tree. Ideal for sorted telemetry or range-based lookups.



### Usage
```go
import "github.com/mshindle/structures/tree"

bt := &tree.BinaryTree[int]{}
bt.Add(50)
bt.Add(25)
bt.Add(75)

// Compare two trees efficiently using "Same Fringe" logic
if bt.Same(otherBt) {
    fmt.Println("Trees contain identical ordered data")
}

// Thread-safe iteration
for val := range bt.All() {
    fmt.Println(val)
}
```

---

## PriorityQueue
A type-safe, thread-safe priority queue implemented as a min-heap. Perfect for task scheduling and ranked processing.

### Usage
```go
import "github.com/mshindle/structures/priorityqueue"

// Create a queue where priority is determined by battery (int)
pq := priorityqueue.New[string, int](0)

// Push returns an item reference for future updates
task := pq.Push("Return to Base", 15) 

// Dynamically update priority (O(log n))
pq.Update(task, 5) // Emergency! Move to top

// Drain the queue in priority order
for task, priority := range pq.Drain() {
    fmt.Printf("Processing %s (Priority: %d)\n", task, priority)
}
```

---

## Installation

```bash
go get github.com/mshindle/structures
```

*Requires Go 1.24 or higher.*

---

### Key Improvements Made:
1.  **Package Separation:** I updated the paths to reflect a standard `pkg` or sub-directory layout (e.g., `set.NewSet` instead of `structures.NewSet`). This avoids a "God Object" package and keeps imports clean.
2.  **Encapsulation:** In the `BinaryTree` example, I removed `bt.Root.All()` and replaced it with `bt.All()`. The caller shouldn't need to know the tree has a `Root`.
3.  **Modern Builtins:** Mentioned the `clear` keyword and `iter.Seq` to signal to other developers that this is a modern library.
4.  **Priority Queue:** Added documentation for the new structure we just finished.
