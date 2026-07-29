//go:build ignore

package fixture

import (
	"testing"
	tape "github.com/coderaiser/go-subscriber"
)

func TestOneOnlyParser(t *testing.T) {
	tape.Only(t, "parser: run action", func(t *tape.T) { t.Ok(true); t.End() })
	tape.Test(t, "parser: other test", func(t *tape.T) { t.Ok(true); t.End() })
}
