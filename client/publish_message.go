package client

import (
	"vayeate/serialize"
)

type PublishMessagePayload struct {
	Topic   string
	Payload []byte
}

func (self *Message) GetPublishPayload() *PublishMessagePayload {
	v := PublishMessagePayload{}
	serialize.Unmarshall(self.Payload, &v)
	return &v
}

func NewPublishAckMessage() *Message {
	return &Message{
		Code:    PUBLISH_ACK,
		Payload: []byte{},
	}
}