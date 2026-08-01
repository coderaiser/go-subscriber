package statemachine

import (
	"context"
	"fmt"
	"log/slog"
)

// TransitionDef is a single row in the transition table.
type TransitionDef struct {
	From  string
	Event string
	To    string
}

// TransitionSource provides the raw transition definitions.
type TransitionSource interface {
	Transitions() []TransitionDef
}

// MemorySource holds transition definitions in memory.
type MemorySource struct {
	Defs []TransitionDef
}

func (m *MemorySource) Transitions() []TransitionDef { return m.Defs }

// Adapter persists current state for an entity.
type Adapter[S comparable] interface {
	Get(id string) (*S, error)
	Set(id string, state S) error
}

// Context passed to hook functions.
type Context[S, E comparable] struct {
	ID    string
	From  S
	Event E
	To    S
	Data  any
}

// HookFunc is called after a successful transition.
type HookFunc[S, E comparable] func(ctx Context[S, E]) error

// Machine is a generic state machine.
type Machine[S, E comparable] struct {
	transitions map[S]map[E]S
	adapter     Adapter[S]
	initial     *S
	hooks       map[S]map[E][]HookFunc[S, E]
	log         *slog.Logger
}

// New creates a new Machine from a TransitionSource.
func New[S, E comparable](
	source TransitionSource,
	parseState func(string) (S, error),
	parseEvent func(string) (E, error),
	adapter Adapter[S],
) (*Machine[S, E], error) {
	m := &Machine[S, E]{
		adapter:     adapter,
		transitions: make(map[S]map[E]S),
		hooks:       make(map[S]map[E][]HookFunc[S, E]),
		log:         slog.New(noopHandler{}),
	}

	for _, d := range source.Transitions() {
		from, err := parseState(d.From)
		if err != nil {
			return nil, fmt.Errorf("parse from state %q: %w", d.From, err)
		}
		evt, err := parseEvent(d.Event)
		if err != nil {
			return nil, fmt.Errorf("parse event %q: %w", d.Event, err)
		}
		to, err := parseState(d.To)
		if err != nil {
			return nil, fmt.Errorf("parse to state %q: %w", d.To, err)
		}
		if m.transitions[from] == nil {
			m.transitions[from] = make(map[E]S)
		}
		m.transitions[from][evt] = to
	}

	return m, nil
}

// WithLogger sets the logger on the machine.
func (m *Machine[S, E]) WithLogger(log *slog.Logger) {
	m.log = log
}

// SetTransitions is used in tests to inject a raw transition map.
func (m *Machine[S, E]) SetTransitions(t map[S]map[E]S) {
	m.transitions = t
}

// WithInitial sets a fallback initial state for unknown IDs.
func (m *Machine[S, E]) WithInitial(s S) {
	m.initial = &s
}

// Hook registers a hook for a given (state, event) pair.
func (m *Machine[S, E]) Hook(state S, event E, fn HookFunc[S, E]) {
	if m.hooks[state] == nil {
		m.hooks[state] = make(map[E][]HookFunc[S, E])
	}
	m.hooks[state][event] = append(m.hooks[state][event], fn)
}

// Apply transitions from the current state for id using event.
func (m *Machine[S, E]) Apply(id string, event E, data any) (S, error) {
	var zero S

	ptr, err := m.adapter.Get(id)
	if err != nil {
		return zero, fmt.Errorf("adapter.Get: %w", err)
	}

	var current S
	if ptr == nil {
		if m.initial != nil {
			current = *m.initial
		}
		// When no state and no initial, use zero value
	} else {
		current = *ptr
	}

	evtMap, ok := m.transitions[current]
	if !ok {
		return m.handleInvalid(id, current, event)
	}

	next, ok := evtMap[event]
	if !ok {
		return m.handleInvalid(id, current, event)
	}

	if err := m.adapter.Set(id, next); err != nil {
		return zero, fmt.Errorf("adapter.Set: %w", err)
	}

	m.log.Info("transition", "id", id, "from", current, "event", event, "to", next)

	if hooks, ok := m.hooks[current][event]; ok {
		ctx := Context[S, E]{ID: id, From: current, Event: event, To: next, Data: data}
		for _, fn := range hooks {
			if err := fn(ctx); err != nil {
				return next, err
			}
		}
	}

	return next, nil
}

func (m *Machine[S, E]) handleInvalid(id string, current S, event E) (S, error) {
	var zero S
	err := fmt.Errorf("no transition from %v on event %v", current, event)
	return zero, err
}

// Validate checks that every state has at least one outgoing transition.
func (m *Machine[S, E]) Validate() error {
	for state, evtMap := range m.transitions {
		if len(evtMap) == 0 {
			return fmt.Errorf("state %v has no transitions", state)
		}
	}
	return nil
}

// Memory is an in-memory Adapter.
type Memory[S comparable] struct {
	data map[string]S
}

func NewMemory[S comparable]() *Memory[S] {
	return &Memory[S]{data: make(map[string]S)}
}

func (m *Memory[S]) Get(id string) (*S, error) {
	v, ok := m.data[id]
	if !ok {
		return nil, nil
	}
	return &v, nil
}

func (m *Memory[S]) Set(id string, state S) error {
	m.data[id] = state
	return nil
}

type noopHandler struct{}

func (noopHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (noopHandler) Handle(context.Context, slog.Record) error { return nil }
func (noopHandler) WithAttrs(_ []slog.Attr) slog.Handler      { return noopHandler{} }
func (noopHandler) WithGroup(_ string) slog.Handler           { return noopHandler{} }
