package queue

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runCommonQueueTests(t *testing.T, createQueue func(capacity int) Queue[int]) {
	t.Helper()

	t.Run("new queue has zero length", func(t *testing.T) {
		q := createQueue(1)
		assert.Equal(t, 0, q.Length())
	})

	t.Run("cannot pop from empty queue", func(t *testing.T) {
		q := createQueue(1)
		_, err := q.Pop()
		assert.ErrorIs(t, err, ErrQueueEmpty)
	})

	t.Run("push increments length", func(t *testing.T) {
		q := createQueue(2)
		assert.NoError(t, q.Push(1))
		assert.Equal(t, 1, q.Length())

		assert.NoError(t, q.Push(2))
		assert.Equal(t, 2, q.Length())
	})

	t.Run("pop decrements length", func(t *testing.T) {
		q := createQueue(2)
		assert.NoError(t, q.Push(1))

		_, err := q.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 0, q.Length())
	})

	t.Run("peek does not change length", func(t *testing.T) {
		q := createQueue(2)
		assert.NoError(t, q.Push(1))

		_, err := q.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 1, q.Length())
	})

	t.Run("push and pop single item", func(t *testing.T) {
		q := createQueue(1)
		x := rand.Int()
		assert.NoError(t, q.Push(x))

		y, err := q.Pop()
		assert.NoError(t, err)
		assert.Equal(t, x, y)
	})

	t.Run("push and peek single item", func(t *testing.T) {
		q := createQueue(1)
		x := rand.Int()
		assert.NoError(t, q.Push(x))

		y, err := q.Peek()
		assert.NoError(t, err)
		assert.Equal(t, x, y)
	})

	t.Run("pop returns item in front of queue", func(t *testing.T) {
		q := createQueue(2)
		assert.NoError(t, q.Push(10))
		assert.NoError(t, q.Push(20))

		x, err := q.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 10, x)
	})

	t.Run("peek returns item in front of queue", func(t *testing.T) {
		q := createQueue(2)
		assert.NoError(t, q.Push(10))
		assert.NoError(t, q.Push(20))

		x, err := q.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 10, x)
	})

	t.Run("can push and pop more items than initial length", func(t *testing.T) {
		q := createQueue(2)

		assert.NoError(t, q.Push(10))
		assert.NoError(t, q.Push(20))
		_, err := q.Pop()
		assert.NoError(t, err)

		assert.NoError(t, q.Push(30))
		x, err := q.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 20, x)

		assert.NoError(t, q.Push(40))
		x, err = q.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 30, x)
	})
}

func runBoundedQueueTests(t *testing.T, createQueue func(capacity int) Queue[int]) {
	t.Helper()

	runCommonQueueTests(t, createQueue)

	t.Run("cannot push to full queue", func(t *testing.T) {
		q := createQueue(1)
		assert.NoError(t, q.Push(1))

		err := q.Push(2)
		assert.ErrorIs(t, err, ErrQueueFull)
	})
}

func runUnboundedQueueTests(t *testing.T, createQueue func(capacity int) Queue[int]) {
	t.Helper()

	runCommonQueueTests(t, createQueue)

	t.Run("resize after pushing with no pops", func(t *testing.T) {
		q := createQueue(2)
		const itemCount = 4

		for i := range itemCount {
			assert.NoError(t, q.Push(i+1))
		}

		for i := range itemCount {
			x, err := q.Pop()
			assert.NoError(t, err)
			assert.Equal(t, i+1, x)
		}
	})

	t.Run("resize after pushing with pop", func(t *testing.T) {
		t.Run("initial capacity of two", func(t *testing.T) {
			q := createQueue(2)
			assert.NoError(t, q.Push(1))
			assert.NoError(t, q.Push(2))
			_, _ = q.Pop()
			assert.NoError(t, q.Push(3))
			assert.NoError(t, q.Push(4))

			for i := range 3 {
				x, err := q.Pop()
				assert.NoError(t, err)
				assert.Equal(t, i+2, x)
			}
		})

		t.Run("initial capacity of five", func(t *testing.T) {
			q := createQueue(5)
			assert.NoError(t, q.Push(1))
			assert.NoError(t, q.Push(2))
			assert.NoError(t, q.Push(3))
			assert.NoError(t, q.Push(4))
			assert.NoError(t, q.Push(5))
			_, _ = q.Pop()
			_, _ = q.Pop()
			_, _ = q.Pop()
			assert.NoError(t, q.Push(6))
			assert.NoError(t, q.Push(7))
			assert.NoError(t, q.Push(8))
			assert.NoError(t, q.Push(9))

			for i := range 5 {
				x, err := q.Pop()
				assert.NoError(t, err)
				assert.Equal(t, i+4, x)
			}
		})
	})
}

func TestBoundedRingBufferQueue(t *testing.T) {
	runBoundedQueueTests(t, newBoundedRingBufferQueue[int])
}

func TestUnboundedRingBufferQueue(t *testing.T) {
	runUnboundedQueueTests(t, newUnboundedRingBufferQueue[int])
}
