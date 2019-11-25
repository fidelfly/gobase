package gox

import "errors"

type StateSignal chan int

func NewStateSignal() StateSignal {
	return make(StateSignal)
}

func (ss StateSignal) SendSignal(state int) {
	ss <- state
}

func (ss StateSignal) Wait(state int) error {
	for {
		s, ok := <-ss
		if !ok {
			return errors.New("something is wrong")
		}
		if s == state {
			return nil
		}
	}
}
