package containerx

import (
	"testing"
)

func TestNewQueue(t *testing.T) {
	q := NewQueue(3)

	q.Enqueue(0)
	printQueue(t, q)
	q.Enqueue(1)
	printQueue(t, q)
	q.Enqueue(3)
	printQueue(t, q)
	t.Logf("Peek() = %v", q.Peek())
	t.Logf("Dequeue() = %v", q.Dequeue())
	printQueue(t, q)

	t.Logf("Dequeue() = %v", q.Dequeue())
	printQueue(t, q)

	q.Enqueue(4)
	printQueue(t, q)
	q.Enqueue(5)
	printQueue(t, q)
	q.Enqueue(6)
	printQueue(t, q)

}

func printQueue(t *testing.T, s Queue) {
	if cs, ok := s.(*circleQueue); ok {
		t.Logf("%s\n", cs.DebugString())
	} else {
		t.Logf("Len() = %d %v\n", s.Len(), s)
	}
}
