package client

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
	"vayeate/sync"

	"github.com/google/uuid"
)

var connectionTimeout = 60 * time.Second
var readTimeout = 60 * time.Second
var writeTimeout = 5 * time.Second

type Client struct {
	ID            string                       `json:"id"`
	SessionID     string                       `json:"session_id"`
	Address       string                       `json:"address"`
	LatencyMS     int64                        `json:"latency_ms"`
	ConnectedAt   time.Time                    `json:"connected_at"`
	LastMessageAt time.Time                    `json:"last_message_at"`
	Topics        sync.SyncSet[string, string] `json:"-"`

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

	self := Client{
		ID:            payload.ClientID,
		SessionID:     uuid.NewString(),
		Address:       conn.RemoteAddr().String(),
		LatencyMS:     time.Now().UnixMilli() - m.SentAt,
		ConnectedAt:   time.Now(),
		LastMessageAt: time.Now(),
		Topics:        sync.NewSyncSet[string, string](),
		open:          true,
		conn:          conn,
		reader:        reader,
	}

	self.pingTimer = time.AfterFunc(connectionTimeout, self.onConnectionTimeout)
	return &self, nil
}

func (self *Client) Close() {
	self.pingTimer.Stop()
	self.conn.Close()
	self.open = false
}

func (self *Client) Read() (*Message, error) {
	self.conn.SetReadDeadline(time.Now().Add(readTimeout))

	if !self.open {
		return nil, io.EOF
	}

	m, err := ReadMessage(self.reader)

	if err != nil {
		return nil, err
	}

	self.LatencyMS = time.Now().UnixMilli() - m.SentAt
	self.LastMessageAt = time.Now()
	self.pingTimer.Reset(connectionTimeout)

	return m, nil
}

func (self *Client) Write(m *Message) error {
	self.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	_, err := self.conn.Write(m.Serialize())
	return err
}

func (self *Client) onConnectionTimeout() {
	self.Write(NewErrorMessage("connection closed due to inactivity"))
	self.Close()
}
