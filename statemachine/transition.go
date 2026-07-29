package statemachine

import "fmt"

// TransitionDef defines a single state transition.
type TransitionDef struct {
	From  string
	Event string
	To    string
}

// Validate checks a set of transition definitions for basic correctness.
func Validate(defs []TransitionDef) error {
	for _, d := range defs {
		if d.From == "" {
			return fmt.Errorf("transition has empty From")
		}
		if d.Event == "" {
			return fmt.Errorf("transition has empty Event")
		}
		if d.To == "" {
			return fmt.Errorf("transition has empty To")
		}
	}
	return nil
}
