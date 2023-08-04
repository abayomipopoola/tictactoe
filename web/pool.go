package web

import (
	"fmt"
	"sync"
	"time"
)

var tick = []string{"-1", "X", "O"}

type PlayerPool struct {
	play    map[string]string
	players int
	mu      sync.Mutex
	max     int
}

func NewPlayerPool(max int) *PlayerPool {
	return &PlayerPool{
		play:    make(map[string]string, 2),
		players: 0,
		max:     max,
	}
}

func (p *PlayerPool) CanPlay() (bool, string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.players < p.max {
		p.players++
		c := p.players
		p.play[tick[c]] = fmt.Sprintf("%d", time.Now().Unix())
		return true, tick[c]
	}
	return false, tick[0]
}
