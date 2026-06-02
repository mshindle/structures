package structures

import (
	"errors"
	"sync"
)

var ErrFull = errors.New("ring buffer is full")
var ErrEmpty = errors.New("ring buffer is empty")

// RingBuffer is a thread-safe, fixed-size circular buffer.
type RingBuffer[T any] struct {
	mu    sync.Mutex
	data  []T
	head  int // Index of the first element (for Read)
	tail  int // Index of the next available slot (for Write)
	size  int // Current number of elements
	limit int // Maximum capacity
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		data:  make([]T, capacity),
		limit: capacity,
	}
}

// Push adds an element to the tail. Returns ErrFull if no space.
func (rb *RingBuffer[T]) Push(v T) error {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.size == rb.limit {
		return ErrFull
	}

	rb.data[rb.tail] = v
	rb.tail = (rb.tail + 1) % rb.limit
	rb.size++
	return nil
}

// Pop removes an element from the head. Returns ErrEmpty if empty.
func (rb *RingBuffer[T]) Pop() (T, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.size == 0 {
		var zero T
		return zero, ErrEmpty
	}

	val := rb.data[rb.head]

	// Optional: Clear reference for GC if T is a pointer
	var zero T
	rb.data[rb.head] = zero

	rb.head = (rb.head + 1) % rb.limit
	rb.size--
	return val, nil
}

func (rb *RingBuffer[T]) Len() int {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.size
}

func (rb *RingBuffer[T]) Capacity() int {
	return rb.limit
}

// TryPush attempts to add an element.
// It returns false if the buffer is full, rather than an error,
// allowing for high-frequency polling without error-string overhead.
func (rb *RingBuffer[T]) TryPush(v T) bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.size == rb.limit {
		return false
	}

	rb.data[rb.tail] = v
	rb.tail = (rb.tail + 1) % rb.limit
	rb.size++
	return true
}

// OverwritePush adds an element to the tail.
// If the buffer is full, it discards the oldest element (at the head)
// to make room. This is the "LMAX/Disruptor" style for telemetry.
func (rb *RingBuffer[T]) OverwritePush(v T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.size == rb.limit {
		// Drop the oldest data by advancing the head
		rb.head = (rb.head + 1) % rb.limit
		rb.size--
	}

	rb.data[rb.tail] = v
	rb.tail = (rb.tail + 1) % rb.limit
	rb.size++
}

// PeekTail returns the most recently added element without removing it.
// This is the element at (tail - 1) in the circular buffer.
func (rb *RingBuffer[T]) PeekTail() (T, error) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	if rb.size == 0 {
		var zero T
		return zero, ErrEmpty
	}

	// Calculate the index of the last element added
	index := (rb.tail - 1 + rb.limit) % rb.limit
	return rb.data[index], nil
}
