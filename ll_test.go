package structures

import "testing"

var listElements = []string{"alpha", "bravo", "charlie", "delta", "echo"}

func listFromStrings(lstr []string) *List[string] {
	ll := NewList[string]()
	for i := range lstr {
		ll.PushBack(lstr[i])
	}
	return ll
}

func TestList_Create(t *testing.T) {
	ll := listFromStrings(listElements)

	if ll.Len() != len(listElements) {
		t.Fatalf("got len %d, expected %d", ll.Len(), len(listElements))
	}

	i := 0
	for val := range ll.All() {
		if val != listElements[i] {
			t.Fatalf("got value '%v', expected '%s'", val, listElements[i])
		}
		i++
	}
}

func TestList_PushFront(t *testing.T) {
	ll := NewList[string]()

	for i := range listElements {
		ll.PushFront(listElements[i])
	}

	if ll.Len() != len(listElements) {
		t.Fatalf("got len %d, expected %d", ll.Len(), len(listElements))
	}

	i := len(listElements) - 1
	for val := range ll.All() {
		if val != listElements[i] {
			t.Fatalf("got value '%v', expected '%s'", val, listElements[i])
		}
		i--
	}
}

func TestList_PopFront(t *testing.T) {
	ll := listFromStrings(listElements)

	for _, expected := range listElements {
		val, ok := ll.PopFront()
		if !ok {
			t.Fatalf("expected value, got none")
		}
		if val != expected {
			t.Fatalf("got value '%v', expected '%s'", val, expected)
		}
	}

	val, ok := ll.PopFront()
	if ok {
		t.Fatalf("expected no value, got '%v'", val)
	}
}

func TestList_PushBack(t *testing.T) {
	ll := NewList[string]()

	for i := range listElements {
		ll.PushBack(listElements[i])
	}

	if ll.Len() != len(listElements) {
		t.Fatalf("got len %d, expected %d", ll.Len(), len(listElements))
	}

	i := 0
	for val := range ll.All() {
		if val != listElements[i] {
			t.Fatalf("got value '%v', expected '%s'", val, listElements[i])
		}
		i++
	}
}

func TestList_Len(t *testing.T) {
	ll := NewList[int]()
	if ll.Len() != 0 {
		t.Fatalf("expected empty list len 0, got %d", ll.Len())
	}

	ll.PushBack(1)
	if ll.Len() != 1 {
		t.Fatalf("expected len 1, got %d", ll.Len())
	}

	ll.PushFront(2)
	if ll.Len() != 2 {
		t.Fatalf("expected len 2, got %d", ll.Len())
	}

	ll.PopFront()
	if ll.Len() != 1 {
		t.Fatalf("expected len 1 after pop, got %d", ll.Len())
	}

	ll.PopFront()
	if ll.Len() != 0 {
		t.Fatalf("expected len 0 after second pop, got %d", ll.Len())
	}

	ll.PopFront()
	if ll.Len() != 0 {
		t.Fatalf("expected len 0 after pop on empty, got %d", ll.Len())
	}
}
