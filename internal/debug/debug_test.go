package debug_test

import (
	"testing"

	Test "github.com/coderaiser/go-subscriber/internal/tape"
	"github.com/coderaiser/go-subscriber/internal/debug"
)

func TestLoggerDisabledByDefault(t *testing.T) {
	Test.Test(t, "debug: logger is noop when DEBUG not set", func(t *Test.T) {
		t.TB().Setenv("DEBUG", "")
		log := debug.Logger("subscriber:engine")
		t.Ok(log != nil)
		t.End()
	})
}

func TestLoggerEnabledWildcard(t *testing.T) {
	Test.Test(t, "debug: logger is active for subscriber:*", func(t *Test.T) {
		t.TB().Setenv("DEBUG", "subscriber:*")
		log := debug.Logger("subscriber:engine")
		t.Ok(log != nil)
		t.End()
	})
}

func TestLoggerEnabledExact(t *testing.T) {
	Test.Test(t, "debug: logger is active for exact namespace match", func(t *Test.T) {
		t.TB().Setenv("DEBUG", "subscriber:statemachine")
		log := debug.Logger("subscriber:statemachine")
		t.Ok(log != nil)
		t.End()
	})
}

func TestLoggerDisabledOtherNamespace(t *testing.T) {
	Test.Test(t, "debug: logger is noop for non-matching namespace", func(t *Test.T) {
		t.TB().Setenv("DEBUG", "subscriber:engine")
		log := debug.Logger("subscriber:statemachine")
		t.Ok(log != nil)
		t.End()
	})
}

func TestLoggerJSONFormat(t *testing.T) {
	Test.Test(t, "debug: logger uses JSON handler when LOG_FORMAT=json", func(t *Test.T) {
		t.TB().Setenv("DEBUG", "subscriber:*")
		t.TB().Setenv("LOG_FORMAT", "json")
		log := debug.Logger("subscriber:engine")
		t.Ok(log != nil)
		t.End()
	})
}

func TestLoggerTextFormat(t *testing.T) {
	Test.Test(t, "debug: logger uses text handler when LOG_FORMAT=text", func(t *Test.T) {
		t.TB().Setenv("DEBUG", "subscriber:*")
		t.TB().Setenv("LOG_FORMAT", "text")
		log := debug.Logger("subscriber:engine")
		t.Ok(log != nil)
		t.End()
	})
}

func TestLoggerGlobalWildcard(t *testing.T) {
	Test.Test(t, "debug: logger is active for bare * wildcard", func(t *Test.T) {
		t.TB().Setenv("DEBUG", "*")
		log := debug.Logger("subscriber:engine")
		t.Ok(log != nil)
		t.End()
	})
}
