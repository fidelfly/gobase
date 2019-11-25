package containerx

import (
	"fmt"
	"sync"
)

const defaultStackCapacity = 5

type Stack interface {
	Len() int
	Pop() interface{}
	Push(interface{})
	Peek() interface{}
}

type comStack struct {
	data []interface{}
	lock sync.RWMutex
}

func (cs *comStack) Len() int {
	return len(cs.data)
}

func (cs *comStack) Pop() interface{} {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	sLen := cs.Len()
	if sLen == 0 {
		return nil
	}
	v := cs.data[sLen-1]
	cs.data[sLen-1] = nil
	cs.data = cs.data[:(sLen - 1)]
	return v
}

func (cs *comStack) Push(item interface{}) {
	cs.lock.Lock()
	defer cs.lock.Unlock()
	cs._ensureCapacity()
	cs.data = append(cs.data, item)
}

func (cs *comStack) Peek() interface{} {
	sLen := cs.Len()
	if sLen == 0 {
		return nil
	}
	return cs.data[sLen-1]
}

func (cs *comStack) _ensureCapacity() {
	if len(cs.data) == cap(cs.data) {
		newCap := 2 * cap(cs.data)
		if newCap == 0 {
			newCap = defaultStackCapacity
		}

		orgs := cs.data
		cs.data = make([]interface{}, len(orgs), newCap)
		copy(cs.data, orgs)
	}
}

func (cs *comStack) String() string {
	return fmt.Sprintf("%v", cs.data)
}

func (cs *comStack) DebugString() string {
	return fmt.Sprintf("len=%d cap=%d %v", len(cs.data), cap(cs.data), cs.data)
}

func NewStack(capacity int) Stack {
	if capacity == 0 {
		capacity = defaultStackCapacity
	}
	return &comStack{
		data: make([]interface{}, 0, capacity),
	}
}
