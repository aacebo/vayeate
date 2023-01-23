package server

import (
	"fmt"
	"net"

	"vayeate/common"
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
	sockets  *common.SyncMap[*Socket]
	queues   *common.SyncMap[*queue.Queue]
}

func New(port int) (*Server, error) {
	id := uuid.NewString()
	sockets := common.NewSyncMap[*Socket]()
	queues := common.NewSyncMap[*queue.Queue]()
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
	return self.queues.Get(name)
}

func (self *Server) GetQueues(pattern string) []*queue.Queue {
	return self.queues.GetMatching(pattern)
}

func (self *Server) AddQueue(q *queue.Queue) *queue.Queue {
	existing := self.queues.Get(q.Name)

	if existing != nil {
		return existing
	}

	self.queues.Set(q.Name, q)
	go q.Start(self.onConsume(q))
	return q
}

func (self *Server) onConnection(
	conn net.Conn,
	onConnect func(s *Socket),
	onFrame func(s *Socket, f *frame.Frame),
	onError func(err error),
) {
	s := NewSocket(conn)
	self.sockets.Set(s.ID, s)

	if onConnect != nil {
		onConnect(s)
	}

	defer func() {
		s.Close()
		self.sockets.Del(s.ID)
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

func (self *Server) onConsume(q *queue.Queue) func(id string, payload []byte) {
	return func(id string, payload []byte) {
		s := self.sockets.Get(id)

		if s == nil {
			return
		}

		f := frame.NewDelegate([]byte(q.Name), payload)
		s.Write(f)
	}
}
