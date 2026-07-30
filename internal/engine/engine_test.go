package engine_test

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	Test "github.com/coderaiser/go-subscriber/internal/tape"
	"github.com/coderaiser/go-subscriber/internal/engine"
	"github.com/coderaiser/go-subscriber/internal/store"
)

var epoch = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func newEngine(t *testing.T, now func() time.Time) *engine.Engine {
	t.Helper()
	return engine.New(
		store.NewStateStore(),
		store.NewFactsStore(),
		now,
		slog.Default(),
	)
}

func fixed(t time.Time) func() time.Time { return func() time.Time { return t } }

var errCooloff = errors.New("cooloff")
var errAlreadySubscribed = errors.New("already subscribed")

func TestSubscribeTrial(t *testing.T) {
	Test.Test(t, "engine: subscribe trial sets trial state", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		err := eng.Subscribe("111", "svc1", true)
		t.NoError(err)
		t.End()
	})
}

func TestSubscribePaid(t *testing.T) {
	Test.Test(t, "engine: subscribe paid sets active state", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		err := eng.Subscribe("111", "svc1", false)
		t.NoError(err)
		t.End()
	})
}

func TestSubscribeAlreadySubscribed(t *testing.T) {
	Test.Test(t, "engine: subscribe returns error when already subscribed", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		err := eng.Subscribe("111", "svc1", false)
		t.Error(err)
		t.End()
	})
}

func TestSubscribeCooloffActive(t *testing.T) {
	Test.Test(t, "engine: subscribe returns cooloff error during cooloff", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		eng.Unsubscribe("111", "svc1")
		err := eng.Subscribe("111", "svc1", false)
		t.Ok(err != nil)
		t.End()
	})
}

func TestSubscribeAfterCooloff(t *testing.T) {
	Test.Test(t, "engine: subscribe succeeds after cooloff expires", func(t *Test.T) {
		ss := store.NewStateStore()
		fs := store.NewFactsStore()
		eng := engine.New(ss, fs, fixed(epoch), slog.Default())
		eng.Subscribe("222", "svc1", false)
		eng.Unsubscribe("222", "svc1")
		eng2 := engine.New(ss, fs, fixed(epoch.AddDate(0, 0, 31)), slog.Default())
		err := eng2.Subscribe("222", "svc1", false)
		t.NoError(err)
		t.End()
	})
}

func TestUnsubscribeFromActive(t *testing.T) {
	Test.Test(t, "engine: unsubscribe from active sets terminated", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.Unsubscribe("111", "svc1")
		t.Ok(state == engine.StateTerminated)
		t.End()
	})
}

func TestUnsubscribeFromTrial(t *testing.T) {
	Test.Test(t, "engine: unsubscribe from trial sets terminated", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", true)
		state := eng.Unsubscribe("111", "svc1")
		t.Ok(state == engine.StateTerminated)
		t.End()
	})
}

func TestUnsubscribeFromSuspended(t *testing.T) {
	Test.Test(t, "engine: unsubscribe from suspended sets terminated", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		eng.OnChargeResult("111", "svc1", engine.ResultPermanent)
		state := eng.Unsubscribe("111", "svc1")
		t.Ok(state == engine.StateTerminated)
		t.End()
	})
}

func TestUnsubscribeNoSubscription(t *testing.T) {
	Test.Test(t, "engine: unsubscribe with no subscription does not panic", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		state := eng.Unsubscribe("nobody", "svc1")
		t.Ok(state == engine.StateTerminated)
		t.End()
	})
}

func TestChargeSuccess(t *testing.T) {
	Test.Test(t, "engine: charge success sets active", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.OnChargeResult("111", "svc1", engine.ResultSuccess)
		t.Ok(state == engine.StateActive)
		t.End()
	})
}

func TestChargeLowBalanceBelowThreshold(t *testing.T) {
	Test.Test(t, "engine: low balance below threshold keeps current state", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.OnChargeResult("111", "svc1", engine.ResultLowBalance)
		t.Ok(state == engine.StateActive)
		t.End()
	})
}

func TestChargeLowBalanceAtThreshold(t *testing.T) {
	Test.Test(t, "engine: low balance at threshold suspends", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		eng.OnChargeResult("111", "svc1", engine.ResultLowBalance)
		eng.OnChargeResult("111", "svc1", engine.ResultLowBalance)
		state := eng.OnChargeResult("111", "svc1", engine.ResultLowBalance)
		t.Ok(state == engine.StateSuspended)
		t.End()
	})
}

func TestChargePermanent(t *testing.T) {
	Test.Test(t, "engine: permanent failure suspends", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.OnChargeResult("111", "svc1", engine.ResultPermanent)
		t.Ok(state == engine.StateSuspended)
		t.End()
	})
}

