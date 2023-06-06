package client

import (
	"bufio"
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id"`
	Address       string    `json:"address"`
	ConnectedAt   time.Time `json:"connected_at"`
	LastMessageAt time.Time `json:"last_message_at"`

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
		return nil, errors.New("unauthorized: first message must be of type `CONNECT`")
	}

	payload := m.GetConnectPayload()

	if payload.Username != username || payload.Password != password {
		return nil, errors.New("unauthorized: invalid username/password")
	}

	sessionId := uuid.NewString()

	return &Client{
		ID:            payload.ClientID,
		SessionID:     sessionId,
		Address:       conn.RemoteAddr().String(),
		ConnectedAt:   time.Now(),
		LastMessageAt: time.Now(),
		conn:          conn,
		reader:        reader,
	}, nil
}

func (self *Client) Close() {
	self.conn.Close()
}

func (self *Client) Read() (*Message, error) {
	m, err := ReadMessage(self.reader)

	if err != nil {
		return nil, err
	}

	self.LastMessageAt = time.Now()
	return m, nil
}

func (self *Client) Write(m *Message) error {
	_, err := self.conn.Write(m.Serialize())
	return err
}
