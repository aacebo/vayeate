package common

import (
	"regexp"
	"sync"
)

type SyncMap[V interface{}] struct {
	mu      sync.RWMutex
	content map[string]V
}

func NewSyncMap[V interface{}]() *SyncMap[V] {
	return &SyncMap[V]{
		sync.RWMutex{},
		map[string]V{},
	}
}

func (self *SyncMap[V]) Get(key string) V {
	self.mu.RLock()
	v := self.content[key]
	self.mu.RUnlock()
	return v
}

func (self *SyncMap[V]) Set(key string, value V) {
	self.mu.Lock()
	self.content[key] = value
	self.mu.Unlock()
}

func (self *SyncMap[V]) Del(key string) {
	self.mu.Lock()
	delete(self.content, key)
	self.mu.Unlock()
}

func (self *SyncMap[V]) Len() int {
	self.mu.Lock()
	l := len(self.content)
	self.mu.Unlock()
	return l
}

func (self *SyncMap[V]) Iterate(callback func(k string, v V)) {
	self.mu.Lock()

	for k, v := range self.content {
		callback(k, v)
	}

	self.mu.Unlock()
}

func (self *SyncMap[V]) GetMatching(pattern string) []V {
	self.mu.RLock()
	slice := []V{}

	for k, v := range self.content {
		match, _ := regexp.MatchString(pattern, k)

		if match == true {
			slice = append(slice, v)
		}
	}

	self.mu.RUnlock()
	return slice
}
