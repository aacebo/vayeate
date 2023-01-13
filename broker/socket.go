package broker

import (
	"bufio"
	"io"
	"net"

	"vayeate/frame"

	"github.com/google/uuid"
)

const delimiter = byte(':')

type Socket struct {
	id     string
	conn   net.Conn
	reader *bufio.Reader
}

func NewSocket(conn net.Conn) *Socket {
	id := uuid.NewString()
	reader := bufio.NewReader(conn)
	self := Socket{id, conn, reader}
	return &self
}

func (self *Socket) GetID() string {
	return self.id
}

func (self *Socket) Close() {
	self.conn.Close()
}

func (self *Socket) Read() (*frame.Frame, error) {
	data := []byte{}

	for {
		chunk, err := self.reader.ReadBytes(byte(':'))
		data = append(data, chunk...)

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	if len(data) == 0 {
		return nil, nil
	}

	f, err := frame.Decode(data)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (self *Socket) Write(f *frame.Frame) error {
	_, err := self.conn.Write(f.Encode())
	return err
}
