package node

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"vayeate/client"
	"vayeate/logger"

	"github.com/google/uuid"
)

type Node struct {
	ID         string `json:"id"`
	ClientPort int    `json:"client_port"`
	Username   string `json:"-"`
	Password   string `json:"-"`

	log            *logger.Logger
	clientListener net.Listener
}

func New(clientPort string, username string, password string) (*Node, error) {
	id := uuid.NewString()
	cp, err := strconv.Atoi(clientPort)

	if err != nil {
		return nil, err
	}

	cl, err := net.Listen("tcp", fmt.Sprintf(":%d", cp))

	if err != nil {
		return nil, err
	}

	self := &Node{
		ID:             id,
		ClientPort:     cp,
		Username:       username,
		Password:       password,
		log:            logger.New(fmt.Sprintf("vayeate:node:%s", id)),
		clientListener: cl,
	}

	return self, nil
}

func (self *Node) Close() {
	self.clientListener.Close()
}

func (self *Node) Listen() error {
	for {
		conn, err := self.clientListener.Accept()

		if err != nil {
			return err
		}

		go self.onClientConnection(conn)
	}
}

func (self *Node) onClientConnection(conn net.Conn) {
	c, err := client.FromConnection(self.Username, self.Password, conn)

	if err != nil {
		self.log.Warn(err)
		conn.Close()
		return
	}

	defer c.Close()

	for {
		m, err := c.Read()

		if err != nil {
			if err == io.EOF {
				return
			}

			self.log.Warnln(err)
			continue
		}

		if m.Code == client.PUBLISH {
			self.log.Infoln(*m)
			self.log.Infoln(m.GetPublishPayload())
			c.Write(client.NewPublishAckMessage())
		}
	}
}
