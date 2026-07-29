package statemachine

// Adapter stores current state per entity id.
//
// Get returns:
//   - nil, nil  → id never set
//   - &v, nil   → state found
//   - nil, err  → storage failure
type Adapter[S any] interface {
	Get(id string) (*S, error)
	Set(id string, state S) error
}
