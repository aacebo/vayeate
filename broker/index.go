package broker

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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
		data, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			log.Warn(err)
			return
		}

		trimmed := strings.TrimSpace(string(data))

		if trimmed == "STOP" {
			break
		}

		if trimmed == "PING" {
			conn.Write([]byte("PONG"))
		}
	}

	conn.Close()
}
