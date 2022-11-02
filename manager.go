package connutil

import (
	"errors"
	"sync"
	"time"
)

var m struct {
	mu      *sync.RWMutex
	clients map[string]*Client
}

func init() {
	m.mu = new(sync.RWMutex)
	m.clients = make(map[string]*Client)
}

func NewClient(id string, conn Conn, listener Listener, timeout ...time.Duration) (*Client, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var c *Client
	var isNew bool
	if c = m.clients[id]; c != nil {
		c.reset(conn)
	} else {
		var t time.Duration
		if len(timeout) > 0 && timeout[0] > 0 {
			t = timeout[0]
		} else {
			t = time.Minute
		}
		c = newClient(id, conn, t)
		c.setListener(listener)
		go c.ping(c.hb)
		m.clients[id] = c
		isNew = true
	}
	return c, isNew
}

func GetClient(id string) (*Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if c := m.clients[id]; c != nil {
		return c, nil
	}
	return nil, errors.New("not found")
}

func remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, id)
}

func Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.clients)
}
