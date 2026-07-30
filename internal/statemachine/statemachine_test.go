package statemachine_test

import (
	"errors"
	"testing"

	Test "github.com/coderaiser/go-subscriber/internal/tape"
	"github.com/coderaiser/go-subscriber/internal/statemachine"
)

type state int

const (
	stateIdle state = iota
	stateRunning
	stateDone
)

type event int

const (
	eventRun event = iota
	eventFinish
	eventFail
)

func parseState(s string) (state, error) {
	switch s {
	case "idle":
		return stateIdle, nil
	case "running":
		return stateRunning, nil
	case "done":
		return stateDone, nil
	}
	return 0, errors.New("unknown state: " + s)
}

func parseEvent(e string) (event, error) {
	switch e {
	case "run":
		return eventRun, nil
	case "finish":
		return eventFinish, nil
	case "fail":
		return eventFail, nil
	}
	return 0, errors.New("unknown event: " + e)
}

func newMachine(t *testing.T) *statemachine.Machine[state, event] {
	t.Helper()
	src := &statemachine.MemorySource{
		Defs: []statemachine.TransitionDef{
			{From: "idle", Event: "run", To: "running"},
			{From: "running", Event: "finish", To: "done"},
			{From: "running", Event: "fail", To: "done"},
		},
	}
	m, err := statemachine.New(src, parseState, parseEvent, statemachine.NewMemory[state](), false)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestMachineValidTransition(t *testing.T) {
	Test.Test(t, "statemachine: valid transition returns next state", func(t *Test.T) {
		m := newMachine(t.TB())
		_, err := m.Apply("e1", eventRun, nil)
		t.NoError(err)
		t.End()
	})
}

func TestMachineStoresState(t *testing.T) {
	Test.Test(t, "statemachine: state is persisted after Apply", func(t *Test.T) {
		mem := statemachine.NewMemory[state]()
		src := &statemachine.MemorySource{
			Defs: []statemachine.TransitionDef{
				{From: "idle", Event: "run", To: "running"},
			},
		}
		m, _ := statemachine.New(src, parseState, parseEvent, mem, false)
		m.Apply("e1", eventRun, nil)
		ptr, _ := mem.Get("e1")
		t.Ok(ptr != nil)
		t.End()
	})
}

func TestMachineStoresStateValue(t *testing.T) {
	Test.Test(t, "statemachine: stored state matches running", func(t *Test.T) {
		mem := statemachine.NewMemory[state]()
		src := &statemachine.MemorySource{
			Defs: []statemachine.TransitionDef{
				{From: "idle", Event: "run", To: "running"},
			},
		}
		m, _ := statemachine.New(src, parseState, parseEvent, mem, false)
		m.Apply("e1", eventRun, nil)
		ptr, _ := mem.Get("e1")
		t.Ok(*ptr == stateRunning)
		t.End()
	})
}

func TestMachineInvalidTransitionNonStrict(t *testing.T) {
	Test.Test(t, "statemachine: invalid transition returns error in non-strict mode", func(t *Test.T) {
		m := newMachine(t.TB())
		_, err := m.Apply("e1", eventFinish, nil)
		t.Error(err)
		t.End()
	})
}

func TestMachineInvalidTransitionStrict(t *testing.T) {
	Test.Test(t, "statemachine: invalid transition panics in strict mode", func(t *Test.T) {
		src := &statemachine.MemorySource{
			Defs: []statemachine.TransitionDef{
				{From: "idle", Event: "run", To: "running"},
			},
		}
		m, _ := statemachine.New(src, parseState, parseEvent, statemachine.NewMemory[state](), true)
		defer func() {
			t.Ok(recover() != nil)
			t.End()
		}()
		m.Apply("e1", eventFinish, nil)
	})
}

func TestMachineHookCalled(t *testing.T) {
	Test.Test(t, "statemachine: hook is called after transition", func(t *Test.T) {
		m := newMachine(t.TB())
		called := false
		m.Hook(stateIdle, eventRun, func(ctx statemachine.Context[state, event]) error {
			called = true
			return nil
		})
		m.Apply("e1", eventRun, nil)
		t.Ok(called)
		t.End()
	})
}

func TestMachineHookError(t *testing.T) {
	Test.Test(t, "statemachine: hook error is returned", func(t *Test.T) {
		m := newMachine(t.TB())
		m.Hook(stateIdle, eventRun, func(ctx statemachine.Context[state, event]) error {
			return errors.New("hook failed")
		})
		_, err := m.Apply("e1", eventRun, nil)
		t.Error(err)
		t.End()
	})
}

func TestMachineWithInitialNoError(t *testing.T) {
	Test.Test(t, "statemachine: WithInitial allows Apply for unknown id", func(t *Test.T) {
		m := newMachine(t.TB())
		m.WithInitial(stateIdle)
		_, err := m.Apply("brand-new", eventRun, nil)
		t.NoError(err)
		t.End()
	})
}

func TestMachineValidate(t *testing.T) {
	Test.Test(t, "statemachine: Validate passes for valid machine", func(t *Test.T) {
		m := newMachine(t.TB())
		t.NoError(m.Validate())
		t.End()
	})
}

func TestMachineValidateEmpty(t *testing.T) {
	Test.Test(t, "statemachine: Validate fails for state with no transitions", func(t *Test.T) {
		m := &statemachine.Machine[state, event]{}
		m.SetTransitions(map[state]map[event]state{
			stateIdle: {},
		})
		t.Error(m.Validate())
		t.End()
	})
}

func TestMemoryGetMissingReturnsNil(t *testing.T) {
	Test.Test(t, "statemachine: Memory.Get returns nil for unknown id", func(t *Test.T) {
		mem := statemachine.NewMemory[state]()
		ptr, _ := mem.Get("unknown")
		t.Ok(ptr == nil)
		t.End()
	})
}

func TestMemorySetAndGetNoError(t *testing.T) {
	Test.Test(t, "statemachine: Memory.Set does not error", func(t *Test.T) {
		mem := statemachine.NewMemory[state]()
		err := mem.Set("e1", stateRunning)
		t.NoError(err)
		t.End()
	})
}

func TestMemorySetAndGetCorrectValue(t *testing.T) {
	Test.Test(t, "statemachine: Memory.Get returns what was Set", func(t *Test.T) {
		mem := statemachine.NewMemory[state]()
		mem.Set("e1", stateRunning)
		ptr, _ := mem.Get("e1")
		t.Ok(*ptr == stateRunning)
		t.End()
	})
}

func TestNewWithBadFromState(t *testing.T) {
	Test.Test(t, "statemachine: New returns error for unknown From state", func(t *Test.T) {
		src := &statemachine.MemorySource{
			Defs: []statemachine.TransitionDef{
				{From: "bad", Event: "run", To: "running"},
			},
		}
		_, err := statemachine.New(src, parseState, parseEvent, statemachine.NewMemory[state](), false)
		t.Error(err)
		t.End()
	})
}

func TestNewWithBadEvent(t *testing.T) {
	Test.Test(t, "statemachine: New returns error for unknown Event", func(t *Test.T) {
		src := &statemachine.MemorySource{
			Defs: []statemachine.TransitionDef{
				{From: "idle", Event: "bad", To: "running"},
			},
		}
		_, err := statemachine.New(src, parseState, parseEvent, statemachine.NewMemory[state](), false)
		t.Error(err)
		t.End()
	})
}

func TestNewWithBadToState(t *testing.T) {
	Test.Test(t, "statemachine: New returns error for unknown To state", func(t *Test.T) {
		src := &statemachine.MemorySource{
			Defs: []statemachine.TransitionDef{
				{From: "idle", Event: "run", To: "bad"},
			},
		}
		_, err := statemachine.New(src, parseState, parseEvent, statemachine.NewMemory[state](), false)
		t.Error(err)
		t.End()
	})
}

func TestAdapterGetError(t *testing.T) {
	Test.Test(t, "statemachine: Apply returns error when adapter.Get fails", func(t *Test.T) {
		src := &statemachine.MemorySource{
			Defs: []statemachine.TransitionDef{
				{From: "idle", Event: "run", To: "running"},
			},
		}
		m, _ := statemachine.New(src, parseState, parseEvent, &errGetAdapter{}, false)
		_, err := m.Apply("x", eventRun, nil)
		t.Error(err)
		t.End()
	})
}

func TestAdapterSetError(t *testing.T) {
	Test.Test(t, "statemachine: Apply returns error when adapter.Set fails", func(t *Test.T) {
		src := &statemachine.MemorySource{
			Defs: []statemachine.TransitionDef{
				{From: "idle", Event: "run", To: "running"},
			},
		}
		m, _ := statemachine.New(src, parseState, parseEvent, &errSetAdapter{}, false)
		_, err := m.Apply("x", eventRun, nil)
		t.Error(err)
		t.End()
	})
}

type errGetAdapter struct{}

func (a *errGetAdapter) Get(id string) (*state, error) { return nil, errors.New("get failed") }
func (a *errGetAdapter) Set(id string, s state) error  { return nil }

type errSetAdapter struct{}

func (a *errSetAdapter) Get(id string) (*state, error) { s := stateIdle; return &s, nil }
func (a *errSetAdapter) Set(id string, s state) error  { return errors.New("set failed") }
