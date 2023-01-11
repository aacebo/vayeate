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
	port 	 int
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

func (self *Broker) Listen(handler func()) {
	for {
		conn, err := self.listener.Accept()

		if err != nil {
			log.Error(err)
		}

		go self.onConnection(conn)
	}
}

func (self *Broker) Close() {
	self.listener.Close()
}

func (self *Broker) onConnection(conn net.Conn) {
	for {
		data, err := io.ReadAll(conn)

		if err != nil {
			log.Warn(err)
			return
		}

		if len(data) < 4 {
			log.Warn("invalid frame")
			return
		}

		f := frame.Decode(data)

		if f.GetType() == "STOP" {
			break
		}

		if f.GetType() == "PING" {
			conn.Write(frame.New("PONG", nil).Encode())
		}
	}

	conn.Close()
}
