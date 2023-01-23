package common

import "sync"

type SyncSet[V comparable] struct {
	mu      sync.RWMutex
	idx     int
	content []V
}

func NewSyncSet[V comparable]() *SyncSet[V] {
	return &SyncSet[V]{
		idx:     0,
		content: []V{},
	}
}

func (self *SyncSet[V]) Reset() {
	self.mu.Lock()
	self.idx = 0
	self.mu.Unlock()
}

func (self *SyncSet[V]) Next() V {
	self.mu.Lock()
	l := len(self.content)

	if self.idx == l-1 {
		panic("idx out of range")
	}

	self.idx++
	v := self.content[self.idx]

	self.mu.Unlock()
	return v
}

func (self *SyncSet[V]) CyclicNext() V {
	self.mu.Lock()
	l := len(self.content)

	if self.idx == l-1 {
		self.idx = 0
	} else {
		self.idx++
	}

	v := self.content[self.idx]

	self.mu.Unlock()
	return v
}

func (self *SyncSet[V]) Current() V {
	self.mu.RLock()
	v := self.content[self.idx]
	self.mu.RUnlock()
	return v
}

func (self *SyncSet[V]) Add(v V) {
	self.mu.Lock()
	l := len(self.content)

	if self.idx == l-1 || l == 0 {
		self.content = append(self.content, v)
	} else {
		self.content = append(self.content[:self.idx+1], self.content[self.idx:]...)
		self.content[self.idx] = v
	}

	self.mu.Unlock()
}

func (self *SyncSet[V]) Len() int {
	self.mu.RLock()
	l := len(self.content)
	self.mu.RUnlock()
	return l
}
