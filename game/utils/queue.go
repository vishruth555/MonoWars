package utils

type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

func (q *Queue[T]) Dequeue() T {
	var item T
	if q.IsEmpty() {
		return item
	}
	item = q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue[T]) Size() int {
	return len(q.items)
}
