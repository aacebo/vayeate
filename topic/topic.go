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
	Messages    sync.SyncQueue[[]byte]               `json:"-"`
	Subscribers sync.SyncSet[string, *client.Client] `json:"-"`

	log *logger.Logger
}

func New(name string) *Topic {
	self := Topic{
		Name:        name,
		CreatedAt:   time.Now(),
		Messages:    sync.NewSyncQueue[[]byte](),
		Subscribers: sync.NewSyncSet[string, *client.Client](),
		log:         logger.New(fmt.Sprintf("vayeate:topic:%s", name)),
	}

	go func() {
		for {
			if self.Subscribers.Len() == 0 || self.Messages.Len() == 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			payload := self.Messages.Pop()
			c := self.Subscribers.Next()
			err := c.Write(client.NewConsumeMessage(self.Name, payload))

			if err != nil {
				continue
			}

			// for m := range c.Messages {
			// 	if m.Code == client.CONSUME_ACK {
			// 		p := m.GetConsumeAckPayload()

			// 		if p.Topic == self.Name {
			// 			self.Queue.Pop()
			// 			break
			// 		}
			// 	}
			// }
		}
	}()

	return &self
}
