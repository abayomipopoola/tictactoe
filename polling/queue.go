package polling

import "sync"

type Queue[T any] struct {
	items []T
	lock  sync.RWMutex
	cap   int
}

func NewQueue[T any](cap int) *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0, cap),
		cap:   cap,
	}
}

func (q *Queue[T]) Enqueue(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// If queue is at capacity, remove the oldest item.
	if len(q.items) == q.cap {
		q.items = q.items[1:]
	}

	q.items = append(q.items, item)
}

func (q *Queue[T]) Copy() []T {
	q.lock.RLock()
	defer q.lock.RUnlock()

	// Return a copy of the slice
	return append([]T(nil), q.items...)
}
