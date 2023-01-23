package server

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

// timeout connections after 60s of inactivity
const timeout = 60 * time.Second

type Socket struct {
	ID        string
	Closed    bool
	StartedAt int64

	pingTimer *time.Timer
	log       *logger.Logger
	reader    *bufio.Reader
	conn      net.Conn
}

func NewSocket(conn net.Conn) *Socket {
	id := uuid.NewString()
	closed := false
	reader := bufio.NewReader(conn)
	now := time.Now().Unix()
	log := logger.New(fmt.Sprintf("socket:%s", id))
	self := Socket{id, closed, now, nil, log, reader, conn}
	self.pingTimer = time.AfterFunc(timeout, onTimeout(&self))

	return &self
}

func (self *Socket) Close() {
	self.Closed = true
	self.conn.Close()
	self.pingTimer.Stop()
}

func (self *Socket) GetRemoteAddress() string {
	return self.conn.RemoteAddr().String()
}

func (self *Socket) Read() (*frame.Frame, error) {
	f, err := frame.Decode(self.reader)

	if err == io.EOF {
		return nil, nil
	}

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
