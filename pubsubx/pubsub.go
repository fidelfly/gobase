package pubsubx

import (
	"github.com/cskr/pubsub"

	"github.com/fidelfly/gox/logx"
)

type PubXSub struct {
	ps *pubsub.PubSub
}

type SubscriberHandler func(msg interface{}) error
type Subscriber chan interface{}

func New(capacity int) *PubXSub {
	return &PubXSub{ps: pubsub.New(capacity)}
}

func (pxs *PubXSub) Subscribe(topic string, handlers ...SubscriberHandler) Subscriber {
	subch := pxs.ps.Sub(topic)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.Error(err)
			}
		}()
		for {
			if msg, ok := <-subch; ok {
				for _, handler := range handlers {
					_ = handler(msg)
				}
			} else {
				break
			}
		}

	}()
	return subch
}

func (pxs *PubXSub) UnSubscribe(sub Subscriber, topics ...string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.Error(err)
			}
		}()
		pxs.ps.Unsub(sub, topics...)
	}()
}

func (pxs *PubXSub) Publish(msg interface{}, topics ...string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.Error(err)
			}
		}()
		pxs.ps.Pub(msg, topics...)
	}()
}
