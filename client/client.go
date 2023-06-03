package client

import (
	"bufio"
	"net"
)

type Client struct {
	ID string `json:"id"`

	conn   net.Conn
	reader *bufio.Reader
}

func FromConnection(username string, password string, conn net.Conn) *Client {
	reader := bufio.NewReader(conn)
	self := Client{"", conn, reader}
	return &self
}

func (self *Client) Close() {
	self.conn.Close()
}

func (self *Client) Read() (*Message, error) {
	return ReadMessage(self.reader)
}
