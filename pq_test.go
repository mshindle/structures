package structures

import (
	"slices"
	"testing"
)

func TestPriorityQueue_Basic(t *testing.T) {
	pq := NewPriorityQueue[string, int](0)

	pq.Push("orange", 3)
	pq.Push("apple", 1)
	pq.Push("banana", 2)

	if pq.Len() != 3 {
		t.Errorf("expected Len 3, got %d", pq.Len())
	}

	expected := []string{"apple", "banana", "orange"}
	i := 0
	for val, _ := range pq.Drain() {
		if val != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], val)
		}
		i++
	}

	if pq.Len() != 0 {
		t.Errorf("expected Len 0, got %d", pq.Len())
	}
}

func TestPriorityQueue_Update(t *testing.T) {
	pq := NewPriorityQueue[string, int](0)

	item := pq.Push("initial", 10)
	pq.Push("other", 5)

	// Use the Update receiver
	pq.Update(item, 1)

	if item.Priority() != 1 {
		t.Errorf("expected priority 1, got %d", item.Priority())
	}

	val, priority := pq.Pop()
	if val != "initial" {
		t.Errorf("expected initial, got %s", val)
	}
	if priority != 1 {
		t.Errorf("expected priority 1, got %d", priority)
	}
}

func TestPriorityQueue_Pop(t *testing.T) {
	pq := NewPriorityQueue[string, int](0)
	pq.Push("a", 10)
	pq.Push("b", 5)
	pq.Push("c", 15)

	val, _ := pq.Pop()
	if val != "b" {
		t.Errorf("expected b, got %s", val)
	}
}

func TestPriorityQueue_EmptyPop(t *testing.T) {
	pq := NewPriorityQueue[string, int](0)
	val, priority := pq.Pop()
	if val != "" || priority != 0 {
		t.Errorf("expected zero values for empty pop, got %v, %v", val, priority)
	}
}

func TestPriorityQueue_Drain(t *testing.T) {
	pq := NewPriorityQueue[int, int](0)
	items := []int{5, 3, 8, 1}
	for _, v := range items {
		pq.Push(v, v)
	}

	var got []int
	for v, _ := range pq.Drain() {
		got = append(got, v)
	}

	slices.Sort(items)
	if !slices.Equal(got, items) {
		t.Errorf("Drain() = %v, want %v", got, items)
	}
}

func TestPriorityQueue_SamePriority(t *testing.T) {
	pq := NewPriorityQueue[string, int](0)

	pq.Push("a", 1)
	pq.Push("b", 1)

	if pq.Len() != 2 {
		t.Errorf("expected Len 2, got %d", pq.Len())
	}

	pq.Pop()
	pq.Pop()

	if pq.Len() != 0 {
		t.Errorf("expected Len 0, got %d", pq.Len())
	}
}
