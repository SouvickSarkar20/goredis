package store

import (
	"sync"
)

// Store is our concurrent-safe in-memory database.
type Store struct {
	// A standard Go map holding our Key-Value pairs.
	// We store the keys as strings, and the values as strings.
	data map[string]string

	// RWMutex allows MULTIPLE readers at the same time (GET),
	// but only ONE writer at a time (SET).
	// This is much faster than a standard Mutex if we have many GET requests.
	mu sync.RWMutex
}

// NewStore creates and initializes a new instance of our database.
func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

// Set stores a value in the map, acquiring a write lock to ensure safety.
func (s *Store) Set(key, value string) {
	// 1. Lock the Mutex for writing. 
	// If any other goroutine is currently reading or writing, this blocks until they finish.
	s.mu.Lock()

	// 2. Ensure the lock is ALWAYS released when this function finishes,
	// even if the code below were to panic.
	defer s.mu.Unlock()

	// 3. Perform the actual write to the map.
	s.data[key] = value
}

// Get retrieves a value from the map, acquiring a read lock.
func (s *Store) Get(key string) (string, bool) {
	// 1. Lock the Mutex for reading.
	// Multiple goroutines can hold a RLock at the same time.
	// But if someone is currently holding a write Lock(), this blocks until they finish writing.
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 2. Read from the map.
	val, exists := s.data[key]
	return val, exists
}

// Delete removes a key from the map, acquiring a write lock.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
