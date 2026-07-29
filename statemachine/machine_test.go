package statemachine

import (
	"errors"
	"testing"

	"github.com/coderaiser/go-subscriber/statemachine/adapters"
)

type testState int

const (
	stateIdle testState = iota
	stateRunning
	stateDone
)

type testEvent int

const (
	eventRun testEvent = iota
	eventFinish
	eventFail
)

func parseState(s string) (testState, error) {
	switch s {
	case "idle":
		return stateIdle, nil
	case "running":
		return stateRunning, nil
	case "done":
		return stateDone, nil
	default:
		return 0, errors.New("unknown state")
	}
}

func parseEvent(e string) (testEvent, error) {
	switch e {
	case "run":
		return eventRun, nil
	case "finish":
		return eventFinish, nil
	case "fail":
		return eventFail, nil
	default:
		return 0, errors.New("unknown event")
	}
}

func TestNewValidMachine(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
			{From: "running", Event: "finish", To: "done"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, err := New(src, parseState, parseEvent, adapter, false)
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("expected machine")
	}
}

func TestApplyValidTransition(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	next, err := m.Apply("test-1", eventRun, nil)
	if err != nil {
		t.Fatal(err)
	}
	if next != stateRunning {
		t.Errorf("want %v, got %v", stateRunning, next)
	}
}

func TestApplyInvalidTransitionNonStrict(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	_, err := m.Apply("test-1", eventFinish, nil)
	if err == nil {
		t.Fatal("expected error for invalid transition")
	}
}

func TestApplyInvalidTransitionStrictPanics(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, true)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for strict invalid transition")
		}
	}()

	m.Apply("test-1", eventFinish, nil)
}

func TestHookCalledOnTransition(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	called := false
	m.Hook(stateIdle, eventRun, func(ctx Context[testState, testEvent]) error {
		called = true
		return nil
	})

	m.Apply("test-1", eventRun, nil)
	if !called {
		t.Error("expected hook to be called")
	}
}

func TestHookErrorReturned(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	m.Hook(stateIdle, eventRun, func(ctx Context[testState, testEvent]) error {
		return errors.New("hook failed")
	})

	_, err := m.Apply("test-1", eventRun, nil)
	if err == nil || err.Error() != "hook failed" {
		t.Fatalf("expected hook error, got %v", err)
	}
}

func TestApplyStoresState(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	m.Apply("test-1", eventRun, nil)
	ptr, err := adapter.Get("test-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ptr == nil {
		t.Fatal("expected stored state, got nil")
	}
	if *ptr != stateRunning {
		t.Errorf("want %v, got %v", stateRunning, *ptr)
	}
}

func TestUnknownIdGetsNil(t *testing.T) {
	adapter := adapters.NewMemory[testState]()
	ptr, err := adapter.Get("unknown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ptr != nil {
		t.Errorf("expected nil, got %v", ptr)
	}
}

func TestNewWithParseStateError(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "invalid", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	_, err := New(src, parseState, parseEvent, adapter, false)
	if err == nil {
		t.Fatal("expected error for invalid state")
	}
}

func TestNewWithParseEventError(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "invalid", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	_, err := New(src, parseState, parseEvent, adapter, false)
	if err == nil {
		t.Fatal("expected error for invalid event")
	}
}

func TestNewWithParseToStateError(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "invalid"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	_, err := New(src, parseState, parseEvent, adapter, false)
	if err == nil {
		t.Fatal("expected error for invalid to state")
	}
}

func TestValidatePasses(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	if err := m.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestHookNotCalledForUnknownTransition(t *testing.T) {
	src := &MemorySource{
		Defs: []TransitionDef{
			{From: "idle", Event: "run", To: "running"},
		},
	}
	adapter := adapters.NewMemory[testState]()
	m, _ := New(src, parseState, parseEvent, adapter, false)

	called := false
	m.Hook(stateRunning, eventFinish, func(ctx Context[testState, testEvent]) error {
		called = true
		return nil
	})

	m.Apply("test-1", eventFinish, nil)
	if called {
		t.Error("hook should not be called for invalid transition")
	}
}

// -- error adapter tests --

type errGetAdapter struct{}

func (a errGetAdapter) Get(id string) (*testState, error) { return nil, errors.New("get failed") }
func (a errGetAdapter) Set(id string, s testState) error  { return nil }

type errSetAdapter struct{}

func (a errSetAdapter) Get(id string) (*testState, error) { s := stateIdle; return &s, nil }
func (a errSetAdapter) Set(id string, s testState) error  { return errors.New("set failed") }

func TestApplyGetError(t *testing.T) {
	src := &MemorySource{Defs: []TransitionDef{{From: "idle", Event: "run", To: "running"}}}
	m, _ := New(src, parseState, parseEvent, errGetAdapter{}, false)
	_, err := m.Apply("x", eventRun, nil)
	if err == nil {
		t.Fatal("expected error from Get failure")
	}
}

func TestApplySetError(t *testing.T) {
	src := &MemorySource{Defs: []TransitionDef{{From: "idle", Event: "run", To: "running"}}}
	m, _ := New(src, parseState, parseEvent, errSetAdapter{}, false)
	_, err := m.Apply("x", eventRun, nil)
	if err == nil {
		t.Fatal("expected error from Set failure")
	}
}

func TestWithInitialUsedForUnknownId(t *testing.T) {
	src := &MemorySource{Defs: []TransitionDef{{From: "idle", Event: "run", To: "running"}}}
	m, _ := New(src, parseState, parseEvent, adapters.NewMemory[testState](), false)
	m.WithInitial(stateIdle)
	next, err := m.Apply("brand-new", eventRun, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if next != stateRunning {
		t.Errorf("want running, got %v", next)
	}
}

func TestValidateEmptyTransitions(t *testing.T) {
	m := &Machine[testState, testEvent]{
		transitions: map[testState]map[testEvent]testState{
			stateIdle: {},
		},
	}
	if err := m.Validate(); err == nil {
		t.Fatal("expected error for state with no transitions out")
	}
}

// -- FileSource tests --

func TestFileSourceLoadsTransitions(t *testing.T) {
	src := FileSource{Path: "testdata/runner.toml"}
	defs, err := src.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(defs) == 0 {
		t.Fatal("expected transitions")
	}
}

func TestFileSourceJSON(t *testing.T) {
	src := FileSource{Path: "testdata/runner.json"}
	defs, err := src.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(defs) == 0 {
		t.Fatal("expected transitions")
	}
}

// errSource returns an error on Load.
type errSource struct{}

func (s errSource) Load() ([]TransitionDef, error) {
	return nil, errors.New("source failed")
}

func TestNewWithSourceError(t *testing.T) {
	_, err := New[testState, testEvent](errSource{}, parseState, parseEvent, adapters.NewMemory[testState](), false)
	if err == nil {
		t.Fatal("expected error from source failure")
	}
}

func TestFileSourceUnsupportedExtension(t *testing.T) {
	src := FileSource{Path: "testdata/runner.yaml"}
	_, err := src.Load()
	if err == nil {
		t.Fatal("expected error for unsupported extension")
	}
}

func TestFileSourceInvalidJSON(t *testing.T) {
	src := FileSource{Path: "testdata/invalid.json"}
	_, err := src.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFileSourceMissingJSON(t *testing.T) {
	src := FileSource{Path: "testdata/nonexistent.json"}
	_, err := src.Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFileSourceMissingTOML(t *testing.T) {
	src := FileSource{Path: "nonexistent.toml"}
	_, err := src.Load()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
