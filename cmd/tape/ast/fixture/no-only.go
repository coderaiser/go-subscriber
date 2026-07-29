//go:build ignore

package fixture

import (
	"testing"
	tape "github.com/coderaiser/go-subscriber"
)

func TestNoOnlyParser(t *testing.T) {
	tape.Test(t, "parser: run action", func(t *tape.T) { t.Ok(true); t.End() })
	tape.Test(t, "parser: fail action", func(t *tape.T) { t.Ok(true); t.End() })
}
