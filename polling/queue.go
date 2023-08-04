package polling

import "sync"

type Queue[T any] struct {
	items []T
	lock  *sync.RWMutex
	cap   int
}

func NewQueue[T any](cap int) *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0, cap),
		lock:  new(sync.RWMutex),
		cap:   cap,
	}
}

func (q *Queue[T]) Enqueue(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if l := len(q.items); l == 0 {
		q.items = append(q.items, item)
	} else {
		to := q.cap - 1
		if l < q.cap {
			to = l
		}
		q.items = append([]T{item}, q.items[:to]...)
	}
}

func (q *Queue[T]) Copy() []T {
	q.lock.RLock()
	defer q.lock.RUnlock()

	copied := make([]T, len(q.items))
	copy(copied, q.items)

	return copied
}
