package broker

import (
	"fmt"
	"net"
	"vayeate/frame"
	"vayeate/logger"
	"vayeate/queue"
)

var log = logger.New("broker")

type Broker struct {
	port     int
	listener net.Listener
	sockets  map[string]*Socket
	queues   map[string]*queue.Queue
}

func New(port int) *Broker {
	sockets := map[string]*Socket{}
	queues := map[string]*queue.Queue{}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Error(err)
	}

	self := Broker{port, listener, sockets, queues}
	return &self
}

func (self *Broker) Listen() {
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
			err := socket.Write(frame.NewPong())

			if err != nil {
				log.Warn(err)
				return
			}
		} else if f.IsAssert() {
			q := queue.New(f.GetBody())
			self.queues[q.ID] = q
		}
	}
}
