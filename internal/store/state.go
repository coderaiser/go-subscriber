package store

import (
	"sync"

	"github.com/coderaiser/go-subscriber/internal/engine"
)

// StateStore is the statemachine adapter.
// It stores current engine.State per composite key (msisdn:serviceID).
type StateStore struct {
	mu     sync.RWMutex
	states map[string]engine.State
}

func NewStateStore() *StateStore {
	return &StateStore{states: make(map[string]engine.State)}
}

func (s *StateStore) Get(id string) (*engine.State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.states[id]
	if !ok {
		return nil, nil
	}
	return &v, nil
}

func (s *StateStore) Set(id string, state engine.State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[id] = state
	return nil
}
