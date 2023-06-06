package client

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
)

var timeout = 60 * time.Second

type Client struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id"`
	Address       string    `json:"address"`
	ConnectedAt   time.Time `json:"connected_at"`
	LastMessageAt time.Time `json:"last_message_at"`

	open      bool
	conn      net.Conn
	reader    *bufio.Reader
	pingTimer *time.Timer
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

	self := Client{
		ID:            payload.ClientID,
		SessionID:     sessionId,
		Address:       conn.RemoteAddr().String(),
		ConnectedAt:   time.Now(),
		LastMessageAt: time.Now(),
		open:          true,
		conn:          conn,
		reader:        reader,
	}

	self.pingTimer = time.AfterFunc(timeout, self.onConnectionTimeout)
	return &self, nil
}

func (self *Client) Close() {
	self.pingTimer.Stop()
	self.conn.Close()
	self.open = false
}

func (self *Client) Read() (*Message, error) {
	if !self.open {
		return nil, io.EOF
	}

	m, err := ReadMessage(self.reader)

	if err != nil {
		return nil, err
	}

	self.LastMessageAt = time.Now()
	self.pingTimer.Reset(timeout)
	return m, nil
}

func (self *Client) Write(m *Message) error {
	_, err := self.conn.Write(m.Serialize())
	return err
}

func (self *Client) onConnectionTimeout() {
	self.Write(NewErrorMessage("connection closed due to inactivity"))
	self.Close()
}
