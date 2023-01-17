package queue

import (
	"fmt"
	"vayeate/logger"

	"github.com/google/uuid"
)

type Queue struct {
	ID       string
	Name     string
	Messages []*Message

	log *logger.Logger
}

func New(name string) *Queue {
	id := uuid.NewString()
	log := logger.New(fmt.Sprintf("queue:%s", id))
	self := Queue{id, name, []*Message{}, log}
	log.Infoln("created")
	return &self
}

func (self *Queue) Push(payload []byte) *Message {
	message := NewMessage(payload)
	self.Messages = append(self.Messages, message)
	return message
}

func (self *Queue) Pop() *Message {
	message := self.Messages[0]
	self.Messages = self.Messages[1:]
	return message
}
