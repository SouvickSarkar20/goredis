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
	s.mu.Lock()
	item, exists := s.data[key]
	s.mu.Unlock()

	if exists && isExpired(item) {
		delete(s.data, key)
		exists = false
	}

	var linkedList *List

	if !exists {
		linkedList = NewList()
		s.data[key] = Item{Value: linkedList}
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
	s.mu.Lock()
	item, exists := s.data[key]
	s.mu.Unlock()

	if !exists {
		return "", false, nil
	}

	if isExpired(item) {
		delete(s.data, key)
		return "", false, nil
	}

	linkedList, ok := item.Value.(*List)
	if !ok {
		return "", false, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	val, ok := linkedList.LPop()

	return val, ok, nil
}
