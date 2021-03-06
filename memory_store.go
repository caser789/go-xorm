// MemoryStore implements CacheStore provide local machine
package xorm

import (
	"sync"

	"github.com/caser789/go-xorm/core"
)

var _ core.CacheStore = NewMemoryStore()

// memory store
type MemoryStore struct {
	store map[interface{}]interface{}
	mutex sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{store: make(map[interface{}]interface{})}
}

func (s *MemoryStore) Put(key string, value interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store[key] = value
	return nil
}

func (s *MemoryStore) Get(key string) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if v, ok := s.store[key]; ok {
		return v, nil
	}

	return nil, ErrNotExist
}

func (s *MemoryStore) Del(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.store, key)
	return nil
}
