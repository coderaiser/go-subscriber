package engine

import "github.com/coderaiser/go-subscriber/internal/statemachine"

// transitions is the complete state machine definition for the subscription
// lifecycle. Every row maps directly to a state and event from the spec.
var transitions = &statemachine.MemorySource{
	Defs: []statemachine.TransitionDef{
		// trial
		{From: StateTrial, Event: EventExpireSuccess, To: StateActive},
		{From: StateTrial, Event: EventExpireFail,    To: StateSuspended},
		{From: StateTrial, Event: EventUnsubscribe,   To: StateTerminated},

		// active -- subscribe_paid and renew_success both land here
		{From: StateActive, Event: EventRenewSuccess, To: StateActive},
		{From: StateActive, Event: EventRenewFail,    To: StateSuspended},
		{From: StateActive, Event: EventUnsubscribe,  To: StateTerminated},

		// suspended
		{From: StateSuspended, Event: EventRetrySuccess, To: StateActive},
		{From: StateSuspended, Event: EventKickOut,      To: StateRemoved},
		{From: StateSuspended, Event: EventUnsubscribe,  To: StateTerminated},
	},
}
