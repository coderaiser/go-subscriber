package store

import "sync"

// Facts holds data that persists across subscriptions for one (msisdn, serviceID).
type Facts struct {
	CooloffUntil    int64 // unix seconds; 0 = no cooloff
	FallbackAttempt int   // low_balance ladder counter; reset on success
	TrialUsed       bool  // one trial per (msisdn, serviceID), permanent
}

// FactsStore is owned by the engine only.
// Rows are never deleted -- cooloff and counter history must survive re-subscription.
type FactsStore struct {
	mu    sync.RWMutex
	facts map[string]Facts
}

func NewFactsStore() *FactsStore {
	return &FactsStore{facts: make(map[string]Facts)}
}

func (f *FactsStore) Get(id string) Facts {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.facts[id]
}

func (f *FactsStore) Set(id string, facts Facts) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.facts[id] = facts
}
