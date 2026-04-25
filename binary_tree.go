package structures

import (
	"cmp"
	"fmt"
	"iter"
)

// BinaryTree represents a binary search tree of ordered elements.
type BinaryTree[T cmp.Ordered] struct {
	Root *Tree[T]
}

// Add inserts a new value into the binary tree.
func (bt *BinaryTree[T]) Add(v T) {
	bt.Root = bt.Root.Insert(v)
}

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

// walk is the private recursive helper.
// It returns false if the caller wants to stop early (e.g., a break).
func (t *Tree[T]) walk(yield func(T) bool) bool {
	if t == nil {
		return true
	}

	// 1. Walk Left
	if !t.Left.walk(yield) {
		return false
	}

	// 2. Yield current Value
	if !yield(t.Value) {
		return false
	}

	// 3. Walk Right
	return t.Right.walk(yield)
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
