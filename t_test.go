package tape

import (
	"os"
	"testing"
)

func TestTB(t *testing.T) {
	Test(t, "tape: TB returns underlying testing.T", func(t *T) {
		t.Ok(t.TB() != nil)
		t.End()
	})
}

func TestSetenv(t *testing.T) {
	Test(t, "tape: Setenv sets env variable", func(t *T) {
		t.Setenv("TAPE_TEST_VAR", "hello")
		t.Equal(os.Getenv("TAPE_TEST_VAR"), "hello")
		t.End()
	})
}
