package connutil

import (
	"fmt"
	"sync"
	"time"
)

type runner struct {
	id        string
	connChan  chan Conn
	writeChan chan []byte
	timeout   time.Duration

	onMessage func(id string, data []byte)

	shutdowns chan chan struct{}
}

func (r *runner) Run() {
	for {
		select {
		case conn := <-r.connChan:
			fmt.Printf("%s (re)connected\n", r.id)
			shutdown := make(chan struct{})
			r.shutdowns <- shutdown
			r.do(conn, shutdown)
			fmt.Printf("%s disconnected\n", r.id)
		case <-time.After(time.Second * 5):
			fmt.Printf("%s shutdown", r.id)
			return
		}
	}
}

func (r *runner) replace(conn Conn) {
	select {
	case shutdown := <-r.shutdowns:
		close(shutdown)
	default:
	}
	r.connChan <- conn
}

func (r *runner) write(b []byte) {
	r.writeChan <- b
}

func (r *runner) do(conn Conn, shutdown chan struct{}) {
	defer conn.Close()
	var res = make(chan []byte)
	go r.accept(conn, res)
	for {
		select {
		case w := <-r.writeChan:
			_, _ = conn.Write(w)
		case msg, ok := <-res:
			if ok {
				r.onMessage(r.id, msg)
			} else {
				return
			}
		case <-time.After(r.timeout):
			return
		case <-shutdown:
			return
		}
	}
}

func (r *runner) accept(conn Conn, ch chan<- []byte) {
	defer close(ch)
	for {
		if data, err := conn.Read(); err != nil {
			break
		} else {
			ch <- data
		}
	}
}

type Pool struct {
	mu      *sync.RWMutex
	runners map[string]*runner
}

func NewPool() *Pool {
	return &Pool{mu: new(sync.RWMutex), runners: make(map[string]*runner)}
}

func (p *Pool) Join(id string, conn Conn, onMessage func(id string, msg []byte), timeout time.Duration) {
	r := p.get(id, onMessage, timeout)
	r.replace(conn)
}

func (p *Pool) get(id string, onMessage func(id string, msg []byte), timeout time.Duration) *runner {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.runners[id] == nil {
		r := &runner{
			id:        id,
			connChan:  make(chan Conn),
			writeChan: make(chan []byte),
			timeout:   timeout,
			onMessage: onMessage,
			shutdowns: make(chan chan struct{}, 4),
		}
		go r.Run()
		p.runners[id] = r
	}
	return p.runners[id]
}

func (p *Pool) Write(id string, data []byte) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if r := p.runners[id]; r != nil {
		r.write(data)
	}
}
