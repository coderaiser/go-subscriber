package statemachine

// Context provides state machine context to a handler.
type Context[S, E any] struct {
	From    S
	Event   E
	Payload any
	Adapter Adapter[S]
}

// Handler is a side-effect function for a transition.
type Handler[S, E any] func(ctx Context[S, E]) error
