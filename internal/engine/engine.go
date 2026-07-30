package engine

import (
	"errors"
	"log/slog"
	"time"

	"github.com/coderaiser/go-subscriber/internal/debug"
	"github.com/coderaiser/go-subscriber/internal/statemachine"
	"github.com/coderaiser/go-subscriber/internal/store"
)

var (
	ErrAlreadySubscribed = errors.New("already subscribed")
	ErrCooloff           = errors.New("cooloff period active")
	ErrTrialAlreadyUsed  = errors.New("trial already used")
)

// Engine holds business logic for subscription lifecycle.
type Engine struct {
	states  *store.StateStore
	facts   *store.FactsStore
	now     func() time.Time
	log     *slog.Logger
	machine *statemachine.Machine[string, string]
}

func New(
	states *store.StateStore,
	facts *store.FactsStore,
	now func() time.Time,
	log *slog.Logger,
) *Engine {
	parseState := func(s string) (string, error) {
		switch s {
		case StateTrial, StateActive, StateSuspended, StateTerminated, StateRemoved:
			return s, nil
		}
		return "", errors.New("unknown state: " + s)
	}
	parseEvent := func(e string) (string, error) {
		switch e {
		case EventExpireSuccess, EventExpireFail,
			EventRenewSuccess, EventRenewFail,
			EventRetrySuccess, EventKickOut,
			EventUnsubscribe:
			return e, nil
		}
		return "", errors.New("unknown event: " + e)
	}

	m, err := statemachine.New(transitions, parseState, parseEvent, states, false)
	if err != nil {
		panic("subscriber: failed to build state machine: " + err.Error())
	}
	m.WithLogger(debug.Logger("subscriber:statemachine"))

	return &Engine{
		states:  states,
		facts:   facts,
		now:     now,
		log:     log,
		machine: m,
	}
}

func key(msisdn, serviceID string) string { return msisdn + ":" + serviceID }

func (e *Engine) getState(id string) State {
	ptr, err := e.states.Get(id)
	if err != nil || ptr == nil {
		return ""
	}
	return *ptr
}

// Subscribe creates a new subscription.
func (e *Engine) Subscribe(msisdn, serviceID string, trial bool) error {
	id := key(msisdn, serviceID)
	s := e.getState(id)
	if s == StateActive || s == StateTrial || s == StateSuspended {
		return ErrAlreadySubscribed
	}

	facts := e.facts.Get(id)
	if facts.CooloffUntil > 0 {
		if e.now().Unix() < facts.CooloffUntil {
			return ErrCooloff
		}
	}

	if trial {
		facts := e.facts.Get(id)
		if facts.TrialUsed {
			return ErrTrialAlreadyUsed
		}
		facts.TrialUsed = true
		e.facts.Set(id, facts)
		return e.states.Set(id, StateTrial)
	}
	return e.states.Set(id, StateActive)
}

// Unsubscribe removes a subscription.
func (e *Engine) Unsubscribe(msisdn, serviceID string) State {
	id := key(msisdn, serviceID)
	state := e.getState(id)
	if state == "" {
		return StateTerminated
	}

	facts := e.facts.Get(id)
	facts.CooloffUntil = e.now().AddDate(0, 0, CooloffDays).Unix()
	e.facts.Set(id, facts)

	next, err := e.machine.Apply(id, EventUnsubscribe, nil)
	if err != nil {
		return StateTerminated
	}
	return next
}

// OnChargeResult processes a charge result from the carrier.
func (e *Engine) OnChargeResult(msisdn, serviceID string, result ChargeResult) State {
	id := key(msisdn, serviceID)
	state := e.getState(id)
	if state == "" {
		return ""
	}

	switch result {
	case ResultSuccess:
		facts := e.facts.Get(id)
		facts.FallbackAttempt = 0
		e.facts.Set(id, facts)
		next, _ := e.machine.Apply(id, EventRenewSuccess, nil)
		return next
	case ResultLowBalance:
		facts := e.facts.Get(id)
		facts.FallbackAttempt++
		e.facts.Set(id, facts)
		if facts.FallbackAttempt < FallbackThreshold {
			return state
		}
		facts.FallbackAttempt = 0
		e.facts.Set(id, facts)
		next, _ := e.machine.Apply(id, EventRenewFail, nil)
		return next
	case ResultPermanent, ResultSubscriberState:
		next, _ := e.machine.Apply(id, EventRenewFail, nil)
		return next
	case ResultRateLimit, ResultSystemError, ResultPending, ResultNoResponse:
		return state
	}
	return state
}

// ExpireTrial handles trial expiration.
func (e *Engine) ExpireTrial(msisdn, serviceID string, success bool) State {
	id := key(msisdn, serviceID)
	state := e.getState(id)
	if state != StateTrial {
		return ""
	}

	var event string
	if success {
		event = EventExpireSuccess
	} else {
		event = EventExpireFail
	}
	next, _ := e.machine.Apply(id, event, nil)
	return next
}

// Retry attempts to retry a suspended subscription.
func (e *Engine) Retry(msisdn, serviceID string, success bool) State {
	id := key(msisdn, serviceID)
	state := e.getState(id)
	if state != StateSuspended {
		return ""
	}

	var event string
	if success {
		event = EventRetrySuccess
	} else {
		event = EventKickOut
	}
	next, _ := e.machine.Apply(id, event, nil)
	return next
}

// KickOut removes a suspended subscriber permanently.
func (e *Engine) KickOut(msisdn, serviceID string) State {
	id := key(msisdn, serviceID)
	state := e.getState(id)
	if state != StateSuspended {
		return ""
	}

	next, _ := e.machine.Apply(id, EventKickOut, nil)
	return next
}

// Facts returns the FactsStore for direct access in tests.
func (e *Engine) Facts() *store.FactsStore { return e.facts }

