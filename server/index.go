package server

import (
	"fmt"
	"net"
	"regexp"

	"vayeate/frame"
	"vayeate/logger"
	"vayeate/queue"

	"github.com/google/uuid"
)

var log = logger.New("server")

type Server struct {
	ID   string
	Port int

	listener net.Listener
	sockets  map[string]*Socket
	queues   map[string]*queue.Queue
}

func New(port int) (*Server, error) {
	id := uuid.NewString()
	sockets := map[string]*Socket{}
	queues := map[string]*queue.Queue{}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	return &Server{
		id,
		port,
		listener,
		sockets,
		queues,
	}, nil
}

func (self *Server) Listen() error {
	for {
		conn, err := self.listener.Accept()

		if err != nil {
			return err
		}

		go self.onConnection(conn)
	}
}

func (self *Server) Close() {
	self.listener.Close()
}

func (self *Server) GetQueues(pattern string) []*queue.Queue {
	queues := []*queue.Queue{}

	for key, q := range self.queues {
		match, _ := regexp.MatchString(pattern, key)

		if match == true {
			queues = append(queues, q)
		}
	}

	return queues
}

func (self *Server) onConnection(conn net.Conn) {
	socket := NewSocket(conn)
	self.sockets[socket.ID] = socket
	defer socket.Close()

	for {
		if socket.Closed == true {
			delete(self.sockets, socket.ID)
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
			q := queue.New(f.GetSubject())
			self.queues[q.Name] = q
			log.Infoln(len(self.queues))
		} else if f.IsProduce() {
			queues := self.GetQueues(f.GetSubject())

			for _, q := range queues {
				q.Push(f.Body)
			}
		} else if f.IsConsume() {
			queues := self.GetQueues(f.GetSubject())

			for _, q := range queues {
				q.Consume(socket.ID)
			}
		}
	}
}
