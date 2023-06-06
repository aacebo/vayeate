package client

import (
	"time"
	"vayeate/serialize"
)

type ErrorMessagePayload struct {
	Reason string
}

func NewErrorMessage(reason string) *Message {
	v := ErrorMessagePayload{reason}

	return &Message{
		Code:    ERROR,
		SentAt:  time.Now().UnixMilli(),
		Payload: serialize.Marshall(v),
	}
}

func (self *Message) GetErrorPayload() *ErrorMessagePayload {
	v := ErrorMessagePayload{}
	serialize.Unmarshall(self.Payload, &v)
	return &v
}
