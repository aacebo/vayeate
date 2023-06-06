package client

import (
	"bufio"
	"errors"
	"net"

	"github.com/google/uuid"
)

type Client struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`

	conn   net.Conn
	reader *bufio.Reader
}

func FromConnection(username string, password string, conn net.Conn) (*Client, error) {
	reader := bufio.NewReader(conn)
	m, err := ReadMessage(reader)

	if err != nil {
		return nil, err
	}

	if m.Code != CONNECT {
		return nil, errors.New("client: first message must be of type `CONNECT`")
	}

	payload := m.GetConnectPayload()

	if payload.Username != username || payload.Password != password {
		return nil, errors.New("client: unauthorized")
	}

	sessionId := uuid.NewString()
	conn.Write(NewConnectAckMessage(sessionId).Serialize())

	return &Client{
		ID:        payload.ClientID,
		SessionID: sessionId,
		conn:      conn,
		reader:    reader,
	}, nil
}

func (self *Client) Close() {
	self.conn.Close()
}

func (self *Client) Read() (*Message, error) {
	return ReadMessage(self.reader)
}

func (self *Client) Write(m *Message) error {
	_, err := self.conn.Write(m.Serialize())
	return err
}
