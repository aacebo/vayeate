package client

import "time"

func NewPingAckMessage() *Message {
	return &Message{
		Code:    PING_ACK,
		SentAt:  time.Now().UnixMilli(),
		Payload: []byte{},
	}
}
