package config

import (
	"os"
	"testing"
	"time"
)

func TestCheckScopesDefault(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_CHECK_SCOPES") })
	os.Unsetenv("TAPE_CHECK_SCOPES")
	if !CheckScopes() {
		t.Error("default should be true")
	}
}

func TestCheckScopesDisabled(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_CHECK_SCOPES") })
	os.Setenv("TAPE_CHECK_SCOPES", "0")
	if CheckScopes() {
		t.Error("should be false when env is 0")
	}
}

func TestCheckAssertionsCountDefault(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT") })
	os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	if !CheckAssertionsCount() {
		t.Error("default should be true")
	}
}

func TestCheckAssertionsCountDisabled(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT") })
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "false")
	if CheckAssertionsCount() {
		t.Error("should be false when env is false")
	}
}

func TestCheckEndDefault(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_CHECK_END") })
	os.Unsetenv("TAPE_CHECK_END")
	if !CheckEnd() {
		t.Error("default should be true")
	}
}

func TestCheckEndDisabled(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_CHECK_END") })
	os.Setenv("TAPE_CHECK_END", "0")
	if CheckEnd() {
		t.Error("should be false when env is 0")
	}
}

func TestTimeoutDefault(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_TIMEOUT") })
	os.Unsetenv("TAPE_TIMEOUT")
	if d := Timeout(); d != 3*time.Second {
		t.Errorf("default should be 3s, got %v", d)
	}
}

func TestTimeoutCustom(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_TIMEOUT") })
	os.Setenv("TAPE_TIMEOUT", "5s")
	if d := Timeout(); d != 5*time.Second {
		t.Errorf("should be 5s, got %v", d)
	}
}

func TestTimeoutInvalid(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_TIMEOUT") })
	os.Setenv("TAPE_TIMEOUT", "invalid")
	if d := Timeout(); d != 3*time.Second {
		t.Errorf("invalid should fallback to 3s, got %v", d)
	}
}

func TestStrictTransitionsDefault(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_STRICT_TRANSITIONS") })
	os.Unsetenv("TAPE_STRICT_TRANSITIONS")
	if !StrictTransitions() {
		t.Error("default should be true")
	}
}

func TestStrictTransitionsDisabled(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_STRICT_TRANSITIONS") })
	os.Setenv("TAPE_STRICT_TRANSITIONS", "0")
	if StrictTransitions() {
		t.Error("should be false when env is 0")
	}
}

func TestEnvBoolTrue(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_CHECK_SCOPES") })
	os.Setenv("TAPE_CHECK_SCOPES", "true")
	if !CheckScopes() {
		t.Error("expected true for 'true'")
	}
}

func TestEnvBoolInvalid(t *testing.T) {
	t.Cleanup(func() { os.Unsetenv("TAPE_CHECK_SCOPES") })
	os.Setenv("TAPE_CHECK_SCOPES", "invalid")
	if !CheckScopes() {
		t.Error("expected default true for invalid value")
	}
}
