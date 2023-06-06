package client

import (
	"time"
	"vayeate/serialize"
)

type SubscribeMessagePayload struct {
	Topic string
}

func (self *Message) GetSubscribePayload() *SubscribeMessagePayload {
	v := SubscribeMessagePayload{}
	serialize.Unmarshall(self.Payload, &v)
	return &v
}

func NewSubscribeAckMessage() *Message {
	return &Message{
		Code:    SUBSCRIBE_ACK,
		SentAt:  time.Now().UnixMilli(),
		Payload: []byte{},
	}
}
