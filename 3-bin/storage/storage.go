package storage

import "sync"

// InMemory stores created bin IDs locally
type InMemory struct {
	mu  sync.RWMutex
	ids map[string]struct{}
}

// NewInMemory creates a new in-memory storage
func NewInMemory() *InMemory {
	return &InMemory{ids: make(map[string]struct{})}
}

// Add records an id
func (s *InMemory) Add(id string) {
	s.mu.Lock()
	s.ids[id] = struct{}{}
	s.mu.Unlock()
}

// Has checks if id exists
func (s *InMemory) Has(id string) bool {
	s.mu.RLock()
	_, ok := s.ids[id]
	s.mu.RUnlock()
	return ok
}

// All returns all ids (copy)
func (s *InMemory) All() []string {
	s.mu.RLock()
	out := make([]string, 0, len(s.ids))
	for id := range s.ids {
		out = append(out, id)
	}
	s.mu.RUnlock()
	return out
}
