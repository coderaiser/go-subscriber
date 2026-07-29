//go:build coverage

package tape

import "testing"

// TestScopeCheckFailsInsideTest verifies that Test() Fatalfs on invalid scope.
// Only built with -tags coverage. Intentionally fails.
func TestScopeCheckFailsInsideTest(t *testing.T) {
	t.Setenv("TAPE_CHECK_SCOPES", "1")
	failed := false
	t.Run("outer", func(outer *testing.T) {
		defer func() { failed = outer.Failed() }()
		Test(outer, "no scope here", func(t *T) {
			t.Ok(true)
			t.End()
		})
	})
	if !failed {
		t.Fatal("expected Test() to fail on invalid scope name")
	}
}

// TestHitTwiceWithCountEnabled verifies that hit() Fatalfs on second call.
// Only built with -tags coverage. Intentionally fails.
func TestHitTwiceWithCountEnabled(t *testing.T) {
	failed := false
	t.Run("inner", func(inner *testing.T) {
		defer func() { failed = inner.Failed() }()
		assertOne(inner)
		hit(inner) // first — ok
		hit(inner) // second — Fatalf
	})
	if !failed {
		t.Fatal("expected second hit to fail the test")
	}
}
