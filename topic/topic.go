package topic

import (
	"fmt"
	"os"
	"time"
	"vayeate/client"
	"vayeate/logger"
	"vayeate/sync"
)

type Topic struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`

	log         *logger.Logger
	file        *os.File
	messages    sync.SyncQueue[[]byte]
	subscribers sync.SyncSet[string, *client.Client]
}

func New(name string) *Topic {
	now := time.Now()
	os.MkdirAll(fmt.Sprintf("topics/%s", name), 0700)
	file, _ := os.OpenFile(
		fmt.Sprintf("topics/%s/%d_%d_%d.log", name, now.Month(), now.Day(), now.Year()),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	self := Topic{
		Name:        name,
		CreatedAt:   time.Now(),
		log:         logger.New(fmt.Sprintf("vayeate:topic:%s", name)),
		file:        file,
		messages:    sync.NewSyncQueue[[]byte](),
		subscribers: sync.NewSyncSet[string, *client.Client](),
	}

	go func() {
		for {
			if self.subscribers.Len() == 0 || self.messages.Len() == 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			payload := self.messages.Pop()

			go func() {
				self.subscribers.ForEach(func(_ int, c *client.Client) {
					c.Write(client.NewConsumeMessage(self.Name, payload))
				})
			}()
		}
	}()

	return &self
}

func (self *Topic) MessagesLen() int {
	return self.messages.Len()
}

func (self *Topic) SubscribersLen() int {
	return self.subscribers.Len()
}

func (self *Topic) Push(payload []byte) {
	self.messages.Push(payload)

	go func() {
		self.file.Write(append(payload, '\x04'))
	}()
}

func (self *Topic) Subscribe(c *client.Client) {
	self.subscribers.Add(c.ID, c)
}

func (self *Topic) UnSubscribe(c *client.Client) {
	self.subscribers.Del(c.ID)
}
