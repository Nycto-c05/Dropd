package idgen

import (
	"sync"
	"time"
)

type IDGenerator interface {
	Next() int64
}

type Generator struct {
	lastTs  int64
	counter int64
	mtx     sync.Mutex
}

func NewIDGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Next() int64 {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	now := time.Now().UnixMilli()

	if now == g.lastTs {
		g.counter++
	} else {
		g.lastTs = now
		g.counter = 0
	}

	//handles 4096 ID generations in a millisecond
	return (now << 12) | g.counter //Alex Xu - Design a Distributed ID Generator

}
