// Package queue provides a ring buffer-based queue implementation.
package queue

import "errors"

var (
	// ErrQueueEmpty is an error returned when an attempt is made to take an element from an empty queue.
	ErrQueueEmpty = errors.New("cannot take element from empty queue")

	// ErrQueueFull is an error returned when an attempt is made to add an element to a full queue.
	ErrQueueFull = errors.New("queue is full and cannot accept more elements")
)

type Queue[Element any] interface {
	// Push adds an element to the end of the queue. If the queue cannot accept more elements, the ErrQueueFull error is returned.
	Push(Element) error

	// Pop removes and returns the first element of the queue. If the queue is empty, the ErrQueueEmpty error is returned.
	Pop() (Element, error)

	// Peek returns the first element of the queue. If the queue is empty, the ErrQueueEmpty error is returned.
	Peek() (Element, error)

	// Length returns the number of elements in the queue.
	Length() int
}

type ringBufferQueue[Element any] struct {
	items   []Element
	front   int
	length  int
	bounded bool
}

// NewBoundedQueue returns a new queue with a maximum specific capacity.
func NewBoundedQueue[Element any](capacity int) Queue[Element] {
	return newBoundedRingBufferQueue[Element](capacity)
}

// NewUnboundedQueue returns a new queue with the specific initial capacity.
// The queue will resize its internal storage if its current capacity is exceeded.
// This implementation will double the internal capacity during each resize operation.
func NewUnboundedQueue[Element any](initialCapacity int) Queue[Element] {
	return newUnboundedRingBufferQueue[Element](initialCapacity)
}

const defaultRingBufferQueueCapacity = 2

func newBoundedRingBufferQueue[Element any](capacity int) Queue[Element] {
	if capacity == 0 {
		capacity = defaultRingBufferQueueCapacity
	}

	return &ringBufferQueue[Element]{
		items:   make([]Element, capacity),
		bounded: true,
	}
}

func newUnboundedRingBufferQueue[Element any](initialCapacity int) Queue[Element] {
	if initialCapacity == 0 {
		initialCapacity = defaultRingBufferQueueCapacity
	}

	return &ringBufferQueue[Element]{
		items:   make([]Element, initialCapacity),
		bounded: false,
	}
}

func (q *ringBufferQueue[Element]) expand() {
	if q.length < cap(q.items) {
		return
	}

	newCapacity := cap(q.items) * 2
	newItems := make([]Element, newCapacity)

	copyCount := copy(newItems, q.items[q.front:q.length])
	if q.front > 0 {
		copy(newItems[copyCount:], q.items[0:q.front])
	}

	q.items = newItems
	q.front = 0
}

func (q *ringBufferQueue[Element]) Push(item Element) error {
	if q.length == cap(q.items) {
		if q.bounded {
			return ErrQueueFull
		}

		q.expand()
	}

	back := (q.front + q.length) % cap(q.items)
	q.items[back] = item
	q.length++

	return nil
}

func (q *ringBufferQueue[Element]) Pop() (Element, error) {
	item, err := q.Peek()
	if err != nil {
		return item, err
	}

	q.front = (q.front + 1) % cap(q.items)
	q.length--

	return item, nil
}

func (q *ringBufferQueue[Element]) Peek() (Element, error) {
	var item Element

	if q.length == 0 {
		return item, ErrQueueEmpty
	}

	return q.items[q.front], nil
}

func (q *ringBufferQueue[Element]) Length() int {
	return q.length
}
