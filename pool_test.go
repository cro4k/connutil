package connutil

import (
	"fmt"
	"testing"
	"time"
)

type IDConn struct {
	id string
}

func (c *IDConn) Read() ([]byte, error) {
	time.Sleep(time.Second)
	return []byte(c.id), nil
}

func (c *IDConn) Write(b []byte) (int, error) {
	return len(b), nil
}

func (c *IDConn) Close() error {
	return nil
}

func onMessage(id string, msg []byte) {
	fmt.Printf("%s: %s\n", id, string(msg))
}

func TestPool(t *testing.T) {
	p := NewPool()
	conn1 := &IDConn{id: "111"}
	conn2 := &IDConn{id: "222"}

	p.Join("1", conn1, onMessage, time.Second*5)
	time.Sleep(time.Second * 5)
	p.Join("1", conn2, onMessage, time.Second*5)
	time.Sleep(time.Second * 5)
}
