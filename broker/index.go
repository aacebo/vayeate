package broker

import (
	"fmt"
	"io"
	"net"
	"vayeate/frame"
	"vayeate/logger"
)

var log = logger.New("broker")

type Broker struct {
	port     int
	listener net.Listener
}

func New(port int) *Broker {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Error(err)
	}

	self := Broker{port, listener}
	return &self
}

func (self *Broker) Listen(handler func(f *frame.Frame)) {
	for {
		conn, err := self.listener.Accept()

		if err != nil {
			log.Error(err)
		}

		go self.onConnection(conn, handler)
	}
}

func (self *Broker) Close() {
	self.listener.Close()
}

func (self *Broker) onConnection(conn net.Conn, handler func(f *frame.Frame)) {
	defer conn.Close()

	for {
		data, err := io.ReadAll(conn)

		if err != nil {
			log.Warn(err)
			return
		}

		if len(data) == 0 {
			continue
		}

		if len(data) < 4 {
			log.Warn("invalid frame")
			return
		}

		f, err := frame.Decode(data)

		if err != nil {
			log.Warn(err)
			return
		}

		if f.IsClose() {
			break
		} else if f.IsPing() {
			conn.Write(frame.Pong().Encode())
		} else {
			handler(f)
		}
	}
}
