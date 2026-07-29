package assert

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// TB is the subset of testing.TB used by assert functions.
// Allows mockT injection in tests.
type TB interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Fatal(args ...any)
}

// Equal asserts got == want using reflect.DeepEqual.
func Equal(t TB, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:  %#v\nwant: %#v", got, want)
	}
}

// NotEqual asserts got != want for primitives and pointers only.
func NotEqual(t TB, got, want any) {
	t.Helper()
	if isPrimitive(got) && reflect.DeepEqual(got, want) {
		t.Errorf("expected values to differ, both are: %#v", got)
	}
}

// DeepEqual asserts deep equality.
func DeepEqual(t TB, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:  %#v\nwant: %#v", got, want)
	}
}

// NotDeepEqual asserts values are not deeply equal.
func NotDeepEqual(t TB, got, want any) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		t.Errorf("expected values to differ, both deep-equal: %#v", got)
	}
}

// Ok asserts result is truthy.
func Ok(t TB, result any) {
	t.Helper()
	if !truthy(result) {
		t.Errorf("expected truthy, got: %#v", result)
	}
}

// NotOk asserts result is falsy.
func NotOk(t TB, result any) {
	t.Helper()
	if truthy(result) {
		t.Errorf("expected falsy, got: %#v", result)
	}
}

// Match asserts result matches pattern.
func Match(t TB, result string, pattern any) {
	t.Helper()
	re, err := toRegexp(pattern)
	if err != nil {
		t.Errorf("Match: invalid pattern: %v", err)
		return
	}
	if !re.MatchString(result) {
		t.Errorf("%q does not match %q", result, re)
	}
}

// NotMatch asserts result does not match pattern.
func NotMatch(t TB, result string, pattern any) {
	t.Helper()
	re, err := toRegexp(pattern)
	if err != nil {
		t.Errorf("NotMatch: invalid pattern: %v", err)
		return
	}
	if re.MatchString(result) {
		t.Errorf("%q should not match %q", result, re)
	}
}

// Pass logs a message.
func Pass(t TB, message string) {
	t.Helper()
	t.Errorf("pass: %s", message)
}

// Fail logs a failure.
func Fail(t TB, message string) {
	t.Helper()
	t.Errorf("fail: %s", message)
}

// NoError asserts err is nil.
func NoError(t TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Error asserts err is non-nil.
func Error(t TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

// Contains asserts s contains sub.
func Contains(t TB, s, sub string) {
	t.Helper()
	if !strings.Contains(s, sub) {
		t.Errorf("%q does not contain %q", s, sub)
	}
}

// NotContains asserts s does not contain sub.
func NotContains(t TB, s, sub string) {
	t.Helper()
	if strings.Contains(s, sub) {
		t.Errorf("%q should not contain %q", s, sub)
	}
}

// CheckScopeName checks a test name follows "scope: message" format.
// If checkScopes is true and valid is false, it calls t.Fatalf.
func CheckScopeName(t TB, name string, valid, checkScopes bool) {
	t.Helper()
	if checkScopes && !valid {
		t.Fatalf("tape: test name must be 'scope: message', got: %q", name)
	}
}

// CheckEndCalled checks that t.End() was called.
// If checkEnd is true and ended is false, it calls t.Fatal.
func CheckEndCalled(t TB, ended, checkEnd bool) {
	t.Helper()
	if checkEnd && !ended {
		t.Fatal("tape: missing t.End()")
	}
}

func HitCheck(t TB, count int) {
	t.Helper()
	if count > 1 {
		t.Fatalf("too many assertions: got %d, expected 1", count)
	}
}

// ToRegexp is exported for use by the tape package.
func ToRegexp(pattern any) (*regexp.Regexp, error) {
	return toRegexp(pattern)
}

// IsPrimitive is exported for use by the tape package.
func IsPrimitive(v any) bool {
	return isPrimitive(v)
}

// truthy checks if v is truthy.
func truthy(v any) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case int:
		return val != 0
	case string:
		return val != ""
	default:
		return true
	}
}

// toRegexp converts a string or *regexp.Regexp to *regexp.Regexp.
func toRegexp(pattern any) (*regexp.Regexp, error) {
	switch p := pattern.(type) {
	case *regexp.Regexp:
		return p, nil
	case string:
		return regexp.Compile(p)
	default:
		return nil, fmt.Errorf("pattern must be string or *regexp.Regexp, got %T", pattern)
	}
}

// isPrimitive checks if v is a primitive type or pointer.
func isPrimitive(v any) bool {
	switch v.(type) {
	case bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		complex64, complex128,
		string, uintptr:
		return true
	}
	t := reflect.TypeOf(v)
	return t != nil && t.Kind() == reflect.Pointer
}

