package randx

import (
	"math/rand"
	"time"
)

type Pool struct {
	ss     []Source
	queue  []Source
	cw     []int // basis factor of each source, it's used to determine which source would be the next element in queue
	tw     int   // total weight of all sources
	rsi    int   // random start index
	len    int   // length of queue excluding nil
	myRand *rand.Rand
}

type Source1 struct {
	G Generator
	w int
}

func (s *Source1) Get() interface{} {
	return s.G.Get()
}

func (s Source1) Weight() int {
	return s.w
}

type Source interface {
	Generator
	Weight() int
}

type Generator interface {
	Get() interface{}
}

func NewPool(params ...interface{}) *Pool {
	ss := make([]Source, 0)
	var s *Source1
	var addTemporarySource = func() {
		if s != nil {
			ss = append(ss, s)
			s = nil
		}
	}
	for _, param := range params {
		switch v := param.(type) {
		case Source:
			addTemporarySource()
			ss = append(ss, v)
		case Generator:
			s = &Source1{
				G: v,
				w: 1,
			}
		case int:
			if s != nil {
				s.w = v
				addTemporarySource()
			}
		default:
			addTemporarySource()
		}
	}

	p := &Pool{ss: ss, myRand: rand.New(rand.NewSource(time.Now().Unix()))}
	p.resetQueue()
	return p
}

func (p *Pool) RandInt(n int) int {
	return p.myRand.Intn(n)
}

func (p *Pool) resetQueue() {
	p.tw = 0
	for _, s := range p.ss {
		p.tw += s.Weight()
	}
	p.queue = make([]Source, p.tw)
	p.cw = make([]int, len(p.ss))
	p.len = 0
	p.ResetIndex()
}

func (p *Pool) ResetIndex() {
	if p.len == 0 {
		p.rsi = p.myRand.Intn(len(p.ss))
	} else {
		p.rsi = p.myRand.Intn(p.len)
	}

	//p.rsi = rand.Int() % len(p.ss)
}

func (p *Pool) Get() interface{} {
	if p.queue[p.rsi] == nil || (p.rsi < p.tw-1 && p.queue[p.rsi+1] == nil) {
		p.fillQueue(len(p.ss))
	}
	s := p.queue[p.rsi]
	p.rsi++
	if p.rsi == p.tw {
		p.rsi = 0
	}
	return s.Get()
}

func (p *Pool) fillQueue(n int) {
	for n > 0 {
		if p.len == p.tw {
			break
		}
		chosen := -1
		for i := 0; i < len(p.ss); i++ {
			s := p.ss[i]
			p.cw[i] += s.Weight()
			if chosen < 0 || p.cw[i] > p.cw[chosen] {
				chosen = i
			}
		}

		p.queue[p.len] = p.ss[chosen]
		p.len++
		p.cw[chosen] -= p.tw
		n--
	}
}
