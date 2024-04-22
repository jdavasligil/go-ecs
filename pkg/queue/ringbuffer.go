package queue

import "sync"

type RingBuffer[T any] struct {
	mu       sync.Mutex
	back     int
	buf      []T
	length   int
	capacity int
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buf:      make([]T, capacity),
		capacity: capacity,
	}
}

func (rb *RingBuffer[T]) Push(t T) {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	rb.buf[rb.back] = t
	rb.back = (rb.back + 1) % rb.capacity
	if rb.length < rb.capacity {
		rb.length += 1
	}
}

func (rb *RingBuffer[T]) Pop() T {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	if rb.length == 0 {
		var noop T
		return noop
	}

	rb.length -= 1

	return rb.buf[(rb.capacity+(rb.back-rb.length-1))%rb.capacity]
}

func (rb *RingBuffer[T]) Peek() T {
	rb.mu.Lock()
	defer rb.mu.Unlock()
	if rb.length == 0 {
		var noop T
		return noop
	}

	return rb.buf[(rb.capacity+(rb.back-rb.length))%rb.capacity]
}

func (rb *RingBuffer[T]) Len() int {
	return rb.length
}

func (rb *RingBuffer[T]) IsEmpty() bool {
	return rb.length == 0
}

func (rb *RingBuffer[T]) IsFull() bool {
	return rb.length == rb.capacity
}
