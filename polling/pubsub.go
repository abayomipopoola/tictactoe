package polling

import "sync"

type PubSub struct {
	channels map[chan struct{}]struct{}
	lock     sync.RWMutex
}

func NewPubSub() *PubSub {
	return &PubSub{
		channels: make(map[chan struct{}]struct{}),
	}
}

func (p *PubSub) Subscribe() (<-chan struct{}, func()) {
	p.lock.Lock()
	defer p.lock.Unlock()

	c := make(chan struct{}, 1)
	p.channels[c] = struct{}{}
	return c, func() {
		p.lock.Lock()
		defer p.lock.Unlock()

		delete(p.channels, c)
		close(c)
	}
}

func (p *PubSub) Publish() {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for channel := range p.channels {
		channel <- struct{}{}
	}
}
