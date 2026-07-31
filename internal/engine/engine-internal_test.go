package engine

import (
	"log/slog"
	"testing"
	"time"

	. "github.com/coderaiser/go-subscriber/internal/tape"
	"github.com/coderaiser/go-subscriber/internal/statemachine"
	"github.com/coderaiser/go-subscriber/internal/store"
)

var epoch = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func fixed(t time.Time) func() time.Time { return func() time.Time { return t } }

// --- parseState (line 40) ---

func TestParseStateUnknown(t *testing.T) {
	Test(t, "engine: parseState returns error for unknown state", func(t *T) {
		_, err := parseState("bogus")
		t.Error(err)
		t.End()
	})
}

func TestParseStateKnown(t *testing.T) {
	Test(t, "engine: parseState accepts all known states", func(t *T) {
		for _, s := range []string{StateTrial, StateActive, StateSuspended, StateTerminated, StateRemoved} {
			got, err := parseState(s)
			t.NoError(err)
			t.Equal(got, s)
		}
		t.End()
	})
}

// --- parseEvent (line 50) ---

func TestParseEventUnknown(t *testing.T) {
	Test(t, "engine: parseEvent returns error for unknown event", func(t *T) {
		_, err := parseEvent("bogus")
		t.Error(err)
		t.End()
	})
}

func TestParseEventKnown(t *testing.T) {
	Test(t, "engine: parseEvent accepts all known events", func(t *T) {
		for _, e := range []string{EventExpireSuccess, EventExpireFail, EventRenewSuccess, EventRenewFail, EventRetrySuccess, EventKickOut, EventUnsubscribe} {
			got, err := parseEvent(e)
			t.NoError(err)
			t.Equal(got, e)
		}
		t.End()
	})
}

// --- panic in newWithSource (lines 54-55) ---

func TestNewWithSourcePanicsOnBadTransition(t *testing.T) {
	Test(t, "engine: newWithSource panics when transition references unknown state", func(t *T) {
		defer func() {
			r := recover()
			t.Ok(r != nil)
			t.End()
		}()
		broken := &statemachine.MemorySource{
			Defs: []statemachine.TransitionDef{
				{From: "bogus_state", Event: EventUnsubscribe, To: StateTerminated},
			},
		}
		newWithSource(broken, store.NewStateStore(), store.NewFactsStore(), fixed(epoch), slog.Default())
	})
}

// --- Subscribe default ChargeResult (line 100-102) ---

func TestSubscribeUnknownChargeResult(t *testing.T) {
	Test(t, "engine: Subscribe with unknown charge result returns ErrChargeFailed", func(t *T) {
		eng := New(store.NewStateStore(), store.NewFactsStore(), fixed(epoch), slog.Default())
		err := eng.Subscribe("111", "svc1", false, "totally_unknown_result")
		t.Error(err)
		t.End()
	})
}

// --- Unsubscribe machine.Apply error (lines 144-146) ---

func TestUnsubscribeApplyError(t *testing.T) {
	Test(t, "engine: Unsubscribe returns StateTerminated when machine.Apply errors", func(t *T) {
		// Build a machine with no transitions defined for EventUnsubscribe from StateActive,
		// using strict=true so Apply returns an error on undefined transitions.
		ss := store.NewStateStore()
		fs := store.NewFactsStore()
		noUnsubscribe := &statemachine.MemorySource{
			Defs: []statemachine.TransitionDef{
				// active can renew but has no unsubscribe transition
				{From: StateActive, Event: EventRenewSuccess, To: StateActive},
			},
		}
		eng := newWithSource(noUnsubscribe, ss, fs, fixed(epoch), slog.Default())
		// Set state directly so we bypass Subscribe (which needs more transitions)
		ss.Set("111:svc1", StateActive)
		state := eng.Unsubscribe("111", "svc1")
		t.Ok(state == StateTerminated)
		t.End()
	})
}

// --- OnChargeResult final return state fallthrough (line 188) ---

func TestOnChargeResultUnknownResult(t *testing.T) {
	Test(t, "engine: OnChargeResult with unknown result returns current state unchanged", func(t *T) {
		eng := New(store.NewStateStore(), store.NewFactsStore(), fixed(epoch), slog.Default())
		eng.Subscribe("111", "svc1", false)
		state := eng.OnChargeResult("111", "svc1", "totally_unknown_result")
		t.Ok(state == StateActive)
		t.End()
	})
}