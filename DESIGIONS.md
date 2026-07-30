# Design decisions

Ambiguities in the spec resolved before writing code.
Nothing decided by accident.

## 1. Retry count contradiction

Spec says "2 attempts/day" in one place and "1/day, 45 total" elsewhere.
**Decision: 2/day**, 45-total cap as global ceiling.
Scheduler simulation ignores time windows — inject `/retry` manually.

## 2. SUBSCRIBER_STATE carrier result

Spec leaves the reaction as an "open decision".
**Decision: treat as failure** — transitions to `suspended` via `renew_fail`.
Same behaviour as `LOW_BALANCE` ladder exhausted.

## 3. SYSTEM_ERROR consumes retry slot?

Not specified.
**Decision: does NOT consume a slot.** The error is platform-side, not
subscriber-side. Mirrors `RATE_LIMIT` behaviour.

## 4. Re-subscribe after TERMINATED or REMOVED

Both states require cooloff.
**Decision: `FactsStore` row is permanent.** `CooloffUntil` is always checked
on subscribe regardless of current subscription state.

## 5. Non-Hosted services

Spec describes a separate code path.
**Decision: out of scope.** All services in the simulator are Hosted.

## 6. Scheduler not implemented

Real timers add complexity with no benefit for a simulator.
**Decision: time-driven events injected via HTTP** — `/renew`, `/retry`,
`/expire-trial`, `/kick-out` simulate what a real scheduler would trigger.
The machine table stays a faithful 9-row map of the spec's 5 real states.

## 8. engine.New() panics on build failure

`engine.New()` builds the state machine from a hardcoded transition table.
If that fails, something is fundamentally broken — not a runtime condition.
**Decision: panic immediately.** Returning nil would hide a programming error
behind a nil pointer dereference elsewhere.

