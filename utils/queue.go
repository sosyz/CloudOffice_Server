package utils

type Queue[T any] struct {
	Header *QueueNode[T]
	Tail   *QueueNode[T]
	size   int
}

type QueueNode[T any] struct {
	item T
	next *QueueNode[T]
}

func (q *Queue[T]) Push(item T) {
	node := &QueueNode[T]{item, nil}
	if q.size == 0 {
		q.Header = node
		q.Tail = node
	} else {
		q.Tail.next = node
		q.Tail = node
	}
	q.size++
}

func (q *Queue[T]) Pop() T {
	if q.size == 0 {
		var zero T
		return zero
	}
	item := q.Header.item
	q.Header = q.Header.next
	q.size--
	return item
}

func (q *Queue[T]) Size() int {
	return q.size
}

func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}

func (q *Queue[T]) Clear() {
	q.Header = nil
	q.Tail = nil
	q.size = 0
}
