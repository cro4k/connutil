package readutil

import (
	"encoding/binary"
	"io"
)

func SimpleDecode(reader io.Reader) ([]byte, error) {
	head := make([]byte, 4)
	_, err := reader.Read(head)
	if err != nil {
		return nil, err
	}
	length := int(binary.BigEndian.Uint32(head))
	var data = make([]byte, 0, length)
	var read int
	for read < length {
		buf := make([]byte, 2048)
		n, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			break
		}
		read += n
		data = append(data, buf[:n]...)
	}
	return data, nil
}
