package node

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"vayeate/client"
	"vayeate/logger"
	"vayeate/sync"
	"vayeate/topic"

	"github.com/google/uuid"
)

type Node struct {
	ID         string                               `json:"id"`
	ClientPort int                                  `json:"client_port"`
	Username   string                               `json:"-"`
	Password   string                               `json:"-"`
	Clients    sync.SyncMap[string, *client.Client] `json:"-"`
	Topics     sync.SyncMap[string, *topic.Topic]   `json:"-"`

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
		Clients:        sync.NewSyncMap[string, *client.Client](),
		Topics:         sync.NewSyncMap[string, *topic.Topic](),
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
		conn.Write(client.NewErrorMessage(err.Error()).Serialize())
		conn.Close()
		return
	}

	if self.Clients.Has(c.ID) {
		c.Write(client.NewErrorMessage(fmt.Sprintf("client_id `%s` is already is use", c.ID)))
		c.Close()
		return
	}

	c.Write(client.NewConnectAckMessage(c.SessionID))
	self.Clients.Set(c.ID, c)

	defer func() {
		c.Close()
		self.Clients.Del(c.ID)
		c.Topics.ForEach(func(_ int, topic string) {
			if self.Topics.Has(topic) {
				t := self.Topics.Get(topic)
				t.UnSubscribe(c)
			}
		})
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
			p := m.GetPublishPayload()
			t := self.Topics.Get(p.Topic)

			if t == nil {
				t = topic.New(p.Topic)
				self.Topics.Set(t.Name, t)
			}

			t.Push(p.Payload)
			c.Write(client.NewPublishAckMessage())
		} else if m.Code == client.SUBSCRIBE {
			p := m.GetSubscribePayload()
			t := self.Topics.Get(p.Topic)

			if t == nil {
				t = topic.New(p.Topic)
				self.Topics.Set(t.Name, t)
			}

			t.Subscribe(c)
			c.Topics.Add(t.Name, t.Name)
			c.Write(client.NewSubscribeAckMessage())
		}
	}
}
