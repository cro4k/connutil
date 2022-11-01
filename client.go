package connutil

import (
	"sync"
	"time"
)

type Client struct {
	id      string
	conn    Conn
	mu      *sync.Mutex
	closed  bool
	hb      chan struct{}
	timeout time.Duration

	listener Listener
}

func (c *Client) ping(hb chan struct{}) {
	c.listener.OnConnected(c.id)
	defer c.listener.OnDisconnected(c.id)
	for {
		select {
		case <-hb:
		case <-time.After(c.timeout):
			c.listener.OnRemove(c.id)
			remove(c.id)
			return
		}
	}
}

func (c *Client) receive(conn Conn, hb chan struct{}) {
	defer conn.Close()
	for {
		data, err := conn.Read()
		if err != nil {
			return
		} else {
			c.listener.OnMessage(c.id, data)
			hb <- struct{}{}
		}
	}
}

func (c *Client) reset(conn Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.conn.Close()
	}
	c.closed = false
	c.conn = conn
	c.listener.OnReconnected(c.id)
}

func (c *Client) SetListener(listener Listener) {
	if c.listener == listener {
		return
	}
	c.listener = listener
}

func (c *Client) Write(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.conn.Write(data)
	return err
}

func (c *Client) Run() {
	c.receive(c.conn, c.hb)
}

func newClient(id string, conn Conn, timeout time.Duration) *Client {
	return &Client{
		id:       id,
		conn:     conn,
		mu:       new(sync.Mutex),
		closed:   false,
		hb:       make(chan struct{}),
		timeout:  timeout,
		listener: DefaultListener,
	}
}