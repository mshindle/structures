package structures

import (
	"slices"
	"testing"
)

func TestBinaryTree(t *testing.T) {
	bt := &BinaryTree[int]{}
	bt.Add(5)
	bt.Add(3)
	bt.Add(7)
	bt.Add(3) // Duplicate, should be ignored

	expected := "((3) 5 (7))"
	if got := bt.Root.String(); got != expected {
		t.Errorf("String() = %q, want %q", got, expected)
	}

	values := slices.Collect(bt.Root.All())
	expectedValues := []int{3, 5, 7}
	if !slices.Equal(values, expectedValues) {
		t.Errorf("All() = %v, want %v", values, expectedValues)
	}
}

func TestTree_Invert(t *testing.T) {
	bt := &BinaryTree[int]{}
	bt.Add(5)
	bt.Add(3)
	bt.Add(7)

	bt.Root.Invert()

	expected := "((7) 5 (3))"
	if got := bt.Root.String(); got != expected {
		t.Errorf("After Invert, String() = %q, want %q", got, expected)
	}

	values := slices.Collect(bt.Root.All())
	expectedValues := []int{7, 5, 3}
	if !slices.Equal(values, expectedValues) {
		t.Errorf("After Invert, All() = %v, want %v", values, expectedValues)
	}
}

func TestTree_Same(t *testing.T) {
	t1 := &BinaryTree[int]{}
	for _, v := range []int{5, 3, 7, 2, 4} {
		t1.Add(v)
	}

	t2 := &BinaryTree[int]{}
	for _, v := range []int{5, 3, 7, 2, 4} {
		t2.Add(v)
	}

	if !t1.Root.Same(t2.Root) {
		t.Error("Same() returned false for identical trees")
	}

	t3 := &BinaryTree[int]{}
	for _, v := range []int{5, 3, 7, 2, 4} {
		t3.Add(v)
	}
	if !t1.Root.Same(t3.Root) {
		t.Error("Same() returned false for trees with same elements added in same order")
	}

	t4 := &BinaryTree[int]{}
	for _, v := range []int{5, 7, 3, 4, 2} { // Different structure, same values
		t4.Add(v)
	}
	if !t1.Root.Same(t4.Root) {
		t.Error("Same() returned false for trees with same elements but different structure")
	}

	t5 := &BinaryTree[int]{}
	for _, v := range []int{5, 3, 7, 2} { // Missing 4
		t5.Add(v)
	}
	if t1.Root.Same(t5.Root) {
		t.Error("Same() returned true for trees with different elements")
	}
}

func TestTree_ToChannel(t *testing.T) {
	bt := &BinaryTree[int]{}
	nums := []int{5, 3, 7, 2, 4, 6, 8}
	for _, n := range nums {
		bt.Add(n)
	}

	ch := bt.Root.ToChannel(0)
	var got []int
	for v := range ch {
		got = append(got, v)
	}

	expected := []int{2, 3, 4, 5, 6, 7, 8}
	if !slices.Equal(got, expected) {
		t.Errorf("ToChannel values = %v, want %v", got, expected)
	}
}