func TestChargeRateLimit(t *testing.T) {
	Test.Test(t, "engine: rate limit leaves state unchanged", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.OnChargeResult("111", "svc1", engine.ResultRateLimit)
		t.Ok(state == engine.StateActive)
		t.End()
	})
}

func TestChargeSystemError(t *testing.T) {
	Test.Test(t, "engine: system error leaves state unchanged", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.OnChargeResult("111", "svc1", engine.ResultSystemError)
		t.Ok(state == engine.StateActive)
		t.End()
	})
}

func TestChargePending(t *testing.T) {
	Test.Test(t, "engine: pending leaves state unchanged", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.OnChargeResult("111", "svc1", engine.ResultPending)
		t.Ok(state == engine.StateActive)
		t.End()
	})
}

func TestChargeNoResponse(t *testing.T) {
	Test.Test(t, "engine: no response leaves state unchanged", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.OnChargeResult("111", "svc1", engine.ResultNoResponse)
		t.Ok(state == engine.StateActive)
		t.End()
	})
}

func TestChargeSubscriberState(t *testing.T) {
	Test.Test(t, "engine: subscriber_state result suspends", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.OnChargeResult("111", "svc1", engine.ResultSubscriberState)
		t.Ok(state == engine.StateSuspended)
		t.End()
	})
}

func TestChargeUnknownSub(t *testing.T) {
	Test.Test(t, "engine: charge result for unknown subscription returns empty", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		state := eng.OnChargeResult("nobody", "svc1", engine.ResultSuccess)
		t.Ok(state == "")
		t.End()
	})
}

func TestExpireTrialSuccess(t *testing.T) {
	Test.Test(t, "engine: expire trial success sets active", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", true)
		state := eng.ExpireTrial("111", "svc1", true)
		t.Ok(state == engine.StateActive)
		t.End()
	})
}

func TestExpireTrialFail(t *testing.T) {
	Test.Test(t, "engine: expire trial fail suspends", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", true)
		state := eng.ExpireTrial("111", "svc1", false)
		t.Ok(state == engine.StateSuspended)
		t.End()
	})
}

func TestExpireTrialNotInTrial(t *testing.T) {
	Test.Test(t, "engine: expire trial on non-trial returns empty", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.ExpireTrial("111", "svc1", true)
		t.Ok(state == "")
		t.End()
	})
}

func TestExpireTrialUnknown(t *testing.T) {
	Test.Test(t, "engine: expire trial for unknown returns empty", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		state := eng.ExpireTrial("nobody", "svc1", true)
		t.Ok(state == "")
		t.End()
	})
}

func TestRetrySuccess(t *testing.T) {
	Test.Test(t, "engine: retry success sets active", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		eng.OnChargeResult("111", "svc1", engine.ResultPermanent)
		state := eng.Retry("111", "svc1", true)
		t.Ok(state == engine.StateActive)
		t.End()
	})
}

func TestRetryFail(t *testing.T) {
	Test.Test(t, "engine: retry fail moves to removed", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		eng.OnChargeResult("111", "svc1", engine.ResultPermanent)
		state := eng.Retry("111", "svc1", false)
		t.Ok(state == engine.StateRemoved)
		t.End()
	})
}

func TestRetryNotSuspended(t *testing.T) {
	Test.Test(t, "engine: retry on non-suspended returns empty", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.Retry("111", "svc1", true)
		t.Ok(state == "")
		t.End()
	})
}

func TestRetryUnknown(t *testing.T) {
	Test.Test(t, "engine: retry for unknown returns empty", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		state := eng.Retry("nobody", "svc1", true)
		t.Ok(state == "")
		t.End()
	})
}

func TestKickOut(t *testing.T) {
	Test.Test(t, "engine: kick_out moves suspended to removed", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		eng.OnChargeResult("111", "svc1", engine.ResultPermanent)
		state := eng.KickOut("111", "svc1")
		t.Ok(state == engine.StateRemoved)
		t.End()
	})
}

func TestKickOutNotSuspended(t *testing.T) {
	Test.Test(t, "engine: kick_out on non-suspended returns empty", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		eng.Subscribe("111", "svc1", false)
		state := eng.KickOut("111", "svc1")
		t.Ok(state == "")
		t.End()
	})
}

func TestKickOutUnknown(t *testing.T) {
	Test.Test(t, "engine: kick_out for unknown returns empty", func(t *Test.T) {
		eng := newEngine(t.TB(), fixed(epoch))
		state := eng.KickOut("nobody", "svc1")
		t.Ok(state == "")
		t.End()
	})
}
