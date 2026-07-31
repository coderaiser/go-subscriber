// Package tape provides a minimal test wrapper for group-style table tests.
package tape

import (
	"fmt"
	"strings"
	"testing"
)

// T wraps *testing.T to provide fluent assertion methods.
type T struct {
	parent *testing.T
}

// TB returns the underlying *testing.T.
func (t *T) TB() *testing.T { return t.parent }

// Ok fails if cond is false.
func (t *T) Ok(cond bool) {
	t.Helper()
	if !cond {
		t.parent.Fatal("expected true, got false")
	}
}

// NotOk fails if cond is true.
func (t *T) NotOk(cond bool) {
	t.Helper()
	if cond {
		t.parent.Fatal("expected false, got true")
	}
}

// NoError fails if err is not nil.
func (t *T) NoError(err error) {
	t.Helper()
	if err != nil {
		t.parent.Fatalf("expected no error, got: %v", err)
	}
}

// Error fails if err is nil.
func (t *T) Error(err error) {
	t.Helper()
	if err == nil {
		t.parent.Fatal("expected error, got nil")
	}
}

// Equal fails if got != want.
func (t *T) Equal(got, want any) {
	t.Helper()
	if got != want {
		t.parent.Fatalf("expected %v, got %v", want, got)
	}
}

// NotEqual fails if got == want.
func (t *T) NotEqual(got, want any) {
	t.Helper()
	if got == want {
		t.parent.Fatalf("expected not %v", want)
	}
}

// DeepEqual fails if values are not deeply equal.
func (t *T) DeepEqual(got, want any) {
	t.Helper()
	if !deepEqual(got, want) {
		t.parent.Fatalf("expected %#v, got %#v", want, got)
	}
}

// NotDeepEqual fails if values are deeply equal.
func (t *T) NotDeepEqual(got, want any) {
	t.Helper()
	if deepEqual(got, want) {
		t.parent.Fatalf("expected not %#v", want)
	}
}

// Match fails if s does not contain substr (simple substring match,
// not regex — kept for compatibility with the assertion signature).
func (t *T) Match(s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.parent.Fatalf("expected %q to contain %q", s, substr)
	}
}

// NotMatch fails if s contains substr.
func (t *T) NotMatch(s, substr string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.parent.Fatalf("expected %q to NOT contain %q", s, substr)
	}
}



// Comment is a no-op that prints a message.
func (t *T) Comment(msg string) {
	t.parent.Log(msg)
}

// End marks the end of a subtest. No-op for compatibility.
func (t *T) End() {}

func (t *T) Helper() { t.parent.Helper() }

// deepEqual is a simple deep equality check for basic types and slices.
func deepEqual(a, b any) bool {
	// Use fmt.Sprintf for a simple deep comparison
	return fmt.Sprintf("%#v", a) == fmt.Sprintf("%#v", b)
}

// Test runs a named subtest with a tape.T wrapper.
func Test(parent *testing.T, name string, fn func(t *T)) {
	parent.Run(name, func(st *testing.T) {
		t := &T{parent: st}
		fn(t)
	})
}
