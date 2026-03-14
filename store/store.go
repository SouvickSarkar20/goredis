package store

import (
	"sync"
	"time"
)

type Item struct {
	Value  string
	Expiry time.Time
}

type Store struct {
	data map[string]Item
	mu   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Item),
	}
}

func (s *Store) Set(key, value string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiry time.Time
	if duration > 0 {
		expiry = time.Now().Add(duration)
	}

	s.data[key] = Item{
		Value:  value,
		Expiry: expiry,
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	item, exists := s.data[key]
	s.mu.RUnlock()

	if !exists {
		return "", false
	}

	if !item.Expiry.IsZero() && time.Now().After(item.Expiry) {
		s.Delete(key)
		return "", false
	}

	return item.Value, true
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
