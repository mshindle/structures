package priorityqueue

import (
	"cmp"
	"container/heap"
	"iter"
	"sync"
)

// PriorityItem is now private-field focused to ensure index integrity
type PriorityItem[T any, P cmp.Ordered] struct {
	value    T
	priority P
	index    int
}

func (item *PriorityItem[T, P]) Value() T    { return item.value }
func (item *PriorityItem[T, P]) Priority() P { return item.priority }

// PriorityQueue wraps the internal heap to provide a Type-Safe API
type PriorityQueue[T any, P cmp.Ordered] struct {
	mu   sync.RWMutex
	impl *pqImpl[T, P]
}

func New[T any, P cmp.Ordered](capacity int) *PriorityQueue[T, P] {
	return &PriorityQueue[T, P]{
		impl: &pqImpl[T, P]{items: make([]*PriorityItem[T, P], 0, capacity)},
	}
}

// Type-safe wrappers
func (pq *PriorityQueue[T, P]) Push(value T, priority P) *PriorityItem[T, P] {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	item := &PriorityItem[T, P]{value: value, priority: priority}
	heap.Push(pq.impl, item)
	return item
}

func (pq *PriorityQueue[T, P]) Pop() (T, P) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if pq.impl.Len() == 0 {
		var zeroT T
		var zeroP P
		return zeroT, zeroP
	}

	item := heap.Pop(pq.impl).(*PriorityItem[T, P])
	return item.value, item.priority
}

func (pq *PriorityQueue[T, P]) Update(item *PriorityItem[T, P], priority P) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	item.priority = priority
	heap.Fix(pq.impl, item.index)
}

func (pq *PriorityQueue[T, P]) Len() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return pq.impl.Len()
}

// Drain returns an iterator that pops items from the queue until it's empty.
// This is perfect for processing a batch of prioritized tasks.
func (pq *PriorityQueue[T, P]) Drain() iter.Seq2[T, P] {
	return func(yield func(T, P) bool) {
		for {
			pq.mu.Lock()
			if pq.impl.Len() == 0 {
				pq.mu.Unlock()
				return
			}
			item := heap.Pop(pq.impl).(*PriorityItem[T, P])
			pq.mu.Unlock()

			if !yield(item.value, item.priority) {
				return
			}
		}
	}
}

// --- Internal implementation (unexported) ---

type pqImpl[T any, P cmp.Ordered] struct {
	items []*PriorityItem[T, P]
}

func (i *pqImpl[T, P]) Len() int           { return len(i.items) }
func (i *pqImpl[T, P]) Less(a, b int) bool { return i.items[a].priority < i.items[b].priority }
func (i *pqImpl[T, P]) Swap(a, b int) {
	i.items[a], i.items[b] = i.items[b], i.items[a]
	i.items[a].index = a
	i.items[b].index = b
}
func (i *pqImpl[T, P]) Push(x any) {
	item := x.(*PriorityItem[T, P])
	item.index = len(i.items)
	i.items = append(i.items, item)
}
func (i *pqImpl[T, P]) Pop() any {
	n := len(i.items)
	item := i.items[n-1]
	i.items[n-1] = nil // GC optimization: crucial in Go 1.26
	i.items = i.items[0 : n-1]
	return item
}
