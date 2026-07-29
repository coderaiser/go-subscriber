# Subscriber

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
│    owns: current State per key  │
│  FactsStore  — engine only      │  internal/store/facts.go
│    owns: CooloffUntil, counter  │
└─────────────────────────────────┘
      │
  Logger                           internal/debug/debug.go
  DEBUG=subscriber:*
  DEBUG=subscriber:engine
  DEBUG=subscriber:statemachine
  LOG_FORMAT=json|text
```

## LICENCE

MIT
