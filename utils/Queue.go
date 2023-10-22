package utils

type Queue struct {
	elements []interface{}
}

func New() *Queue {
	return &Queue{nil}
}

func (q *Queue) Dequeue() interface{} {
	if q.isEmpty() {
		return nil
	}
	n := q.elements[0]
	q.elements = q.elements[1:]

	return n
}

func (q *Queue) Enqueue(value interface{}) {
	q.elements = append(q.elements, value)
}

func (q *Queue) Peek() interface{} {
	if q.isEmpty() {
		return nil
	}
	return q.elements[0]
}

func (q *Queue) isEmpty() bool {
	return len(q.elements) == 0
}

func (q *Queue) Len() int {
	return len(q.elements)
}

func (q *Queue) GetElements() []interface{} {
	return q.elements
}
