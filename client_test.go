package connutil

import (
	"testing"
	"time"
)

type mockConn struct {
	data []byte
}

func (c *mockConn) Read() ([]byte, error) {
	time.Sleep(time.Second)
	return c.data, nil
}

func (c *mockConn) Write(b []byte) (int, error) {
	c.data = b
	return len(b), nil
}

func (c *mockConn) Close() error {
	return nil
}

func TestClient(t *testing.T) {
	conn := &mockConn{}
	listener := NewListener(
		func(c *Client) {
			t.Log(c.id, " connected")
		},
		func(c *Client) {
			t.Log(c.id, " disconnected")
		},
		func(c *Client) {
			t.Log(c.id, " reconnected")
		},
		func(c *Client, data []byte) {
			t.Log(c.id, " message", len(data))
		},
		func(c *Client) {
			t.Log(c.id, " removed")
		})
	c, _ := NewClient("1", conn, listener)
	go func(client *Client) {
		time.Sleep(time.Second * 5)
		client.Write([]byte("hello"))
	}(c)
	c.Run()
}
