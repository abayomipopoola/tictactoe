package web

import (
	"sync"
)

type PlayerPool struct {
	players int
	mu      sync.Mutex
	max     int
}

func NewPlayerPool(max int) *PlayerPool {
	return &PlayerPool{
		players: 0,
		max:     max,
	}
}

func (p *PlayerPool) CanPlay() (bool, int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.players < p.max {
		p.players++
		return true, p.players
	}
	return false, p.players
}
