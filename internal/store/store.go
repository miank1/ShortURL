package store

import "sync"

type Store interface {
	Save(code, url string)
	Get(code string) (string, bool)
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

func (s *MemoryStore) Save(code, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[code] = url
}

func (s *MemoryStore) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.data[code]
	return url, ok
}
