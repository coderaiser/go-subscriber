package state

import (
	"errors"
	"testing"

	"github.com/coderaiser/go-subscriber/model"
	"github.com/coderaiser/go-subscriber/statemachine"
	"github.com/coderaiser/go-subscriber/statemachine/adapters"
)

func TestRunEventCreatesRunningState(t *testing.T) {
	s := New()
	_, err := s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	if err != nil {
		t.Fatal(err)
	}
	st, _ := s.Get("TestFoo")
	if st != StateRunning {
		t.Errorf("want %v, got %v", StateRunning, st)
	}
}

func TestPassEventMarksTestPassed(t *testing.T) {
	s := New()
	s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	s.Apply(model.Event{Action: "pass", Test: "TestFoo"})
	st, _ := s.Get("TestFoo")
	if st != StatePassed {
		t.Errorf("want %v, got %v", StatePassed, st)
	}
}

func TestFailEventMarksTestFailed(t *testing.T) {
	s := New()
	s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	s.Apply(model.Event{Action: "fail", Test: "TestFoo"})
	st, _ := s.Get("TestFoo")
	if st != StateFailed {
		t.Errorf("want %v, got %v", StateFailed, st)
	}
}

func TestSkipEventMarksTestSkipped(t *testing.T) {
	s := New()
	s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	s.Apply(model.Event{Action: "skip", Test: "TestFoo"})
	st, _ := s.Get("TestFoo")
	if st != StateSkipped {
		t.Errorf("want %v, got %v", StateSkipped, st)
	}
}

func TestOutputAppendedToLogs(t *testing.T) {
	s := New()
	s.Apply(model.Event{Action: "output", Test: "TestFoo", Output: "line1\n"})
	s.Apply(model.Event{Action: "output", Test: "TestFoo", Output: "line2\n"})
	got := s.GetOutput("TestFoo")
	if got != "line1\nline2\n" {
		t.Errorf("want line1\\nline2\\n, got %s", got)
	}
}

func TestPackageEventNoTestID(t *testing.T) {
	s := New()
	_, err := s.Apply(model.Event{Action: "pass", Package: "mypkg"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidActionError(t *testing.T) {
	s := New()
	_, err := s.Apply(model.Event{Test: "TestFoo", Action: "invalid"})
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestRunTwice(t *testing.T) {
	s := New()
	s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	s.Apply(model.Event{Action: "run", Test: "TestFoo"})
	st, _ := s.Get("TestFoo")
	if st != StateRunning {
		t.Errorf("want %v, got %v", StateRunning, st)
	}
}

func TestGetNonExistent(t *testing.T) {
	s := New()
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent test")
	}
}

func TestApplyInvalidTransitionReturnsCurrentState(t *testing.T) {
	s := New()
	// pass without run — invalid transition
	state, err := s.Apply(model.Event{Action: "pass", Test: "TestA"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != StateIdle {
		t.Errorf("want StateIdle, got %v", state)
	}
}

func TestSummaryHierarchicalName(t *testing.T) {
	s := New()
	// state stored under "pkg" — output maps pkg/TestFoo -> pkg
	s.Apply(model.Event{Action: "run", Test: "pkg"})
	s.Apply(model.Event{Action: "pass", Test: "pkg"})
	// output with subtest name "pkg/TestFoo" — LastIndex extracts "pkg"
	s.Apply(model.Event{Action: "output", Test: "pkg/TestFoo", Output: "ok"})
	passed, _, _ := s.Summary()
	if len(passed) != 1 {
		t.Errorf("want 1 passed, got %d", len(passed))
	}
}

func TestSummaryPassedFailedSkipped(t *testing.T) {
	s := New()
	s.Apply(model.Event{Action: "run", Test: "TestA"})
	s.Apply(model.Event{Action: "pass", Test: "TestA"})
	s.Apply(model.Event{Action: "output", Test: "TestA", Output: "ok"})
	s.Apply(model.Event{Action: "run", Test: "TestB"})
	s.Apply(model.Event{Action: "fail", Test: "TestB"})
	s.Apply(model.Event{Action: "output", Test: "TestB", Output: "fail"})
	s.Apply(model.Event{Action: "run", Test: "TestC"})
	s.Apply(model.Event{Action: "skip", Test: "TestC"})
	s.Apply(model.Event{Action: "output", Test: "TestC", Output: "skip"})
	passed, failed, skipped := s.Summary()
	if len(passed) != 1 {
		t.Errorf("want 1 passed, got %d", len(passed))
	}
	if len(failed) != 1 {
		t.Errorf("want 1 failed, got %d", len(failed))
	}
	if len(skipped) != 1 {
		t.Errorf("want 1 skipped, got %d", len(skipped))
	}
}

func TestParseTestStateUnknown(t *testing.T) {
	_, err := parseTestState("unknown")
	if err == nil {
		t.Fatal("expected error for unknown state")
	}
}

func TestParseTestEventUnknown(t *testing.T) {
	_, err := parseTestEvent("unknown")
	if err == nil {
		t.Fatal("expected error for unknown event")
	}
}

func TestSummaryWithOutputOnly(t *testing.T) {
	s := New()
	// output-only test — never applied with a state transition
	s.Apply(model.Event{Action: "output", Test: "orphan", Output: "some output"})
	passed, failed, skipped := s.Summary()
	if len(passed)+len(failed)+len(skipped) != 0 {
		t.Errorf("expected 0 total, got pass=%d fail=%d skip=%d", len(passed), len(failed), len(skipped))
	}
}

func TestNewPanic(t *testing.T) {
	s := New()
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

// errSource always fails Load
type errSource struct{}

func (s errSource) Load() ([]statemachine.TransitionDef, error) {
	return nil, errors.New("source failed")
}

func TestNewPanicsOnBadSource(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from bad source")
		}
	}()
	newFromSource(errSource{})
}

// errAdapter always fails Get
type errAdapter struct{ *adapters.Memory[TestState] }

func (a errAdapter) Get(id string) (*TestState, error) {
	return nil, errors.New("adapter error")
}

func TestGetAdapterError(t *testing.T) {
	s := New()
	mem, ok := s.adapter.(*adapters.Memory[TestState])
	if !ok {
		t.Fatal("expected memory adapter")
	}
	s.adapter = errAdapter{mem}
	_, err := s.Get("TestFoo")
	if err == nil {
		t.Fatal("expected error from adapter")
	}
}

func TestSummaryAdapterError(t *testing.T) {
	s := New()
	s.Apply(model.Event{Action: "output", Test: "TestFoo", Output: "ok"})
	mem, ok := s.adapter.(*adapters.Memory[TestState])
	if !ok {
		t.Fatal("expected memory adapter")
	}
	s.adapter = errAdapter{mem}
	// Summary silently skips on adapter error — should not panic
	passed, failed, skipped := s.Summary()
	_ = passed
	_ = failed
	_ = skipped
}
