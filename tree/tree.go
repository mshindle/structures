package tree

import (
	"cmp"
	"fmt"
	"iter"
	"sync"
)

// BinaryTree represents a binary search tree of ordered elements.
type BinaryTree[T cmp.Ordered] struct {
	mu   sync.RWMutex
	root *Tree[T]
}

// Add inserts a new value into the binary tree.
func (bt *BinaryTree[T]) Add(v T) {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	bt.root = bt.root.Insert(v)
}

// All returns a thread-safe iterator.
// It holds a Read Lock for the duration of the traversal.
func (bt *BinaryTree[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		bt.mu.RLock()
		defer bt.mu.RUnlock()
		bt.root.walk(yield)
	}
}

// Same compares two BinaryTree wrappers safely.
func (bt *BinaryTree[T]) Same(other *BinaryTree[T]) bool {
	if bt == other {
		return true
	}

	// Acquire read locks for both trees to ensure a consistent snapshot
	bt.mu.RLock()
	defer bt.mu.RUnlock()
	other.mu.RLock()
	defer other.mu.RUnlock()

	return bt.root.Same(other.root)
}

// String provides a safe string representation.
func (bt *BinaryTree[T]) String() string {
	bt.mu.RLock()
	defer bt.mu.RUnlock()
	return bt.root.String()
}

// Invert safely swaps every node in the tree.
// Requires a full Write Lock.
func (bt *BinaryTree[T]) Invert() {
	bt.mu.Lock()
	defer bt.mu.Unlock()
	bt.root.Invert()
}

// ToChannel returns a read-only channel.
// Note: This launches a goroutine that holds a Read Lock
// until the channel is fully consumed or closed.
func (bt *BinaryTree[T]) ToChannel(cap int) <-chan T {
	ch := make(chan T, cap)
	go func() {
		// We use the thread-safe All() which handles its own RLock
		for v := range bt.All() {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

// --- Internal Recursive Tree Node (Non-Thread-Safe) ---

// Tree represents a node in a binary search tree.
type Tree[T cmp.Ordered] struct {
	Left  *Tree[T]
	Value T
	Right *Tree[T]
}

// String returns a string representation of the tree in-order.
func (t *Tree[T]) String() string {
	if t == nil {
		return "()"
	}
	s := ""
	if t.Left != nil {
		s += t.Left.String() + " "
	}
	s += fmt.Sprint(t.Value)
	if t.Right != nil {
		s += " " + t.Right.String()
	}
	return "(" + s + ")"
}

// Insert adds a value to the tree and returns the updated subtree.
// Using a pointer receiver (*Tree[T]) allows us to modify the tree in place.
func (t *Tree[T]) Insert(v T) *Tree[T] {
	if t == nil {
		return &Tree[T]{Value: v}
	}

	if v < t.Value {
		t.Left = t.Left.Insert(v)
	} else if v > t.Value {
		// Note: We use else-if to handle duplicates.
		// If v == t.Value, we do nothing (Standard Set behavior).
		t.Right = t.Right.Insert(v)
	}

	return t
}

func (t *Tree[T]) walk(yield func(T) bool) bool {
	if t == nil {
		return true
	}
	return t.Left.walk(yield) && yield(t.Value) && t.Right.walk(yield)
}

// Invert swaps the left and right children of every node in the tree.
func (t *Tree[T]) Invert() {
	if t == nil {
		return
	}
	t.Left.Invert()
	t.Right.Invert()
	t.Left, t.Right = t.Right, t.Left
}

// All returns an iterator that performs an in-order traversal.
// We use 'any' instead of 'comparable' because we don't need to
// compare values just to walk the tree.
func (t *Tree[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		t.walk(yield)
	}
}

// ToChannel returns a channel that receives the values of the tree in-order.
// The channel is closed once all values have been sent.
func (t *Tree[T]) ToChannel(cap int) <-chan T {
	ch := make(chan T, cap)
	go func() {
		for v := range t.All() {
			ch <- v
		}
		close(ch)
	}()
	return ch
}

// Same reports whether t and other contain the same values in the same order.
func (t *Tree[T]) Same(other *Tree[T]) bool {
	// 1. Pointer equality optimization
	if t == other {
		return true
	}

	// 2. Initial nil checks
	// One is nil, the other isn't (or they'd be caught by t == other)
	if t == nil || other == nil {
		return false
	}

	// Create pulling iterators from our existing All() method
	next1, stop1 := iter.Pull(t.All())
	next2, stop2 := iter.Pull(other.All())

	// Crucial: stop the iterators to prevent goroutine leaks
	defer stop1()
	defer stop2()

	for {
		v1, ok1 := next1()
		v2, ok2 := next2()

		// If both are finished, the trees match
		if !ok1 && !ok2 {
			return true
		}

		// Compare current values
		if ok1 != ok2 || (ok1 && v1 != v2) {
			return false
		}
	}
}
