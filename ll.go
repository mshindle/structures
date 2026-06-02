package structures

import (
	"iter"
	"sync"
)

// List represents a thread-safe singly linked list.
// Optimized for O(1) PushBack, PushFront, and PopFront.
type List[T any] struct {
	mu   sync.RWMutex
	head *node[T]
	tail *node[T]
	len  int
}

type node[T any] struct {
	value T
	next  *node[T]
}

func NewList[T any]() *List[T] {
	return &List[T]{}
}

// PushBack adds an element to the end of the list (O(1)).
// Use this for FIFO Queueing.
func (l *List[T]) PushBack(v T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	n := &node[T]{value: v}
	if l.tail == nil {
		l.head = n
		l.tail = n
	} else {
		l.tail.next = n
		l.tail = n
	}
	l.len++
}

// PushFront adds an element to the beginning of the list (O(1)).
// Use this for LIFO Stacks.
func (l *List[T]) PushFront(v T) {
	l.mu.Lock()
	defer l.mu.Unlock()

	n := &node[T]{value: v, next: l.head}
	l.head = n
	if l.tail == nil {
		l.tail = n
	}
	l.len++
}

// PopFront removes and returns the first element (O(1)).
func (l *List[T]) PopFront() (T, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.head == nil {
		var zero T
		return zero, false
	}

	val := l.head.value
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	}
	l.len--
	return val, true
}

// PeekFront returns the first element without removing it.
func (l *List[T]) PeekFront() (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.head == nil {
		var zero T
		return zero, false
	}
	return l.head.value, true
}

// All returns an iterator for the list.
func (l *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		l.mu.RLock()
		defer l.mu.RUnlock()
		for curr := l.head; curr != nil; curr = curr.next {
			if !yield(curr.value) {
				return
			}
		}
	}
}

func (l *List[T]) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.len
}
