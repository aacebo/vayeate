package client

import (
	"time"
	"vayeate/serialize"
)

type ConsumeMessagePayload struct {
	Topic   string
	Payload []byte
}

func NewConsumeMessage(topic string, payload []byte) *Message {
	v := ConsumeMessagePayload{topic, payload}

	return &Message{
		Code:    CONSUME,
		SentAt:  time.Now().UnixMilli(),
		Payload: serialize.Marshall(v),
	}
}

type ConsumeAckMessagePayload struct {
	Topic string
}

func (self *Message) GetConsumeAckPayload() *ConsumeAckMessagePayload {
	v := ConsumeAckMessagePayload{}
	serialize.Unmarshall(self.Payload, &v)
	return &v
}
