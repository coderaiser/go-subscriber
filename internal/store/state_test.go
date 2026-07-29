package store_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-subscriber/internal/engine"
	"github.com/coderaiser/go-subscriber/internal/store"
)

func TestStateGetMissingReturnsNil(t *testing.T) {
	tape.Test(t, "store/state: Get returns nil for unknown key", func(t *tape.T) {
		s := store.NewStateStore()
		ptr, _ := s.Get("msisdn:svc")
		t.Ok(ptr == nil)
		t.End()
	})
}

func TestStateGetMissingNoError(t *testing.T) {
	tape.Test(t, "store/state: Get does not error for unknown key", func(t *tape.T) {
		s := store.NewStateStore()
		_, err := s.Get("msisdn:svc")
		t.NoError(err)
		t.End()
	})
}

func TestStateSetAndGetNoError(t *testing.T) {
	tape.Test(t, "store/state: Set does not error", func(t *tape.T) {
		s := store.NewStateStore()
		err := s.Set("msisdn:svc", engine.StateActive)
		t.NoError(err)
		t.End()
	})
}

func TestStateGetReturnsSetValue(t *testing.T) {
	tape.Test(t, "store/state: Get returns what was Set", func(t *tape.T) {
		s := store.NewStateStore()
		s.Set("msisdn:svc", engine.StateActive)
		ptr, _ := s.Get("msisdn:svc")
		t.Ok(*ptr == engine.StateActive)
		t.End()
	})
}

func TestStateOverwriteNoError(t *testing.T) {
	tape.Test(t, "store/state: overwrite does not error", func(t *tape.T) {
		s := store.NewStateStore()
		s.Set("k", engine.StateTrial)
		err := s.Set("k", engine.StateSuspended)
		t.NoError(err)
		t.End()
	})
}

func TestStateOverwriteValue(t *testing.T) {
	tape.Test(t, "store/state: Set overwrites previous value", func(t *tape.T) {
		s := store.NewStateStore()
		s.Set("k", engine.StateTrial)
		s.Set("k", engine.StateSuspended)
		ptr, _ := s.Get("k")
		t.Ok(*ptr == engine.StateSuspended)
		t.End()
	})
}

func TestStateIsolation(t *testing.T) {
	tape.Test(t, "store/state: different keys are independent", func(t *tape.T) {
		s := store.NewStateStore()
		s.Set("a", engine.StateActive)
		s.Set("b", engine.StateTerminated)
		pa, _ := s.Get("a")
		pb, _ := s.Get("b")
		t.Ok(*pa != *pb)
		t.End()
	})
}
