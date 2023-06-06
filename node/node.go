package node

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"vayeate/client"
	"vayeate/logger"
	"vayeate/sync"

	"github.com/google/uuid"
)

type Node struct {
	ID         string `json:"id"`
	ClientPort int    `json:"client_port"`
	Username   string `json:"-"`
	Password   string `json:"-"`

	log            *logger.Logger
	clientListener net.Listener
	clients        sync.SyncMap[string, *client.Client]
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
		clients:        sync.NewSyncMap[string, *client.Client](),
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

func (self *Node) GetClients() []*client.Client {
	return self.clients.Slice()
}

func (self *Node) onClientConnection(conn net.Conn) {
	c, err := client.FromConnection(self.Username, self.Password, conn)

	if err != nil {
		conn.Write(client.NewErrorMessage(err.Error()).Serialize())
		conn.Close()
		return
	}

	if self.clients.Has(c.ID) {
		c.Write(client.NewErrorMessage(fmt.Sprintf("client_id `%s` is already is use", c.ID)))
		c.Close()
		return
	}

	c.Write(client.NewConnectAckMessage(c.SessionID))
	self.clients.Set(c.ID, c)

	defer func() {
		c.Close()
		self.clients.Del(c.ID)
	}()

	for {
		m, err := c.Read()

		if err != nil {
			if err == io.EOF {
				return
			}

			self.log.Warnln(err)
			continue
		}

		if m.Code == client.PING {
			c.Write(client.NewPingAckMessage())
		} else if m.Code == client.PUBLISH {
			c.Write(client.NewPublishAckMessage())
		}
	}
}
