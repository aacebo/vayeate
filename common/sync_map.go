package common

import "sync"

type SyncMap[K comparable, V interface{}] struct {
	mu sync.RWMutex
	content map[K]V
}

func NewSyncMap[K comparable, V interface{}]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		sync.RWMutex{},
		map[K]V{},
	}
}

func (self *SyncMap[K, V]) Get(key K) V {
	self.mu.RLock()
	v := self.content[key]
	self.mu.RUnlock()
	return v
}

func (self *SyncMap[K, V]) Set(key K, value V) {
	self.mu.Lock()
	self.content[key] = value
	self.mu.Unlock()
}

func (self *SyncMap[K, V]) Del(key K) {
	self.mu.Lock()
	delete(self.content, key)
	self.mu.Unlock()
}

func (self *SyncMap[K, V]) Iterate(callback func(k K, v V)) {
	self.mu.Lock()

	for k, v := range self.content {
		callback(k, v)
	}

	self.mu.Unlock()
}
