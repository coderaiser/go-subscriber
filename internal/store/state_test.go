package store_test

import (
	"testing"

	"github.com/coderaiser/go-subscriber/internal/store"
	Test "github.com/coderaiser/go-subscriber/internal/tape"
)

func TestStateGetMissingReturnsNil(t *testing.T) {
	Test.Test(t, "store/state: Get returns nil for unknown key", func(t *Test.T) {
		s := store.NewStateStore()
		ptr, _ := s.Get("msisdn:svc")
		t.Ok(ptr == nil)
		t.End()
	})
}

func TestStateGetMissingNoError(t *testing.T) {
	Test.Test(t, "store/state: Get does not error for unknown key", func(t *Test.T) {
		s := store.NewStateStore()
		_, err := s.Get("msisdn:svc")
		t.NoError(err)
		t.End()
	})
}

func TestStateSetAndGetNoError(t *testing.T) {
	Test.Test(t, "store/state: Set does not error", func(t *Test.T) {
		s := store.NewStateStore()
		err := s.Set("msisdn:svc", "active")
		t.NoError(err)
		t.End()
	})
}

func TestStateGetReturnsSetValue(t *testing.T) {
	Test.Test(t, "store/state: Get returns what was Set", func(t *Test.T) {
		s := store.NewStateStore()
		s.Set("msisdn:svc", "active")
		ptr, _ := s.Get("msisdn:svc")
		t.Ok(*ptr == "active")
		t.End()
	})
}

func TestStateOverwriteNoError(t *testing.T) {
	Test.Test(t, "store/state: overwrite does not error", func(t *Test.T) {
		s := store.NewStateStore()
		s.Set("k", "trial")
		err := s.Set("k", "suspended")
		t.NoError(err)
		t.End()
	})
}

func TestStateOverwriteValue(t *testing.T) {
	Test.Test(t, "store/state: Set overwrites previous value", func(t *Test.T) {
		s := store.NewStateStore()
		s.Set("k", "trial")
		s.Set("k", "suspended")
		ptr, _ := s.Get("k")
		t.Ok(*ptr == "suspended")
		t.End()
	})
}

func TestStateIsolation(t *testing.T) {
	Test.Test(t, "store/state: different keys are independent", func(t *Test.T) {
		s := store.NewStateStore()
		s.Set("a", "active")
		s.Set("b", "terminated")
		pa, _ := s.Get("a")
		pb, _ := s.Get("b")
		t.Ok(*pa != *pb)
		t.End()
	})
}

func TestStateAllEmpty(t *testing.T) {
	Test.Test(t, "store/state: All returns empty map when no states set", func(t *Test.T) {
		s := store.NewStateStore()
		all := s.All()
		t.Ok(len(all) == 0)
		t.End()
	})
}

func TestStateAllReturnsAllKeys(t *testing.T) {
	Test.Test(t, "store/state: All returns all set keys", func(t *Test.T) {
		s := store.NewStateStore()
		s.Set("a:svc1", "active")
		s.Set("b:svc1", "trial")
		all := s.All()
		t.Ok(len(all) == 2)
		t.End()
	})
}

func TestStateAllIsSnapshot(t *testing.T) {
	Test.Test(t, "store/state: All returns a copy not a reference", func(t *Test.T) {
		s := store.NewStateStore()
		s.Set("a:svc1", "active")
		all := s.All()
		all["a:svc1"] = "mutated"
		ptr, _ := s.Get("a:svc1")
		t.Ok(*ptr == "active")
		t.End()
	})
}
