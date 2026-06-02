package structures

import (
	"slices"
	"sync"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet[string](0)

	s.Add("apple")
	s.Add("banana")
	s.Add("apple") // Duplicate

	if !s.Has("apple") {
		t.Error("Set should have 'apple'")
	}
	if !s.Has("banana") {
		t.Error("Set should have 'banana'")
	}
	if s.Has("cherry") {
		t.Error("Set should not have 'cherry'")
	}

	if s.AddIfUnique("cherry") != true {
		t.Error("AddIfUnique should return true for new element")
	}
	if s.AddIfUnique("cherry") != false {
		t.Error("AddIfUnique should return false for existing element")
	}

	elements := slices.Collect(s.All())
	if len(elements) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(elements))
	}

	expected := []string{"apple", "banana", "cherry"}
	for _, e := range expected {
		if !slices.Contains(elements, e) {
			t.Errorf("Set missing expected element: %s", e)
		}
	}
}

func TestSet_Concurrent(t *testing.T) {
	s := NewSet[int](0)
	var wg sync.WaitGroup
	n := 1000

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(v int) {
			defer wg.Done()
			s.Add(v)
		}(i)
	}

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(v int) {
			defer wg.Done()
			s.Has(v)
		}(i)
	}

	wg.Wait()

	count := 0
	for range s.All() {
		count++
	}
	if count != n {
		t.Errorf("Expected %d elements, got %d", n, count)
	}
}
