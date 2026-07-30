package store_test

import (
	"testing"

	Test "github.com/coderaiser/go-subscriber/internal/tape"
	"github.com/coderaiser/go-subscriber/internal/store"
)

func TestFactsGetMissing(t *testing.T) {
	Test.Test(t, "store/facts: Get returns zero facts for unknown key", func(t *Test.T) {
		f := store.NewFactsStore()
		facts := f.Get("msisdn:svc")
		t.Ok(facts.CooloffUntil == 0)
		t.End()
	})
}

func TestFactsSetAndGetCooloff(t *testing.T) {
	Test.Test(t, "store/facts: Set then Get returns CooloffUntil", func(t *Test.T) {
		f := store.NewFactsStore()
		f.Set("k", store.Facts{CooloffUntil: 999, FallbackAttempt: 2})
		facts := f.Get("k")
		t.Ok(facts.CooloffUntil == 999)
		t.End()
	})
}

func TestFactsSetAndGetFallback(t *testing.T) {
	Test.Test(t, "store/facts: Set then Get returns FallbackAttempt", func(t *Test.T) {
		f := store.NewFactsStore()
		f.Set("k", store.Facts{CooloffUntil: 999, FallbackAttempt: 2})
		facts := f.Get("k")
		t.Ok(facts.FallbackAttempt == 2)
		t.End()
	})
}

func TestFactsOverwrite(t *testing.T) {
	Test.Test(t, "store/facts: Set overwrites previous facts", func(t *Test.T) {
		f := store.NewFactsStore()
		f.Set("k", store.Facts{CooloffUntil: 1})
		f.Set("k", store.Facts{CooloffUntil: 2})
		facts := f.Get("k")
		t.Ok(facts.CooloffUntil == 2)
		t.End()
	})
}
