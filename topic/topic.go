package topic

import (
	"fmt"
	"time"
	"vayeate/client"
	"vayeate/logger"
	"vayeate/sync"
)

type Topic struct {
	Name        string                               `json:"name"`
	CreatedAt   time.Time                            `json:"created_at"`
	Queue       sync.SyncQueue[[]byte]               `json:"-"`
	Subscribers sync.SyncSet[string, *client.Client] `json:"-"`

	log *logger.Logger
}

func New(name string) *Topic {
	self := Topic{
		Name:        name,
		CreatedAt:   time.Now(),
		Queue:       sync.NewSyncQueue[[]byte](),
		Subscribers: sync.NewSyncSet[string, *client.Client](),
		log:         logger.New(fmt.Sprintf("vayeate:topic:%s", name)),
	}

	go func() {
		for {
			if self.Subscribers.Len() == 0 || self.Queue.Len() == 0 {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			payload := self.Queue.Top()
			c := self.Subscribers.Next()
			err := c.Write(client.NewConsumeMessage(self.Name, payload))

			if err != nil {
				continue
			}

			for m := range c.Messages {
				if m.Code == client.CONSUME_ACK {
					p := m.GetConsumeAckPayload()

					if p.Topic == self.Name {
						self.Queue.Pop()
						break
					}
				}
			}
		}
	}()

	return &self
}
