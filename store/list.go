package store

import (
	"container/list"
	"fmt"
	"sync"
)

type List struct {
	data *list.List
	mu   sync.Mutex
}

func NewList() *List {
	return &List{
		data: list.New(),
	}
}

// method to push in the a particular list
func (l *List) LPush(val string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.data.PushFront(val)
}

func (l *List) LPop() (string, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	front := l.data.Front()
	if front == nil {
		return "", false
	}
	l.data.Remove(front)
	return front.Value.(string), true
}

func (s *Store) LPush(key, value string) error {
	s.mu.RLock()
	item, exists := s.data[key]
	s.mu.RUnlock()

	var linkedList *List

	if !exists {
		linkedList = NewList()
		s.Set(key, linkedList, 0)
	} else {
		var ok bool
		linkedList, ok = item.Value.(*List)
		if !ok {
			return fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}

	linkedList.LPush(value)

	return nil
}

func (s *Store) LPop(key string) (string, bool, error) {
	s.mu.RLock()
	item, exists := s.data[key]
	s.mu.RUnlock()

	if !exists {
		return "", false, nil
	}

	linkedList, ok := item.Value.(*List)
	if !ok {
		return "", false, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	val, ok := linkedList.LPop()

	return val, ok, nil
}
