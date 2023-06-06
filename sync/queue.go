package sync

import "sync"

type SyncQueue[V any] struct {
	mu      sync.RWMutex
	content []V
}

func NewSyncQueue[V any]() SyncQueue[V] {
	return SyncQueue[V]{
		sync.RWMutex{},
		[]V{},
	}
}

func (self *SyncQueue[V]) Push(value V) {
	self.mu.Lock()
	self.content = append(self.content, value)
	self.mu.Unlock()
}

func (self *SyncQueue[V]) Pop() V {
	self.mu.Lock()
	value := self.content[0]
	self.content = self.content[1:]
	self.mu.Unlock()
	return value
}

func (self *SyncQueue[V]) Len() int {
	self.mu.RLock()
	l := len(self.content)
	self.mu.RUnlock()
	return l
}
