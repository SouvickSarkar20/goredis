package store

import "fmt"

type Set struct {
	data map[string]struct{}
}

func NewSet() *Set {
	return &Set{
		data: make(map[string]struct{}),
	}
}

func (s *Set) Add(val string) {
	s.data[val] = struct{}{}
}

func (s *Set) Remove(val string) bool {
	_, exists := s.data[val]
	if exists {
		delete(s.data, val)
	}
	return exists
}

func (s *Set) Members() []string {
	result := make([]string, 0, len(s.data))
	for k := range s.data {
		result = append(result, k)
	}
	return result
}

func (s *Set) IsMember(val string) bool {
	_, exists := s.data[val]
	return exists
}

// STORE METHODS

func (s *Store) SAdd(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.data[key]

	var set *Set

	if !exists {
		set = NewSet()
		s.data[key] = Item{Value: set}
	} else {
		var ok bool
		set, ok = item.Value.(*Set)
		if !ok {
			return fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}

	set.Add(value)
	return nil
}

func (s *Store) SMembers(key string) ([]string, error) {
	s.mu.RLock()
	item, exists := s.data[key]
	s.mu.RUnlock()

	if !exists {
		return []string{}, nil
	}

	set, ok := item.Value.(*Set)
	if !ok {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return set.Members(), nil
}

func (s *Store) SIsMember(key, value string) (bool, error) {
	s.mu.RLock()
	item, exists := s.data[key]
	s.mu.RUnlock()

	if !exists {
		return false, nil
	}

	set, ok := item.Value.(*Set)
	if !ok {
		return false, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return set.IsMember(value), nil
}

func (s *Store) SRem(key string, value string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.data[key]
	if !exists {
		return false, nil
	}

	set, ok := item.Value.(*Set)
	if !ok {
		return false, fmt.Errorf("WRONGTYPE")
	}

	removed := set.Remove(value)

	if len(set.data) == 0 {
		delete(s.data, key)
	}

	return removed, nil
}
