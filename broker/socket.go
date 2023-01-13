package broker

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"vayeate/frame"
	"vayeate/logger"

	"github.com/google/uuid"
)

const timeout = 60 * time.Second

type Socket struct {
	id         *string
	closed     *bool
	reader     *bufio.Reader
	startedAt  *int64
	pingTimer  *time.Timer
	log		   *logger.Logger
	conn       net.Conn
}

func NewSocket(conn net.Conn) *Socket {
	id := uuid.NewString()
	closed := false;
	reader := bufio.NewReader(conn)
	now := time.Now().Unix()
	log := logger.New(fmt.Sprintf("socket:%s", id))
	self := Socket{&id, &closed, reader, &now, nil, log, conn}
	self.pingTimer = time.AfterFunc(timeout, onTimeout(&self))
	log.Debugln("connected")
	return &self
}

func (self *Socket) GetID() string {
	return *self.id
}

func (self *Socket) GetClosed() bool {
	return *self.closed
}

func (self *Socket) Close() {
	*self.closed = true
	self.conn.Close()
}

func (self *Socket) Read() (*frame.Frame, error) {
	data := []byte{}

	for {
		chunk, err := self.reader.ReadBytes(frame.Delimiter)
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

	self.pingTimer.Reset(timeout)
	return f, nil
}

func (self *Socket) Write(f *frame.Frame) error {
	_, err := self.conn.Write(f.Encode())
	return err
}

func onTimeout(self *Socket) func() {
	return func() {
		self.Close()
		self.log.Debugln("closed due to inactivity")
	}
}
