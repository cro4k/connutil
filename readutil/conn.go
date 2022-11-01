package readutil

import "io"

type Conn struct {
	conn    io.ReadWriteCloser
	decoder func(reader io.Reader) ([]byte, error)
}

func (c *Conn) Read() ([]byte, error) {
	return c.decoder(c.conn)
}

func (c *Conn) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func NewConn(conn io.ReadWriteCloser, decoder func(io.Reader) ([]byte, error)) *Conn {
	return &Conn{
		conn:    conn,
		decoder: decoder,
	}
}
