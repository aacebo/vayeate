package broker

import (
	"fmt"
	"net"
	"vayeate/frame"
	"vayeate/logger"
)

var log = logger.New("broker")

type Broker struct {
	port     int
	listener net.Listener
	sockets  map[string]*Socket
}

func New(port int) *Broker {
	sockets := make(map[string]*Socket)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Error(err)
	}

	self := Broker{port, listener, sockets}
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
	socket := NewSocket(conn)
	self.sockets[socket.GetID()] = socket
	defer socket.Close()

	for {
		if socket.GetClosed() {
			delete(self.sockets, socket.GetID())
			return
		}

		f, err := socket.Read()

		if f == nil || err != nil {
			if err != nil {
				log.Warn(err)

				if err == frame.InvalidFormatError {
					return
				}
			}

			continue
		}

		if f.IsClose() {
			return
		} else if f.IsPing() {
			err := socket.Write(frame.Pong())

			if err != nil {
				log.Warn(err)
				return
			}
		} else {
			handler(f)
		}
	}
}
