package sync

import "sync"

type SyncSet[K comparable, V any] struct {
	mu       sync.RWMutex
	i        int
	mapping  map[K]int
	iterable []V
}

func NewSyncSet[K comparable, V any]() SyncSet[K, V] {
	return SyncSet[K, V]{
		sync.RWMutex{},
		0,
		map[K]int{},
		[]V{},
	}
}

func (self *SyncSet[K, V]) Has(key K) bool {
	self.mu.RLock()
	_, ok := self.mapping[key]
	self.mu.RUnlock()
	return ok
}

func (self *SyncSet[K, V]) Len() int {
	self.mu.RLock()
	l := len(self.iterable)
	self.mu.RUnlock()
	return l
}

func (self *SyncSet[K, V]) Add(key K, value V) {
	if self.Has(key) {
		return
	}

	self.mu.Lock()
	self.iterable = append(self.iterable, value)
	self.mapping[key] = len(self.iterable) - 1
	self.mu.Unlock()
}

func (self *SyncSet[K, V]) Del(key K) {
	if !self.Has(key) {
		return
	}

	self.mu.Lock()
	i := self.mapping[key]
	delete(self.mapping, key)
	self.iterable = append(self.iterable[:i], self.iterable[i+1:]...)
	self.mu.Unlock()
}

func (self *SyncSet[K, V]) Next() V {
	self.mu.RLock()
	value := self.iterable[self.i]
	self.i++

	if self.i > len(self.iterable)-1 {
		self.i = 0
	}

	self.mu.RUnlock()
	return value
}
