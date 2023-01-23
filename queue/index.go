package queue

import (
	"fmt"

	"vayeate/common"
	"vayeate/logger"

	"github.com/google/uuid"
)

type Queue struct {
	ID   string
	Name string

	log       *logger.Logger
	consumers *common.SyncSet[string]
	buffer    chan []byte
}

func New(name string) *Queue {
	id := uuid.NewString()
	log := logger.New(fmt.Sprintf("queue:%s", id))
	consumers := common.NewSyncSet[string]()

	return &Queue{
		id,
		name,
		log,
		consumers,
		make(chan []byte),
	}
}

func (self *Queue) Push(payload []byte) {
	self.buffer <- payload
}

func (self *Queue) Consume(id string) {
	self.consumers.Add(id)
}

func (self *Queue) Start(onConsume func(id string, payload []byte)) {
	for {
		payload := <-self.buffer
		id := self.consumers.CyclicNext()
		onConsume(id, payload)
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
