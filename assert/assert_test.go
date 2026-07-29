package assert

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
)

// mockT implements TB without calling runtime.Goexit().
// Used to test failure paths without failing the outer test.
type mockT struct {
	failed  bool
	message string
}

func (m *mockT) Helper()                      {}
func (m *mockT) Errorf(f string, args ...any) { m.failed = true; m.message = fmt.Sprintf(f, args...) }
func (m *mockT) Fatalf(f string, args ...any) { m.failed = true; m.message = fmt.Sprintf(f, args...) }
func (m *mockT) Fatal(args ...any)            { m.failed = true }

// -- Equal --

func TestEqualMatch(t *testing.T) {
	m := &mockT{}
	Equal(m, 42, 42)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestEqualMismatch(t *testing.T) {
	m := &mockT{}
	Equal(m, 1, 2)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- NoError --

func TestNoErrorNil(t *testing.T) {
	m := &mockT{}
	NoError(m, nil)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestNoErrorNonNil(t *testing.T) {
	m := &mockT{}
	NoError(m, errors.New("oops"))
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- Error --

func TestErrorNonNil(t *testing.T) {
	m := &mockT{}
	Error(m, errors.New("oops"))
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestErrorNil(t *testing.T) {
	m := &mockT{}
	Error(m, nil)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- Ok --

func TestOkTrue(t *testing.T) {
	m := &mockT{}
	Ok(m, true)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestOkFalse(t *testing.T) {
	m := &mockT{}
	Ok(m, false)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

func TestOkIntZero(t *testing.T) {
	m := &mockT{}
	Ok(m, 0)
	if !m.failed {
		t.Fatal("expected fail for zero")
	}
}

func TestOkNonZero(t *testing.T) {
	m := &mockT{}
	Ok(m, 42)
	if m.failed {
		t.Fatal("expected pass for non-zero")
	}
}

// -- NotOk --

func TestNotOkFalse(t *testing.T) {
	m := &mockT{}
	NotOk(m, false)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestNotOkTrue(t *testing.T) {
	m := &mockT{}
	NotOk(m, true)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

func TestNotOkNil(t *testing.T) {
	m := &mockT{}
	NotOk(m, nil)
	if m.failed {
		t.Fatal("expected pass for nil")
	}
}

// -- Contains --

func TestContainsMatch(t *testing.T) {
	m := &mockT{}
	Contains(m, "hello world", "world")
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestContainsNoMatch(t *testing.T) {
	m := &mockT{}
	Contains(m, "hello", "xyz")
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- NotContains --

func TestNotContainsMatch(t *testing.T) {
	m := &mockT{}
	NotContains(m, "hello", "xyz")
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestNotContainsNoMatch(t *testing.T) {
	m := &mockT{}
	NotContains(m, "hello world", "world")
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- NotEqual --

func TestNotEqualMatch(t *testing.T) {
	m := &mockT{}
	NotEqual(m, 1, 2)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestNotEqualFail(t *testing.T) {
	m := &mockT{}
	NotEqual(m, 1, 1)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- DeepEqual --

func TestDeepEqualMatch(t *testing.T) {
	m := &mockT{}
	DeepEqual(m, []int{1}, []int{1})
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestDeepEqualFail(t *testing.T) {
	m := &mockT{}
	DeepEqual(m, []int{1}, []int{2})
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- NotDeepEqual --

func TestNotDeepEqualMatch(t *testing.T) {
	m := &mockT{}
	NotDeepEqual(m, []int{1}, []int{2})
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestNotDeepEqualFail(t *testing.T) {
	m := &mockT{}
	NotDeepEqual(m, []int{1}, []int{1})
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- Match --

func TestMatchSuccess(t *testing.T) {
	m := &mockT{}
	Match(m, "hello 123", `hello \d+`)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestMatchFail(t *testing.T) {
	m := &mockT{}
	Match(m, "hello", "world")
	if !m.failed {
		t.Fatal("expected fail")
	}
}

func TestMatchInvalidPattern(t *testing.T) {
	m := &mockT{}
	Match(m, "hello", `[invalid`)
	if !m.failed {
		t.Fatal("expected fail for invalid pattern")
	}
}

// -- NotMatch --

func TestNotMatchSuccess(t *testing.T) {
	m := &mockT{}
	NotMatch(m, "hello", "world")
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestNotMatchFail(t *testing.T) {
	m := &mockT{}
	NotMatch(m, "hello", "hello")
	if !m.failed {
		t.Fatal("expected fail")
	}
}

func TestNotMatchInvalidPattern(t *testing.T) {
	m := &mockT{}
	NotMatch(m, "hello", `[invalid`)
	if !m.failed {
		t.Fatal("expected fail for invalid pattern")
	}
}

// -- Pass --

func TestPass(t *testing.T) {
	m := &mockT{}
	Pass(m, "message")
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- Fail --

func TestFail(t *testing.T) {
	m := &mockT{}
	Fail(m, "message")
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- isPrimitive --

func TestIsPrimitiveTrue(t *testing.T) {
	if !isPrimitive(42) {
		t.Fatal("expected primitive")
	}
}

func TestIsPrimitivePointer(t *testing.T) {
	x := 42
	if !isPrimitive(&x) {
		t.Fatal("expected pointer as primitive")
	}
}

func TestIsPrimitiveSlice(t *testing.T) {
	if isPrimitive([]int{1}) {
		t.Fatal("expected non-primitive")
	}
}

func TestIsPrimitiveNil(t *testing.T) {
	if isPrimitive(nil) {
		t.Fatal("expected non-primitive for nil")
	}
}

// -- truthy -- 

func TestTruthyBoolTrue(t *testing.T) {
	if !truthy(true) {
		t.Fatal("expected truthy")
	}
}

func TestTruthyBoolFalse(t *testing.T) {
	if truthy(false) {
		t.Fatal("expected falsy")
	}
}

func TestTruthyNil(t *testing.T) {
	if truthy(nil) {
		t.Fatal("expected falsy for nil")
	}
}

func TestTruthyIntZero(t *testing.T) {
	if truthy(0) {
		t.Fatal("expected falsy for zero")
	}
}

func TestTruthyStringEmpty(t *testing.T) {
	if truthy("") {
		t.Fatal("expected falsy for empty string")
	}
}

func TestTruthyDefault(t *testing.T) {
	if !truthy(struct{}{}) {
		t.Fatal("expected truthy for struct")
	}
}

// -- toRegexp --

func TestToRegexpFromString(t *testing.T) {
	re, err := toRegexp(`\d+`)
	if err != nil {
		t.Fatal(err)
	}
	if !re.MatchString("123") {
		t.Fatal("expected match")
	}
}

func TestToRegexpFromRegexp(t *testing.T) {
	re, err := toRegexp(regexp.MustCompile(`\d+`))
	if err != nil {
		t.Fatal(err)
	}
	if !re.MatchString("123") {
		t.Fatal("expected match")
	}
}

func TestToRegexpInvalid(t *testing.T) {
	_, err := toRegexp(`[invalid`)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestToRegexpInvalidType(t *testing.T) {
	_, err := toRegexp(42)
	if err == nil {
		t.Fatal("expected error")
	}
}

// -- ToRegexp (exported wrapper) --

func TestToRegexpExported(t *testing.T) {
	re, err := ToRegexp(`\d+`)
	if err != nil {
		t.Fatal(err)
	}
	if !re.MatchString("123") {
		t.Fatal("expected match")
	}
}

// -- IsPrimitive (exported wrapper) --

func TestIsPrimitiveExported(t *testing.T) {
	if !IsPrimitive(42) {
		t.Fatal("expected primitive")
	}
}

// -- HitCheck --

func TestHitCheckPass(t *testing.T) {
	m := &mockT{}
	HitCheck(m, 1)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestHitCheckFail(t *testing.T) {
	m := &mockT{}
	HitCheck(m, 2)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

// -- CheckScopeName --

func TestCheckScopeNamePass(t *testing.T) {
	m := &mockT{}
	CheckScopeName(m, "scope: message", true, true)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestCheckScopeNameFail(t *testing.T) {
	m := &mockT{}
	CheckScopeName(m, "bad", false, true)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

func TestCheckScopeNameDisabled(t *testing.T) {
	m := &mockT{}
	CheckScopeName(m, "bad", false, false)
	if m.failed {
		t.Fatal("expected pass when check disabled")
	}
}

// -- CheckEndCalled --

func TestCheckEndCalledPass(t *testing.T) {
	m := &mockT{}
	CheckEndCalled(m, true, true)
	if m.failed {
		t.Fatal("expected pass")
	}
}

func TestCheckEndCalledFail(t *testing.T) {
	m := &mockT{}
	CheckEndCalled(m, false, true)
	if !m.failed {
		t.Fatal("expected fail")
	}
}

func TestCheckEndCalledDisabled(t *testing.T) {
	m := &mockT{}
	CheckEndCalled(m, false, false)
	if m.failed {
		t.Fatal("expected pass when check disabled")
	}
}
