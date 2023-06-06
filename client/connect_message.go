package client

import (
	"vayeate/serialize"
)

type ConnectMessagePayload struct {
	ClientID string
	Username string
	Password string
}

func (self *Message) GetConnectPayload() *ConnectMessagePayload {
	v := ConnectMessagePayload{}
	serialize.Unmarshall(self.Payload, &v)
	return &v
}

type ConnectAckMessagePayload struct {
	SessionID string
}

func NewConnectAckMessage(sessionId string) *Message {
	v := ConnectAckMessagePayload{sessionId}

	return &Message{
		Code:    CONNECT_ACK,
		Payload: serialize.Marshall(v),
	}
}
