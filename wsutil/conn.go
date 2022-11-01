package wsutil

import "golang.org/x/net/websocket"

type BytesConn struct {
	c *websocket.Conn
}

func NewBytesConn(c *websocket.Conn) *BytesConn {
	return &BytesConn{c: c}
}

func (c *BytesConn) Read() ([]byte, error) {
	var b []byte
	err := websocket.Message.Receive(c.c, &b)
	return b, err
}

func (c *BytesConn) Write(b []byte) (int, error) {
	err := websocket.Message.Send(c.c, b)
	return len(b), err
}

func (c *BytesConn) Close() error {
	return c.c.Close()
}

type StringConn struct {
	c *websocket.Conn
}

func NewStringConn(c *websocket.Conn) *StringConn {
	return &StringConn{c: c}
}

func (c *StringConn) Read() ([]byte, error) {
	var s string
	err := websocket.Message.Receive(c.c, &s)
	return []byte(s), err
}

func (c *StringConn) Write(b []byte) (int, error) {
	err := websocket.Message.Send(c.c, string(b))
	return len(b), err
}

func (c *StringConn) Close() error {
	return c.c.Close()
}
