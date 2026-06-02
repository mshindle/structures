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

## List
A thread-safe singly linked list. Supports O(1) insertions at both ends and removal from the front, making it suitable for stacks and queues.

### Thought process for building List in this format

* **Stable Pointers:** Unlike a slice, `PushBack` or `PushFront` will never move existing nodes in memory. This is critical for your drone system if you're holding onto references of pending tasks.
* **Predictable Latency:** No "growth spikes." Every insertion is just one small allocation and a pointer swap.
* **Clean API:** We've avoided the `interface{}` trap of the standard library's `container/list`.

### Potential features to add
* A `Wait/Notify` mechanism (like a `Cond` variable) so goroutines can block until an item is available in the list

### Usage
```go
import "github.com/mshindle/structures"

l := structures.NewList[int]()
l.PushBack(10)
l.PushFront(5)

// Fast removal from front (O(1))
if val, ok := l.PopFront(); ok {
    fmt.Println("Popped:", val) // 5
}

// Native for-range iteration
for val := range l.All() {
    fmt.Println(val) // 10
}
```

---

## RingBuffer
A fixed-size circular buffer designed for high-throughput telemetry.
Supports lossy ingestion via `OverwritePush` to ensure the system remains live
under heavy load.

### RingBuffer & The LMAX Disruptor Philosophy

The `RingBuffer` implementation in `github.com/mshindle/structures` is designed for high-throughput telemetry ingestion where low latency and predictable performance are non-negotiable. It draws inspiration from the **LMAX Disruptor** architecture.

### (Core) Philosophy: High Throughput, Low Latency

The LMAX architecture achieves performance by avoiding the "stop-the-world" overhead typical of lock contention and frequent garbage collection. Our Go 1.26.1 implementation focuses on three pillars:

#### 1. Mechanical Sympathy
Modern CPUs thrive on cache locality. By using a **fixed-size slice** allocated once at initialization, we ensure that data nodes stay contiguous in memory. This allows the CPU to pre-fetch data effectively, minimizing cache misses compared to the "pointer chasing" inherent in traditional Linked Lists.

#### 2. The Zero-GC Path
In a high-churn environment (e.g., hundreds of drones streaming telemetry), allocating new objects for every "Push" creates massive pressure on the Go Garbage Collector.
* **LinkedList:** Allocates a new `node` struct for every entry.
* **RingBuffer:** Simply overwrites existing slots in a pre-allocated slice.

By reusing memory, we virtually eliminate GC cycles within the buffer logic itself.

#### 3. Predictable Performance
Standard Go slices grow dynamically, which causes occasional "stutter" or latency spikes when the slice needs to reallocate and copy data. Because the `RingBuffer` is fixed-size, every operation is $O(1)$ with a constant time cost.

### Technical Implementation

#### Optimized Methods for Telemetry

Beyond standard `Push` and `Pop`, this library includes specialized methods for real-time systems:

#### `TryPush(v T) bool`
In performance-critical loops, checking a boolean is significantly faster than handling an `error` interface. `TryPush` allows the caller to handle back-pressure (like dropping a packet or sleeping) without the allocation overhead of error strings.

#### `OverwritePush(v T)`
For telemetry, the **freshest data is the best data**. If the buffer is full, `OverwritePush` automatically drops the oldest entry at the head to make room for the new entry at the tail. This ensures the system remains **"Lossy but Live"**—preferring a gap in history over a complete system halt.

#### `PeekTail() (T, error)`
Enables O(1) access to the most recently ingested data without removing it from the processing queue. This is essential for "latest-state" dashboards.

### A Note on Lock-Free Programming

> [!IMPORTANT]
> While the LMAX Disruptor is famous for its lock-free implementation, our current version uses a `sync.Mutex`. In Go 1.26, mutexes are highly optimized for short-duration locks.
>
> Moving to **Atomic Pointers** and **Memory Padding** (to avoid "False Sharing") can offer a performance edge in extreme scenarios, but it introduces significant complexity. For most use cases, including high-frequency drone simulation, the Mutex-based RingBuffer provides the best balance of maintainability and speed.


### v0.0.3 Roadmap: Linear Collections
The `RingBuffer` joins the `List` (Singly Linked List) in the **v0.0.3 "Linear Collections"** update.

* **Use `List`** when you need dynamic growth and stable pointers.
* **Use `RingBuffer`** when you need high-speed, fixed-capacity, or "lossy" telemetry ingestion.

### Integration Suggestion
The `RingBuffer` is the perfect backbone for an **OpenTelemetry** pipeline. You can easily wrap `OverwritePush` to fire an "event" whenever a packet is dropped, giving you instant visibility into where your ingestion bottlenecks are.

### Usage
```go
import "[github.com/mshindle/structures](https://github.com/mshindle/structures)"

rb := ringbuffer.New[float64](1024)

// Standard Push/Pop
rb.Push(42.5)
val, _ := rb.Pop()

// High-frequency telemetry (latest data prioritized)
rb.OverwritePush(101.2)

// Peek at the most recent entry
latest, _ := rb.PeekTail()
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
5.  **List:** Integrated the new generic linked list with native iteration.
