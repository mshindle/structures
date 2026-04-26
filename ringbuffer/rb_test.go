package ringbuffer

import (
	"errors"
	"sync"
	"testing"
)

func TestRingBuffer_New(t *testing.T) {
	rb := New[int](5)
	if rb.Capacity() != 5 {
		t.Errorf("expected capacity 5, got %d", rb.Capacity())
	}
	if rb.Len() != 0 {
		t.Errorf("expected length 0, got %d", rb.Len())
	}
}

func TestRingBuffer_PushPop(t *testing.T) {
	rb := New[int](3)

	// Test basic Push/Pop
	if err := rb.Push(1); err != nil {
		t.Fatalf("failed to push: %v", err)
	}
	if rb.Len() != 1 {
		t.Errorf("expected len 1, got %d", rb.Len())
	}

	val, err := rb.Pop()
	if err != nil {
		t.Fatalf("failed to pop: %v", err)
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
	if rb.Len() != 0 {
		t.Errorf("expected len 0, got %d", rb.Len())
	}

	// Test FIFO
	rb.Push(10)
	rb.Push(20)
	rb.Push(30)

	if rb.Len() != 3 {
		t.Errorf("expected len 3, got %d", rb.Len())
	}

	if err := rb.Push(40); !errors.Is(err, ErrFull) {
		t.Errorf("expected ErrFull, got %v", err)
	}

	v1, _ := rb.Pop()
	v2, _ := rb.Pop()
	v3, _ := rb.Pop()

	if v1 != 10 || v2 != 20 || v3 != 30 {
		t.Errorf("FIFO failed: got %d, %d, %d", v1, v2, v3)
	}

	if _, err := rb.Pop(); !errors.Is(err, ErrEmpty) {
		t.Errorf("expected ErrEmpty, got %v", err)
	}
}

func TestRingBuffer_TryPush(t *testing.T) {
	rb := New[int](2)

	if ok := rb.TryPush(1); !ok {
		t.Error("TryPush failed to push to empty buffer")
	}
	if ok := rb.TryPush(2); !ok {
		t.Error("TryPush failed to push to non-full buffer")
	}
	if ok := rb.TryPush(3); ok {
		t.Error("TryPush succeeded on full buffer")
	}
}

func TestRingBuffer_OverwritePush(t *testing.T) {
	rb := New[int](3)

	rb.OverwritePush(1)
	rb.OverwritePush(2)
	rb.OverwritePush(3)
	// Buffer: [1, 2, 3], head=0, tail=0 (circularly next), size=3

	rb.OverwritePush(4)
	// Should drop 1. Buffer effectively: [2, 3, 4], head=1, tail=1, size=3

	v1, _ := rb.Pop()
	if v1 != 2 {
		t.Errorf("expected 2, got %d", v1)
	}

	rb.OverwritePush(5)
	rb.OverwritePush(6)
	// After pop 2, len was 2.
	// Push 5: len 3. [5, 3, 4], head=2, tail=1
	// Push 6: len 3, drops 3. [5, 6, 4], head=0, tail=2

	v2, _ := rb.Pop() // should be 4
	v3, _ := rb.Pop() // should be 5
	v4, _ := rb.Pop() // should be 6

	if v2 != 4 || v3 != 5 || v4 != 6 {
		t.Errorf("OverwritePush sequence failed: got %d, %d, %d", v2, v3, v4)
	}
}

func TestRingBuffer_PeekTail(t *testing.T) {
	rb := New[int](3)

	if _, err := rb.PeekTail(); !errors.Is(err, ErrEmpty) {
		t.Errorf("expected ErrEmpty on empty PeekTail, got %v", err)
	}

	rb.Push(100)
	if val, _ := rb.PeekTail(); val != 100 {
		t.Errorf("expected 100, got %d", val)
	}

	rb.Push(200)
	if val, _ := rb.PeekTail(); val != 200 {
		t.Errorf("expected 200, got %d", val)
	}

	rb.Pop()
	if val, _ := rb.PeekTail(); val != 200 {
		t.Errorf("expected 200 after pop, got %d", val)
	}
}

func TestRingBuffer_Concurrency(t *testing.T) {
	rb := New[int](100)
	var wg sync.WaitGroup
	numOps := 1000

	// Producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numOps; i++ {
			for !rb.TryPush(i) {
				// busy wait/yield
			}
		}
	}()

	// Consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 0
		for count < numOps {
			if _, err := rb.Pop(); err == nil {
				count++
			}
		}
	}()

	wg.Wait()
	if rb.Len() != 0 {
		t.Errorf("expected empty buffer after concurrent ops, got %d", rb.Len())
	}
}
