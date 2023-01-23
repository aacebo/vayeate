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

type Server struct {
	ID   string
	Port int

	log      *logger.Logger
	listener net.Listener
	sockets  map[string]*Socket
	queues   map[string]*queue.Queue
}

func New(port int) (*Server, error) {
	id := uuid.NewString()
	sockets := map[string]*Socket{}
	queues := map[string]*queue.Queue{}
	log := logger.New(fmt.Sprintf("server:%s", id))
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	return &Server{
		id,
		port,
		log,
		listener,
		sockets,
		queues,
	}, nil
}

func (self *Server) Listen(
	onConnect func(s *Socket),
	onFrame func(s *Socket, f *frame.Frame),
	onError func(err error),
) error {
	for {
		conn, err := self.listener.Accept()

		if err != nil {
			return err
		}

		go self.onConnection(conn, onConnect, onFrame, onError)
	}
}

func (self *Server) Close() {
	self.listener.Close()
}

func (self *Server) GetQueue(name string) *queue.Queue {
	return self.queues[name]
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

func (self *Server) AddQueue(q *queue.Queue) *queue.Queue {
	if self.queues[q.Name] != nil {
		return self.queues[q.Name]
	}

	self.queues[q.Name] = q
	return q
}

func (self *Server) onConnection(
	conn net.Conn,
	onConnect func(s *Socket),
	onFrame func(s *Socket, f *frame.Frame),
	onError func(err error),
) {
	s := NewSocket(conn)
	self.sockets[s.ID] = s

	if onConnect != nil {
		onConnect(s)
	}

	defer func() {
		s.Close()
		delete(self.sockets, s.ID)
	}()

	for {
		if s.Closed == true {
			return
		}

		f, err := s.Read()

		if f == nil || err != nil {
			if err != nil {
				self.log.Warn(err)

				if onError != nil {
					onError(err)
				}

				if err == frame.InvalidFormatError {
					return
				}
			}

			continue
		}

		if f.IsClose() {
			return
		} else if f.IsPing() {
			self.onPing(s)
		} else if onFrame != nil {
			onFrame(s, f)
		}
	}
}

func (self *Server) onPing(s *Socket) {
	err := s.Write(frame.NewPong())

	if err != nil {
		self.log.Warn(err)
		return
	}
}
