// Written by ChatGPT, not reviewed and may have potential issues
package collections

type OverflowQueue[T any] struct {
	data     []T
	capacity int
	front    int
	rear     int
	size     int
}

func NewOverflowQueue[T any](capacity int) *OverflowQueue[T] {
	return &OverflowQueue[T]{
		data:     make([]T, capacity),
		capacity: capacity,
		front:    0,
		rear:     0,
		size:     0,
	}
}

func (q *OverflowQueue[T]) Enqueue(item T) {
	if q.size == q.capacity {
		// 如果队列已满，覆盖最早的元素
		q.front = (q.front + 1) % q.capacity
	} else {
		q.size++
	}
	q.data[q.rear] = item
	q.rear = (q.rear + 1) % q.capacity
}

func (q *OverflowQueue[T]) Dequeue() (T, bool) {
	var zeroValue T
	if q.size == 0 {
		return zeroValue, false // 队列为空
	}
	item := q.data[q.front]
	q.front = (q.front + 1) % q.capacity
	q.size--
	return item, true
}

func (q *OverflowQueue[T]) IsEmpty() bool {
	return q.size == 0
}

func (q *OverflowQueue[T]) IsFull() bool {
	return q.size == q.capacity
}

func (q *OverflowQueue[T]) Size() int {
	return q.size
}

// Iterator 迭代器结构体
type Iterator[T any] struct {
	queue     *OverflowQueue[T]
	current   int
	traversed int
}

// NewIterator 构造迭代器
func (q *OverflowQueue[T]) NewIterator() *Iterator[T] {
	return &Iterator[T]{queue: q, current: q.front, traversed: 0}
}

// Next 获取下一个元素
func (it *Iterator[T]) Next() (T, bool) {
	var zeroValue T
	if it.traversed >= it.queue.size {
		return zeroValue, false // 已遍历完所有元素
	}
	item := it.queue.data[it.current]
	it.current = (it.current + 1) % it.queue.capacity
	it.traversed++
	return item, true
}

// ReverseIterator 反向迭代器结构体
type ReverseIterator[T any] struct {
	queue     *OverflowQueue[T]
	current   int
	traversed int
}

// NewReverseIterator 构造反向迭代器
func (q *OverflowQueue[T]) NewReverseIterator() *ReverseIterator[T] {
	return &ReverseIterator[T]{queue: q, current: (q.rear - 1 + q.capacity) % q.capacity, traversed: 0}
}

// Next 获取下一个反向元素
func (it *ReverseIterator[T]) Next() (T, bool) {
	var zeroValue T
	if it.traversed >= it.queue.size {
		return zeroValue, false // 已遍历完所有元素
	}
	item := it.queue.data[it.current]
	it.current = (it.current - 1 + it.queue.capacity) % it.queue.capacity
	it.traversed++
	return item, true
}
