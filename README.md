# Subscriber [![License][LicenseIMGURL]][LicenseURL] [![Build Status][BuildStatusIMGURL]][BuildStatusURL] [![Coverage Status][CoverageIMGURL]][CoverageURL]

[BuildStatusURL]: https://github.com/coderaiser/go-subscriber/actions/workflows/test.yml
[BuildStatusIMGURL]: https://github.com/coderaiser/go-subscriber/actions/workflows/test.yml/badge.svg
[LicenseURL]: https://tldrlegal.com/license/mit-license "MIT License"
[LicenseIMGURL]: https://img.shields.io/badge/license-MIT-317BF9.svg?style=flat
[CoverageURL]: https://coveralls.io/github/coderaiser/go-subscriber?branch=master
[CoverageIMGURL]: https://coveralls.io/repos/coderaiser/go-subscriber/badge.svg?branch=master&service=github

Subscription engine simulator. HTTP service that tracks every `(msisdn, service)`pair through its billing lifecycle: trial → active → suspended → terminated/removed.

Built as a work sample for a Telco SDP assignment.

## Architecture

Four layers. Each only talks to the layer directly below it.

```
┌─────────────────────────────────┐
│         HTTP Server             │  cmd/subscriber/server.go
│  decodes requests, encodes      │  cmd/subscriber/main.go
│  responses, calls engine        │
├─────────────────────────────────┤
│           Engine                │  internal/engine/engine.go
│  business rules, ladder,        │  internal/engine/types.go
│  cooloff, calls statemachine    │  internal/engine/transitions.go
├─────────────────────────────────┤
│        StateMachine             │  internal/statemachine/statemachine.go
│  pure (state, event) → state    │  vendored from coderaiser/go-tape
│  no side effects                │  trimmed to stdlib only
├─────────────────────────────────┤
│        Two Stores               │
│  StateStore  — machine adapter  │  internal/store/state.go
│  FactsStore  — engine only      │  internal/store/facts.go
└─────────────────────────────────┘
      │
  Logger                           internal/debug/debug.go
  DEBUG=subscriber:*
  DEBUG=subscriber:engine
  DEBUG=subscriber:statemachine
  LOG_FORMAT=json|text
```

## Design decisions

**Onion architecture** — each layer has one job and one direction of dependency.
HTTP knows about Engine. Engine knows about StateMachine and Stores. StateMachine
knows nothing about business logic. Changing the HTTP framework touches only
the handler layer. Swapping persistence touches only the store layer.

**Two stores, not one** — `StateStore` is the state machine's adapter: it owns
the current state string per key and satisfies `Adapter[string]`. `FactsStore`is the engine's private store: it owns `CooloffUntil` and `FallbackAttempt`,
which persist across subscriptions and have no place in the machine. Keeping
them separate means the machine stays generic and the engine owns its invariants.

**StateMachine vendored from go-tape** — `internal/statemachine/statemachine.go`is a single-file vendor of the statemachine package from [`go-tape`](https://github.com/coderaiser/go-tape). TOML support removed (external dep). Memory adapter merged into the same file.
No `go.mod` dependency — the code is ours to read and explain in a review.

**Ladder counter in the engine, not the machine** — the spec's `LOW_BALANCE`fallback (3 attempts before suspending) does not map to FSM states. Modelling
it as `fallback_1`/`fallback_2` states would invent states the spec does not
have. Instead the counter lives in `FactsStore` and the engine checks it before
deciding which event to fire. The machine table stays a faithful 9-row map of
the spec. See `DECISIONS.md` for the full reasoning.

## State machine

```
(no state) ── subscribe trial ──────────► trial
(no state) ── subscribe paid ───────────► active
trial      ── expire_success ───────────► active
trial      ── expire_fail ──────────────► suspended
trial      ── unsubscribe ──────────────► terminated
active     ── renew_success (self-loop) ► active
active     ── renew_fail ───────────────► suspended
active     ── unsubscribe ──────────────► terminated
suspended  ── retry_success ────────────► active
suspended  ── kick_out ─────────────────► removed
suspended  ── unsubscribe ──────────────► terminated
terminated/removed ── (cooloff 30d) ───► can re-subscribe
```

## Endpoints

| Method | Path              | Notes                                               |
|--------|-------------------|-----------------------------------------------------|
| `POST` | `/subscribe`      | trial or paid; checks cooloff and duplicate         |
| `POST` | `/unsubscribe`    | any state → terminated; starts 30-day cooloff       |
| `POST` | `/charge-result`  | injects carrier answer; drives ladder + transitions |
| `POST` | `/renew`          | scheduler simulation: renewal tick                  |
| `POST` | `/expire-trial`   | scheduler simulation: trial expiry                  |
| `POST` | `/retry`          | scheduler simulation: retry window                  |
| `POST` | `/kick-out`       | scheduler simulation: retries exhausted             |
| `GET`  | `/state/{msisdn}` | all subscriptions for a msisdn                      |
| `POST` | `/send-mt`        | stub — returns `{"status":"ok","stub":true}`        |
| `GET`  | `/healthz`        | liveness probe                                      |
| `GET`  | `/readyz`         | readiness probe                                     |

## Running locally

```sh
go run ./cmd/subscriber
curl http://localhost:8080/healthz
# with debug logging
DEBUG=subscriber:* go run ./cmd/subscriber
# json logs (production style)
LOG_FORMAT=json PORT=8080 go run ./cmd/subscriber
```

## Flags

```
-v, --version   print version and exit
-h, --help      print this help and exit
```

## Testing

```sh
# run all tests
go test ./...
# with coverage
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out
```

## Docker

```sh
docker build -t go-subscriber .
docker run --rm go-subscriber --version
docker run --rm -p 8080:8080 go-subscriber
```

## Manual testing

Run the demo script (requires `jq`):

```sh
go run ./cmd/subscriber &
bash demo.sh
```

Or open `demo.http` in VS Code (REST Client extension) or JetBrains HTTP Client.

## Release

Push a tag to trigger build → test → coverage check → goreleaser publishes
6-platform binaries to GitHub Releases.

```sh
git tag v1.0.0
git push origin v1.0.0
```

## Future

- `GET /ui` — embedded HTML interface for manual testing
- `GET /swagger` — OpenAPI spec

## Licence

MIT
