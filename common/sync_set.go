package common

import "sync"

type SyncSet[V comparable] struct {
	mu       sync.RWMutex
	idx      int
	content  map[V]int
	iterable []V
}

func NewSyncSet[V comparable]() *SyncSet[V] {
	return &SyncSet[V]{
		idx:      -1,
		content:  map[V]int{},
		iterable: []V{},
	}
}

func (self *SyncSet[V]) Reset() {
	self.mu.Lock()
	self.idx = 0
	self.mu.Unlock()
}

func (self *SyncSet[V]) Next() V {
	self.mu.Lock()
	l := len(self.iterable)

	if self.idx == l-1 {
		panic("idx out of range")
	}

	self.idx++
	v := self.iterable[self.idx]
	self.mu.Unlock()
	return v
}

func (self *SyncSet[V]) CyclicNext() V {
	self.mu.Lock()
	l := len(self.iterable)

	if self.idx == l-1 {
		self.idx = 0
	} else {
		self.idx++
	}

	v := self.iterable[self.idx]
	self.mu.Unlock()
	return v
}

func (self *SyncSet[V]) Current() V {
	self.mu.RLock()
	v := self.iterable[self.idx]
	self.mu.RUnlock()
	return v
}

func (self *SyncSet[V]) Add(v V) {
	self.mu.Lock()

	if _, ok := self.content[v]; ok {
		return
	}

	if self.idx == -1 {
		self.idx = 0
	}

	l := len(self.iterable)

	if self.idx == l-1 || l == 0 {
		self.iterable = append(self.iterable, v)
	} else {
		self.iterable = append(self.iterable[:self.idx+1], self.iterable[self.idx:]...)
		self.iterable[self.idx] = v
	}

	self.content[v] = self.idx
	self.mu.Unlock()
}

func (self *SyncSet[V]) Del(v V) {
	self.mu.Lock()

	if self.content[v] > -1 {
		delete(self.content, v)
		self.iterable = append(
			self.iterable[:self.content[v]],
			self.iterable[self.content[v]+1:]...,
		)
	}

	self.mu.Unlock()
}

func (self *SyncSet[V]) Len() int {
	self.mu.RLock()
	l := len(self.iterable)
	self.mu.RUnlock()
	return l
}
