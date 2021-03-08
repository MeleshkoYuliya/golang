package notifier

import (
	"sync"
)

type PubSub interface {
	Subscribe(topic int) <-chan interface{}
	Publish(topic int, data interface{})
	Close()
}

type pubsub struct {
	mu     sync.RWMutex
	subs   map[int][]chan interface{}
	closed bool
}

func NewPubSub() *pubsub {
	ps := &pubsub{}
	ps.subs = make(map[int][]chan interface{})
	return ps
}

func (ps *pubsub) Subscribe(topic int) <-chan interface{} {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan interface{}, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *pubsub) Publish(topic int, data interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		go func(ch chan interface{}) {
			ch <- data
		}(ch)
	}
}

func (ps *pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}
