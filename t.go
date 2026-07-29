package tape

import (
	"regexp"
	"testing"

	"github.com/coderaiser/go-subscriber/assert"
)

// T wraps *testing.T with supertape-compatible assertions.
type T struct {
	t     *testing.T
	ended bool
}

func newT(t *testing.T) *T { return &T{t: t} }

// TB returns the underlying *testing.T for use with helpers
// that require it directly (e.g. t.TB().TempDir(), t.TB().Fatal(...)).
func (tt *T) TB() *testing.T {
	return tt.t
}

// Setenv sets an environment variable for the duration of the test.
// It is automatically restored when the test ends.
func (tt *T) Setenv(key, value string) {
	tt.t.Helper()
	tt.t.Setenv(key, value)
}

// Equal asserts result == expected using strict equality.
// For primitives and pointers only.
// Use DeepEqual for structs, slices, and maps.
func (tt *T) Equal(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	assert.Equal(tt.t, result, expected)
}

// NotEqual asserts result != expected.
// For primitives and pointers only.
func (tt *T) NotEqual(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	assert.NotEqual(tt.t, result, expected)
}

// DeepEqual asserts deep equality using reflect.DeepEqual.
// Use for structs, slices, and maps.
func (tt *T) DeepEqual(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	assert.DeepEqual(tt.t, result, expected)
}

// NotDeepEqual asserts values are not deeply equal.
func (tt *T) NotDeepEqual(result, expected any) {
	tt.t.Helper()
	hit(tt.t)
	assert.NotDeepEqual(tt.t, result, expected)
}

// Ok asserts result is truthy.
func (tt *T) Ok(result any) {
	tt.t.Helper()
	hit(tt.t)
	assert.Ok(tt.t, result)
}

// NotOk asserts result is falsy.
func (tt *T) NotOk(result any) {
	tt.t.Helper()
	hit(tt.t)
	assert.NotOk(tt.t, result)
}

// Match asserts result matches pattern.
// pattern may be a string or *regexp.Regexp.
func (tt *T) Match(result string, pattern any) {
	tt.t.Helper()
	hit(tt.t)
	assert.Match(tt.t, result, pattern)
}

// NotMatch asserts result does not match pattern.
// pattern may be a string or *regexp.Regexp.
func (tt *T) NotMatch(result string, pattern any) {
	tt.t.Helper()
	hit(tt.t)
	assert.NotMatch(tt.t, result, pattern)
}

// Pass generates an unconditional passing assertion.
func (tt *T) Pass(message ...string) {
	tt.t.Helper()
	hit(tt.t)
	msg := "(unnamed assert)"
	if len(message) > 0 {
		msg = message[0]
	}
	tt.t.Log("pass:", msg)
}

// Fail generates an unconditional failing assertion.
// message may be a string or error.
func (tt *T) Fail(message any) {
	tt.t.Helper()
	hit(tt.t)
	switch msg := message.(type) {
	case string:
		tt.t.Errorf("fail: %s", msg)
	case error:
		tt.t.Errorf("fail: %v", msg)
	default:
		tt.t.Errorf("fail: %v", message)
	}
}

// Comment prints a TAP comment without counting as an assertion.
func (tt *T) Comment(message string) {
	tt.t.Helper()
	tt.t.Log("#", message)
}

// Error asserts err is non-nil.
func (tt *T) Error(err error) {
	tt.t.Helper()
	hit(tt.t)
	assert.Error(tt.t, err)
}

// NoError asserts err is nil.
func (tt *T) NoError(err error) {
	tt.t.Helper()
	hit(tt.t)
	assert.NoError(tt.t, err)
}

// Contains asserts s contains sub.
func (tt *T) Contains(s, sub string) {
	tt.t.Helper()
	hit(tt.t)
	assert.Contains(tt.t, s, sub)
}

// NotContains asserts s does not contain sub.
func (tt *T) NotContains(s, sub string) {
	tt.t.Helper()
	hit(tt.t)
	assert.NotContains(tt.t, s, sub)
}

// End marks the test as intentionally complete.
// Required when TAPE_CHECK_END is enabled (default: on).
func (tt *T) End() {
	tt.ended = true
}

// toRegexp is kept for backward compatibility with tests.
func toRegexp(pattern any) (*regexp.Regexp, error) {
	return assert.ToRegexp(pattern)
}

// isPrimitive is kept for backward compatibility with tests.
func isPrimitive(v any) bool {
	return assert.IsPrimitive(v)
}
