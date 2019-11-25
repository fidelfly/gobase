package containerx

import (
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack(5)
	priceStack(t, s)

	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	priceStack(t, s)

	t.Logf("Peek() = %v\n", s.Peek())
	v := s.Pop()
	t.Logf("Pop() = %v\n", v)
	priceStack(t, s)

	s.Push(4)
	priceStack(t, s)
	s.Push(5)
	priceStack(t, s)
	s.Push(6)
	priceStack(t, s)
	s.Push(7)
	priceStack(t, s)

	t.Logf("%v", s)
}

func priceStack(t *testing.T, s Stack) {
	if cs, ok := s.(*comStack); ok {
		t.Logf("%s\n", cs.DebugString())
	} else {
		t.Logf("Len() = %d %v\n", s.Len(), s)
	}
}
