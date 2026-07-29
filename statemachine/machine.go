package statemachine

import "fmt"

type Machine[S comparable, E comparable] struct {
	transitions map[S]map[E]S
	hooks       map[string]Handler[S, E]
	initial     S
	adapter     Adapter[S]
	strict      bool
}

func New[S, E comparable](
	source TransitionSource,
	parseState func(string) (S, error),
	parseEvent func(string) (E, error),
	adapter Adapter[S],
	strict bool,
) (*Machine[S, E], error) {
	defs, err := source.Load()
	if err != nil {
		return nil, err
	}
	m := &Machine[S, E]{
		transitions: make(map[S]map[E]S),
		hooks:       make(map[string]Handler[S, E]),
		adapter:     adapter,
		strict:      strict,
	}
	for _, d := range defs {
		from, err := parseState(d.From)
		if err != nil {
			return nil, fmt.Errorf("invalid state %q: %w", d.From, err)
		}
		event, err := parseEvent(d.Event)
		if err != nil {
			return nil, fmt.Errorf("invalid event %q: %w", d.Event, err)
		}
		to, err := parseState(d.To)
		if err != nil {
			return nil, fmt.Errorf("invalid state %q: %w", d.To, err)
		}
		if m.transitions[from] == nil {
			m.transitions[from] = make(map[E]S)
		}
		m.transitions[from][event] = to
	}
	return m, nil
}

// WithInitial sets the state used when adapter.Get returns nil.
func (m *Machine[S, E]) WithInitial(s S) *Machine[S, E] {
	m.initial = s
	return m
}

func (m *Machine[S, E]) Hook(from S, event E, h Handler[S, E]) {
	key := fmt.Sprintf("%v:%v", from, event)
	m.hooks[key] = h
}

// Apply order:
//  1. adapter.Get → nil,nil uses m.initial; nil,err returns error
//  2. look up transition → panic (strict) or error
//  3. adapter.Set → store BEFORE hook
//  4. run hook if registered
//  5. return next, hookErr
func (m *Machine[S, E]) Apply(id string, event E, payload any) (S, error) {
	ptr, err := m.adapter.Get(id)
	if err != nil {
		var zero S
		return zero, fmt.Errorf("adapter.Get: %w", err)
	}
	current := m.initial
	if ptr != nil {
		current = *ptr
	}
	next, ok := m.transitions[current][event]
	if !ok {
		err := fmt.Errorf("invalid transition: from %v event %v", current, event)
		if m.strict {
			panic(err.Error())
		}
		return current, err
	}
	if err := m.adapter.Set(id, next); err != nil {
		return current, fmt.Errorf("adapter.Set: %w", err)
	}
	key := fmt.Sprintf("%v:%v", current, event)
	if h, ok := m.hooks[key]; ok {
		if err := h(Context[S, E]{
			From: current, Event: event,
			Payload: payload, Adapter: m.adapter,
		}); err != nil {
			return next, err
		}
	}
	return next, nil
}

func (m *Machine[S, E]) Validate() error {
	for from, events := range m.transitions {
		if len(events) == 0 {
			return fmt.Errorf("state %v has no transitions out", from)
		}
	}
	return nil
}
