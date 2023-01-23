package queue

import (
	"fmt"

	"vayeate/logger"

	"github.com/google/uuid"
)

type Queue struct {
	ID   string
	Name string

	log  *logger.Logger
	push chan []byte
}

func New(name string) *Queue {
	id := uuid.NewString()
	log := logger.New(fmt.Sprintf("queue:%s", id))

	return &Queue{
		id,
		name,
		log,
		make(chan []byte),
	}
}

func (self *Queue) Push(payload []byte) {
	self.push <- payload
}

func (self *Queue) Pop(id string) {

}

func (self *Queue) Start() {
	for {
		payload := <-self.push
		self.log.Infof("received message %s", string(payload))
	}
}

// func (self *Queue) Push(payload []byte) *Message {
// 	message := NewMessage(payload)
// 	self.messages = append(self.messages, message)
// 	return message
// }

// func (self *Queue) Pop() *Message {
// 	message := self.messages[0]
// 	self.messages = self.messages[1:]
// 	return message
// }
