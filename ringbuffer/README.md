# RingBuffer & The LMAX Disruptor Philosophy

The `RingBuffer` implementation in `github.com/mshindle/structures` is designed for high-throughput telemetry ingestion where low latency and predictable performance are non-negotiable. It draws inspiration from the **LMAX Disruptor** architecture.

## 核心 (Core) Philosophy: High Throughput, Low Latency

The LMAX architecture achieves performance by avoiding the "stop-the-world" overhead typical of lock contention and frequent garbage collection. Our Go 1.26.1 implementation focuses on three pillars:

### 1. Mechanical Sympathy
Modern CPUs thrive on cache locality. By using a **fixed-size slice** allocated once at initialization, we ensure that data nodes stay contiguous in memory. This allows the CPU to pre-fetch data effectively, minimizing cache misses compared to the "pointer chasing" inherent in traditional Linked Lists.

### 2. The Zero-GC Path
In a high-churn environment (e.g., hundreds of drones streaming telemetry), allocating new objects for every "Push" creates massive pressure on the Go Garbage Collector.
* **LinkedList:** Allocates a new `node` struct for every entry.
* **RingBuffer:** Simply overwrites existing slots in a pre-allocated slice.

By reusing memory, we virtually eliminate GC cycles within the buffer logic itself.

### 3. Predictable Performance
Standard Go slices grow dynamically, which causes occasional "stutter" or latency spikes when the slice needs to reallocate and copy data. Because the `RingBuffer` is fixed-size, every operation is $O(1)$ with a constant time cost.

---

## Technical Implementation

### Optimized Methods for Telemetry

Beyond standard `Push` and `Pop`, this library includes specialized methods for real-time systems:

#### `TryPush(v T) bool`
In performance-critical loops, checking a boolean is significantly faster than handling an `error` interface. `TryPush` allows the caller to handle back-pressure (like dropping a packet or sleeping) without the allocation overhead of error strings.

#### `OverwritePush(v T)`
For telemetry, the **freshest data is the best data**. If the buffer is full, `OverwritePush` automatically drops the oldest entry at the head to make room for the new entry at the tail. This ensures the system remains **"Lossy but Live"**—preferring a gap in history over a complete system halt.

#### `PeekTail() (T, error)`
Enables O(1) access to the most recently ingested data without removing it from the processing queue. This is essential for "latest-state" dashboards.

---

## A Note on Lock-Free Programming

> [!IMPORTANT]
> **Senior Warning:** While the LMAX Disruptor is famous for its lock-free implementation, our current version uses a `sync.Mutex`. In Go 1.26, mutexes are highly optimized for short-duration locks.
>
> Moving to **Atomic Pointers** and **Memory Padding** (to avoid "False Sharing") can offer a performance edge in extreme scenarios, but it introduces significant complexity. For most use cases, including high-frequency drone simulation, the Mutex-based RingBuffer provides the best balance of maintainability and speed.

---

## v0.0.3 Roadmap: Linear Collections
The `RingBuffer` joins the `List` (Singly Linked List) in the **v0.0.3 "Linear Collections"** update.

* **Use `List`** when you need dynamic growth and stable pointers.
* **Use `RingBuffer`** when you need high-speed, fixed-capacity, or "lossy" telemetry ingestion.

---

### Integration Suggestion
The `RingBuffer` is the perfect backbone for an **OpenTelemetry** pipeline. You can easily wrap `OverwritePush` to fire an "event" whenever a packet is dropped, giving you instant visibility into where your ingestion bottlenecks are.
