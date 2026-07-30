package main

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestVersionFromJSONValid(t *testing.T) {
	tape.Test(t, "version: VersionFromJSON returns version string", func(t *tape.T) {
		result := VersionFromJSON([]byte(`{"version":"1.2.3"}`))
		t.Equal(result, "1.2.3")
		t.End()
	})
}

func TestVersionFromJSONInvalid(t *testing.T) {
	tape.Test(t, "version: VersionFromJSON returns unknown on invalid JSON", func(t *tape.T) {
		result := VersionFromJSON([]byte(`{invalid`))
		t.Equal(result, "unknown")
		t.End()
	})
}

func TestVersionFromJSONEmpty(t *testing.T) {
	tape.Test(t, "version: VersionFromJSON returns unknown on empty version", func(t *tape.T) {
		result := VersionFromJSON([]byte(`{"version":""}`))
		t.Equal(result, "unknown")
		t.End()
	})
}

func TestVersionLineContainsV(t *testing.T) {
	tape.Test(t, "version: VersionLine starts with v", func(t *tape.T) {
		result := VersionLine()
		t.Match(result, "v")
		t.End()
	})
}

func TestVersionReturnsString(t *testing.T) {
	tape.Test(t, "version: Version returns non-empty string", func(t *tape.T) {
		result := Version()
		t.Ok(result != "")
		t.End()
	})
}
