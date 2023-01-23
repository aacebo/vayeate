package queue

import (
	"fmt"

	"vayeate/logger"

	"github.com/google/uuid"
)

type Queue struct {
	ID   string
	Name string

	messages []*Message
	log      *logger.Logger
}

func New(name string) *Queue {
	id := uuid.NewString()
	log := logger.New(fmt.Sprintf("queue:%s", id))

	return &Queue{
		id,
		name,
		[]*Message{},
		log,
	}
}

func (self *Queue) Push(payload []byte) *Message {
	message := NewMessage(payload)
	self.messages = append(self.messages, message)
	return message
}

func (self *Queue) Pop() *Message {
	message := self.messages[0]
	self.messages = self.messages[1:]
	return message
}
