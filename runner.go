package connutil

import (
	"errors"
	"time"
)

type Option interface {
	apply(r *runner)
}

type OptionFunc func(r *runner)

func (f OptionFunc) apply(r *runner) {
	f(r)
}

func WithTimeout(timeout time.Duration) OptionFunc {
	return func(r *runner) {
		r.timeout = timeout
	}
}

func WithCallback(callback Callback) OptionFunc {
	return func(r *runner) {
		r.callback = callback
	}
}

type runner struct {
	id        string
	connChan  chan Conn
	writeChan chan []byte
	timeout   time.Duration
	callback  Callback
	shutdowns chan chan struct{}
}

func (r *runner) run() {
	for {
		select {
		case conn := <-r.connChan:
			r.callback.OnConnected(r.id)
			shutdown := make(chan struct{})
			r.shutdowns <- shutdown
			r.do(conn, shutdown)
			r.callback.OnDisconnected(r.id)
		case <-time.After(time.Second * 5):
			r.callback.OnShutdown(r.id)
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
				r.callback.OnMessage(r.id, msg)
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

func (r *runner) shutdown() error {
	select {
	case sh := <-r.shutdowns:
		close(sh)
		return nil
	//case <-time.After(time.Second):
	default:
		return errors.New("shutdown failed, not running")
	}
}

func newRunner(id string, opt ...Option) *runner {
	r := &runner{
		id:        id,
		connChan:  make(chan Conn),
		writeChan: make(chan []byte),
		timeout:   time.Second * 15,
		shutdowns: make(chan chan struct{}),
		callback:  &EmptyCallback{},
	}
	for _, o := range opt {
		o.apply(r)
	}
	return r
}
