package loki

type FixedSizeQueue struct {
	size  int
	queue []string
	head  int
	tail  int
}

func NewFixedSizeQueue(size int) *FixedSizeQueue {
	return &FixedSizeQueue{
		size:  size,
		queue: make([]string, size),
		head:  0,
		tail:  0,
	}
}

func (q *FixedSizeQueue) Enqueue(value string) {
	if q.head == (q.tail+1)%q.size {
		q.Dequeue()
	}
	q.queue[q.tail] = value
	q.tail = (q.tail + 1) % q.size
}

func (q *FixedSizeQueue) Dequeue() string {
	value := q.queue[q.head]
	q.head = (q.head + 1) % q.size
	return value
}

func (q *FixedSizeQueue) Contains(value string) bool {
	for i := q.head; i != q.tail; i = (i + 1) % q.size {
		if q.queue[i] == value {
			return true
		}
	}
	return false
}
