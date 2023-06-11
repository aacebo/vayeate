package sync

import "sync"

type Queue[V any] struct {
	mu      sync.RWMutex
	content []V
}

func NewQueue[V any]() Queue[V] {
	return Queue[V]{
		sync.RWMutex{},
		[]V{},
	}
}

func (self *Queue[V]) Top() V {
	self.mu.RLock()
	value := self.content[0]
	self.mu.RUnlock()
	return value
}

func (self *Queue[V]) Push(value V) {
	self.mu.Lock()
	self.content = append(self.content, value)
	self.mu.Unlock()
}

func (self *Queue[V]) Pop() V {
	self.mu.Lock()
	value := self.content[0]
	self.content = self.content[1:]
	self.mu.Unlock()
	return value
}

func (self *Queue[V]) Len() int {
	self.mu.RLock()
	l := len(self.content)
	self.mu.RUnlock()
	return l
}
