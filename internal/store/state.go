package store

import (
	"sync"

)

// StateStore is the statemachine adapter.
// It stores current state string per composite key (msisdn:serviceID).
type StateStore struct {
	mu     sync.RWMutex
	states map[string]string
}

func NewStateStore() *StateStore {
	return &StateStore{states: make(map[string]string)}
}

func (s *StateStore) Get(id string) (*string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.states[id]
	if !ok {
		return nil, nil
	}
	return &v, nil
}

func (s *StateStore) Set(id string, state string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[id] = state
	return nil
}

func (s *StateStore) All() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]string, len(s.states))
	for k, v := range s.states {
		result[k] = v
	}
	return result
}
