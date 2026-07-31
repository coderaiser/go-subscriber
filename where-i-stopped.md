# Where I stopped

What I would do next, in priority order.

***

## 1. Complete the subscribe rules gauntlet

The spec lists eight checks that `subscribe()` must run before creating any state: valid msisdn format, antifraud ID single-use, no same transaction pending, max subscriptions per period, daily subscription cap, cooloff, trial-already-used, already-subscribed. The engine currently enforces cooloff, trial-already-used, and already-subscribed. The remaining five (msisdn validation, antifraud, pending transaction guard, period cap, daily cap) are not implemented. These are the rules most likely to matter in a fraud or abuse scenario, which is the whole reason the gauntlet exists.

## 2. Implement `on_delivery_report` (non-DCB path)

The spec defines two charging methods: DCB (carrier billing via Diameter, confirmed by `on_charge_result`) and non-DCB (SMSC charges the MT, confirmed by delivery report). The simulator only implements DCB. `on_delivery_report` is currently absent. A service configured as non-DCB would silently misbehave.

## 3. Enforce retry time windows

The spec says daily service retries run at two specific slots: first in the morning, second before 21:00. The simulator accepts `/retry` calls at any time — the scheduler simulation is manual injection with no time constraint. Adding a simulated clock (injectable `now` function is already in the engine) and checking the slot window on `/retry` would make the simulator faithful to the spec without adding real timers.

## 4. Implement `reconcile()`

Currently a stub. The spec describes it as a batch job that compares the platform's transaction store against operator CDRs and resolves NO_RESPONSE charges. In a simulator this would mean: accept a list of CDR outcomes, find matching transactions marked NO_RESPONSE, and apply the real result. This is the only way NO_RESPONSE charges get resolved per the spec.

## 5. Seed-based eligibility and per-service config

`users.json` seeds initial state but does not enforce per-service eligibility rules or service-level config (cycle type, fallback prices, free trial on/off). All services currently behave identically. A service config store would let the simulator reflect the spec's service definition table and make the demo scenarios more realistic.

## 6. Persistence across restarts

State is in-memory. For a real deployment the StateStore and FactsStore adapters would need a durable backend. The adapter interface in `statemachine` already supports this — a file-system or database adapter would plug in without changing engine or handler code.
