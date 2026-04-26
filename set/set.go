package set

import (
	"iter"
	"sync"
)

// Set is a thread-safe collection of unique elements.
type Set[T comparable] struct {
	mu sync.RWMutex
	m  map[T]struct{}
}

// NewSet initializes a new Set with an optional capacity.
func New[T comparable](size int) *Set[T] {
	return &Set[T]{
		m: make(map[T]struct{}, size),
	}
}

// Add adds an element to the set. It is thread-safe.
func (s *Set[T]) Add(v T) {
	s.mu.Lock() // Exclusive lock for writing
	defer s.mu.Unlock()
	s.m[v] = struct{}{}
}

// AddIfUnique adds an element to the set only if it is not already present.
// It returns true if the element was added, and false if it was already in the set.
// This method is thread-safe.
func (s *Set[T]) AddIfUnique(v T) bool {
	// 1. Quick check with a Read Lock (Fast Path)
	s.mu.RLock()
	if _, ok := s.m[v]; ok {
		s.mu.RUnlock()
		return false // Already exists
	}
	s.mu.RUnlock()

	// 2. ID wasn't there, so grab the Write Lock (Slow Path)
	s.mu.Lock()
	defer s.mu.Unlock()

	// 3. Double-check! Another goroutine might have
	// added it between the RUnlock and the Lock.
	if _, ok := s.m[v]; ok {
		return false
	}

	s.m[v] = struct{}{}
	return true
}

// Has checks if an element is present in the set.
// It returns true if the element is found, and false otherwise.
// This method is thread-safe.
func (s *Set[T]) Has(v T) bool {
	s.mu.RLock() // Shared lock for reading
	defer s.mu.RUnlock()
	_, ok := s.m[v]
	return ok
}

// Remove deletes an element from the set.
func (s *Set[T]) Remove(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, v)
}

// Len returns the current number of elements.
func (s *Set[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}

// Clear removes all elements from the set.
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Go 1.24+ optimization: clearing a map is now faster than reallocating
	clear(s.m)
}

// All returns an iterator for the set.
// It uses a read lock to ensure thread-safety during iteration.
func (s *Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		// We MUST hold the RLock while iterating to prevent
		// concurrent writes from corrupting the iteration.
		s.mu.RLock()
		defer s.mu.RUnlock()

		for v := range s.m {
			if !yield(v) {
				return // The caller stopped the loop early (e.g., a 'break')
			}
		}
	}
}
