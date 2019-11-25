package containerx

import (
	"fmt"
	"strings"
	"sync"
)

const defaultQueueCapacity = 10

type Queue interface {
	Len() int
	Dequeue() interface{}
	Enqueue(interface{})
	Peek() interface{}
}

type circleQueue struct {
	data   []interface{}
	head   int
	tail   int
	length int
	lock   sync.RWMutex
}

func (cq *circleQueue) Len() int {
	return len(cq.data)
}

func (cq *circleQueue) Dequeue() interface{} {
	if cq.length == 0 {
		return nil
	}
	cq.lock.Lock()
	defer cq.lock.Unlock()
	v := cq.data[cq.head]
	cq.data[cq.head] = nil
	cq.head = (cq.head + 1) % cap(cq.data)
	cq.length--
	return v
}

func (cq *circleQueue) Enqueue(item interface{}) {
	cq.lock.Lock()
	defer cq.lock.Unlock()
	cq._ensureCapacity()
	if len(cq.data) < cap(cq.data) {
		cq.data = append(cq.data, item)
		cq.tail = len(cq.data) - 1
	} else {
		newIndex := (cq.tail + 1) % cap(cq.data)
		cq.data[newIndex] = item
		cq.tail = newIndex
	}
	cq.length++
}

func (cq *circleQueue) Peek() interface{} {
	if cq.length == 0 {
		return nil
	}
	cq.lock.RLock()
	defer cq.lock.RUnlock()
	return cq.data[cq.head]
}

func (cq *circleQueue) _ensureCapacity() {
	if cq.length == cap(cq.data) {
		newCap := 2 * cap(cq.data)
		if newCap == 0 {
			newCap = defaultQueueCapacity
		}

		orgs := cq.data
		cq.data = make([]interface{}, 0, newCap)
		if cq.tail >= cq.head {
			cq.data = append(cq.data, orgs[cq.head:(cq.tail+1)]...)
		} else {
			cq.data = append(cq.data, orgs[cq.head:cq.length]...)
			cq.data = append(cq.data, orgs[:(cq.tail+1)]...)
		}
		cq.head = 0
		cq.tail = cq.length - 1
	}
}

func (cq *circleQueue) String() string {
	str := strings.Builder{}
	str.Write([]byte("["))
	for n := 0; n < cq.length; n++ {
		if n > 0 {
			str.Write([]byte(","))
		}
		index := (cq.head + n) % cap(cq.data)
		str.WriteString(fmt.Sprintf("%v", cq.data[index]))
	}
	str.Write([]byte("]"))
	return str.String()
}

func (cq *circleQueue) DebugString() string {
	return fmt.Sprintf("len=%d cap=%d %v", cq.length, cap(cq.data), cq)
}

func NewQueue(capacity int) Queue {
	if capacity == 0 {
		capacity = defaultQueueCapacity
	}
	return &circleQueue{
		data:   make([]interface{}, 0, capacity),
		head:   0,
		tail:   0,
		length: 0,
		lock:   sync.RWMutex{},
	}
}
