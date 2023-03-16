package connutil

import (
	"errors"
	"sync"
)

type Pool struct {
	mu      *sync.RWMutex
	runners map[string]*runner
	limit   chan struct{}
}

func NewPool(size int) *Pool {
	return &Pool{mu: new(sync.RWMutex), runners: make(map[string]*runner), limit: make(chan struct{}, size)}
}

func (p *Pool) Join(id string, conn Conn, opt ...Option) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	r := p.runners[id]
	if r == nil {
		r = newRunner(id, opt...)
		select {
		case p.limit <- struct{}{}:
			go p.run(id, r)
		default:
			return errors.New("connection up to limit")
		}
	}
	r.replace(conn)
	return nil
}

func (p *Pool) Size() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.runners)
}

func (p *Pool) remove(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.runners, id)
}

func (p *Pool) run(id string, r *runner) {
	defer p.remove(id)
	r.run()
	<-p.limit
}

func (p *Pool) Write(id string, data []byte) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if r := p.runners[id]; r != nil {
		r.write(data)
	}
}

func (p *Pool) Shutdown(id string) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if r := p.runners[id]; r != nil {
		return r.shutdown()
	}
	return errors.New("not found")
}
