package golimit

import "sync"

type GoLimit struct {
	ch chan int
}

func NewGoLimit(max int) *GoLimit {
	return &GoLimit{ch: make(chan int, max)}
}

func (g *GoLimit) Add() {
	g.ch <- 1
}

func (g *GoLimit) Done() {
	<-g.ch
}

type GoLimitPool struct {
	lock  sync.Mutex
	chmap map[string]*GoLimit
}

func NewGoLimitPool() *GoLimitPool {
	return &GoLimitPool{
		lock:  sync.Mutex{},
		chmap: make(map[string]*GoLimit),
	}
}

func (gp *GoLimitPool) GetGoLimit(addr string, max int) *GoLimit {
	gp.lock.Lock()
	defer gp.lock.Unlock()
	if _, ok := gp.chmap[addr]; !ok {
		gp.chmap[addr] = &GoLimit{ch: make(chan int, max)}
	}
	return gp.chmap[addr]
}
