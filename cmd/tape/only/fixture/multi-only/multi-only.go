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
	tape.Only(t, "parser: fail action", func(t *tape.T) {
		t.Ok(true)
		t.End()
	})
}
