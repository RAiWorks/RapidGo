package session

import (
	"sync"
	"time"
)

type memEntry struct {
	Data      map[string]interface{}
	ExpiresAt time.Time
}

// MemoryStore stores sessions in memory. For development and testing only.
type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]memEntry
}

// NewMemoryStore creates a new in-memory session store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{sessions: make(map[string]memEntry)}
}

func (s *MemoryStore) Read(id string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.sessions[id]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return make(map[string]interface{}), nil
	}
	return entry.Data, nil
}

func (s *MemoryStore) Write(id string, data map[string]interface{}, lifetime time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[id] = memEntry{Data: data, ExpiresAt: time.Now().Add(lifetime)}
	return nil
}

func (s *MemoryStore) Destroy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
	return nil
}

func (s *MemoryStore) GC(maxLifetime time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for id, entry := range s.sessions {
		if now.After(entry.ExpiresAt) {
			delete(s.sessions, id)
		}
	}
	return nil
}
