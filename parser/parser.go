package parser

import (
	"encoding/json"
	"fmt"

	"github.com/coderaiser/go-subscriber/model"
)

// Parse parses a go test -json line into an Event.
func Parse(line string) (model.Event, error) {
	var e model.Event
	if err := json.Unmarshal([]byte(line), &e); err != nil {
		return e, fmt.Errorf("parse event: %w", err)
	}
	return e, nil
}
