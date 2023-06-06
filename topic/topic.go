package topic

import (
	"time"
	"vayeate/sync"
)

type Topic struct {
	Name      string                 `json:"name"`
	CreatedAt time.Time              `json:"created_at"`
	Queue     sync.SyncQueue[[]byte] `json:"-"`
}

func New(name string) *Topic {
	return &Topic{
		Name:      name,
		CreatedAt: time.Now(),
		Queue:     sync.NewSyncQueue[[]byte](),
	}
}
