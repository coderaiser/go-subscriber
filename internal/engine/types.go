package engine

// State is the lifecycle state of a (msisdn, serviceID) subscription.
type State = string

const (
	StateTrial      = "trial"
	StateActive     = "active"
	StateSuspended  = "suspended"
	StateTerminated = "terminated"
	StateRemoved    = "removed"
)

// Event is a named trigger fired into the state machine.
type Event = string

const (
	EventSubscribeTrial = "subscribe_trial"
	EventSubscribePaid  = "subscribe_paid"
	EventExpireSuccess  = "expire_success"
	EventExpireFail     = "expire_fail"
	EventRenewSuccess   = "renew_success"
	EventRenewFail      = "renew_fail"
	EventRetrySuccess   = "retry_success"
	EventKickOut        = "kick_out"
	EventUnsubscribe    = "unsubscribe"
)

// ChargeResult is the answer the carrier sends back for a charge attempt.
type ChargeResult = string

const (
	ResultSuccess         = "success"
	ResultLowBalance      = "low_balance"
	ResultPermanent       = "permanent"
	ResultRateLimit       = "rate_limit"
	ResultPending         = "pending"
	ResultSystemError     = "system_error"
	ResultNoResponse      = "no_response"
	ResultSubscriberState = "subscriber_state"
)

// CooloffDays is the number of days before a terminated/removed subscriber
// can re-subscribe.
const CooloffDays = 30

// FallbackThreshold is the number of low_balance results before suspending.
const FallbackThreshold = 3
