package connutil

type Conn interface {
	Read() ([]byte, error)
	Write(b []byte) (int, error)
	Close() error
}
