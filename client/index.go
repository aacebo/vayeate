package client

import (
	"net"

	"vayeate/frame"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(addr string) *Client {
	self := Client{addr, nil}
	return &self
}

func (self *Client) Connect() error {
	conn, err := net.Dial("tcp", self.addr)

	if err != nil {
		return err
	}

	self.conn = conn
	return nil
}

func (self *Client) Close() error {
	_, err := self.conn.Write(frame.NewClose().Encode())

	if err != nil {
		return err
	}

	return nil
}

func (self *Client) Assert(queue string) error {
	_, err := self.conn.Write(frame.NewAssert([]byte(queue)).Encode())

	if err != nil {
		return err
	}

	return nil
}

func (self *Client) Produce(subject string, body string) error {
	_, err := self.conn.Write(frame.NewProduce([]byte(subject), []byte(body)).Encode())

	if err != nil {
		return err
	}

	return nil
}

func (self *Client) Consume(subject string) error {
	_, err := self.conn.Write(frame.NewConsume([]byte(subject)).Encode())

	if err != nil {
		return err
	}

	return nil
}

func (self *Client) Ack(subject string) error {
	_, err := self.conn.Write(frame.NewAck([]byte(subject)).Encode())

	if err != nil {
		return err
	}

	return nil
}
