//go:build ignore

package fixture

import (
	"testing"

	tape "github.com/coderaiser/go-subscriber"
)

func TestParser(t *testing.T) {
	tape.Only(t, "parser: run action", func(t *tape.T) {
		t.Ok(true)
		t.End()
	})
}

func TestRunner(t *testing.T) {
	tape.Only(t, "runner: starts", func(t *tape.T) {
		t.Ok(true)
		t.End()
	})
}
