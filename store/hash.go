package store

import "fmt"

type Hash struct {
	data map[string]string
}

func NewHash() *Hash {
	return &Hash{
		data: make(map[string]string),
	}
}

func (h *Hash) HSet(field, value string) {
	h.data[field] = value
}

func (h *Hash) HGet(field string) (string, bool) {
	val, ok := h.data[field]
	return val, ok
}

// STORE METHODS

func (s *Store) HSet(key, field, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.data[key]

	var hash *Hash

	if !exists {
		hash = NewHash()
		s.data[key] = Item{Value: hash}
	} else {
		var ok bool
		hash, ok = item.Value.(*Hash)
		if !ok {
			return fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}

	hash.HSet(field, value)
	return nil
}

func (s *Store) HGet(key, field string) (string, bool, error) {
	s.mu.RLock()
	item, exists := s.data[key]
	s.mu.RUnlock()

	if !exists {
		return "", false, nil
	}

	hash, ok := item.Value.(*Hash)
	if !ok {
		return "", false, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	val, ok := hash.HGet(field)
	return val, ok, nil
}
