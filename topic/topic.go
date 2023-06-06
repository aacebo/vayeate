package topic

import (
	"time"
	"vayeate/client"
	"vayeate/sync"
)

type Topic struct {
	Name        string                               `json:"name"`
	CreatedAt   time.Time                            `json:"created_at"`
	Queue       sync.SyncQueue[[]byte]               `json:"-"`
	Subscribers sync.SyncSet[string, *client.Client] `json:"-"`
}

func New(name string) *Topic {
	return &Topic{
		Name:        name,
		CreatedAt:   time.Now(),
		Queue:       sync.NewSyncQueue[[]byte](),
		Subscribers: sync.NewSyncSet[string, *client.Client](),
	}
}
