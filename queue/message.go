package queue

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string
	Payload   []byte
	CreatedAt int64
}

func NewMessage(payload []byte) *Message {
	id := uuid.NewString()
	now := time.Now().Unix()
	self := Message{id, payload, now}
	return &self
}
