package adapters

import "sync"

// Memory is an in-memory Adapter.
// Get returns nil,nil for unknown ids — not an error.
type Memory[S any] struct {
	mu     sync.RWMutex
	states map[string]S
}

func NewMemory[S any]() *Memory[S] {
	return &Memory[S]{states: make(map[string]S)}
}

func (m *Memory[S]) Get(id string) (*S, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.states[id]
	if !ok {
		return nil, nil
	}
	return &s, nil
}

func (m *Memory[S]) Set(id string, state S) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.states[id] = state
	return nil
}
