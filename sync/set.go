package sync

import "sync"

type Set[K comparable, V any] struct {
	mu     sync.RWMutex
	i      int
	keys   []K
	values map[K]V
}

func NewSet[K comparable, V any]() Set[K, V] {
	return Set[K, V]{
		sync.RWMutex{},
		0,
		[]K{},
		map[K]V{},
	}
}

func (self *Set[K, V]) Has(key K) bool {
	self.mu.RLock()
	_, ok := self.values[key]
	self.mu.RUnlock()
	return ok
}

func (self *Set[K, V]) Len() int {
	self.mu.RLock()
	l := len(self.keys)
	self.mu.RUnlock()
	return l
}

func (self *Set[K, V]) Add(key K, value V) {
	if self.Has(key) {
		return
	}

	self.mu.Lock()
	self.keys = append(self.keys, key)
	self.values[key] = value
	self.mu.Unlock()
}

func (self *Set[K, V]) Del(key K) {
	if !self.Has(key) {
		return
	}

	self.mu.Lock()

	for i := 0; i < len(self.keys); i++ {
		if key == self.keys[i] {
			self.keys = append(self.keys[:i], self.keys[i+1:]...)
			delete(self.values, key)
			break
		}
	}

	self.mu.Unlock()
}

func (self *Set[K, V]) Next() V {
	self.mu.Lock()

	if self.i > len(self.keys)-1 {
		self.i = 0
	}

	key := self.keys[self.i]
	self.i++

	if self.i > len(self.keys)-1 {
		self.i = 0
	}

	self.mu.Unlock()
	return self.values[key]
}

func (self *Set[K, V]) ForEach(callback func(i int, v V)) {
	self.mu.RLock()

	for i, key := range self.keys {
		callback(i, self.values[key])
	}

	self.mu.RUnlock()
}
