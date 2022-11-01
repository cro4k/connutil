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

func NewClient(id string, conn Conn, onNewClient func(c *Client), timeout ...time.Duration) *Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	var c *Client
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
		go c.ping(c.hb)
		m.clients[id] = c
		if onNewClient != nil {
			onNewClient(c)
		}
	}
	return c
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
