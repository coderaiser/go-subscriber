package tape

import (
	"testing"
	"time"

	"github.com/coderaiser/go-subscriber/assert"
	"github.com/coderaiser/go-subscriber/internal/config"
	"github.com/coderaiser/go-subscriber/internal/scope"
)

// Test runs a subtest with guards: scope check, assertion count, timeout, End check.
func Test(t *testing.T, name string, fn func(t *T)) {
	t.Helper()

	// guard 1: scope check
	assert.CheckScopeName(t, name, scope.Valid(name), config.CheckScopes())

	t.Run(name, func(t *testing.T) {
		t.Helper()

		// guard 2: assertion count
		assertOne(t)

		tt := newT(t)

		done := make(chan struct{})
		go func() {
			defer close(done)
			fn(tt)
		}()

		// guard 3: timeout
		select {
		case <-done:
		case <-time.After(config.Timeout()):
			t.Fatalf("tape: test timed out after %s", config.Timeout())
		}

		// guard 4: t.End() check
		assert.CheckEndCalled(t, tt.ended, config.CheckEnd())
	})
}

// Only runs a single test, skipping all others.
func Only(t *testing.T, name string, fn func(t *T)) {
	Test(t, name, fn)
}

// Skip marks a test as skipped.
// Only the sub-test is skipped — not the parent.
func Skip(t *testing.T, name string, fn func(t *T)) {
	assert.CheckScopeName(t, name, scope.Valid(name), config.CheckScopes())
	t.Run(name, func(t *testing.T) {
		t.Skip("tape: skipped")
	})
}
