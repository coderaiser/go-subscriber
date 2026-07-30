package main

import (
	"strings"
	"testing"

	Test "github.com/coderaiser/go-subscriber/internal/tape"
)

func TestVersionFromJSONValid(t *testing.T) {
	Test.Test(t, "version: VersionFromJSON returns version string", func(t *Test.T) {
		result := VersionFromJSON([]byte(`{"version":"1.2.3"}`))
		t.Equal(result, "1.2.3")
		t.End()
	})
}

func TestVersionFromJSONInvalid(t *testing.T) {
	Test.Test(t, "version: VersionFromJSON returns unknown on invalid JSON", func(t *Test.T) {
		result := VersionFromJSON([]byte(`{invalid`))
		t.Equal(result, "unknown")
		t.End()
	})
}

func TestVersionFromJSONEmpty(t *testing.T) {
	Test.Test(t, "version: VersionFromJSON returns unknown on empty version", func(t *Test.T) {
		result := VersionFromJSON([]byte(`{"version":""}`))
		t.Equal(result, "unknown")
		t.End()
	})
}

func TestVersionLineContainsV(t *testing.T) {
	Test.Test(t, "version: VersionLine starts with v", func(t *Test.T) {
		result := VersionLine()
		t.Match(result, "v")
		t.End()
	})
}

func TestVersionReturnsString(t *testing.T) {
	Test.Test(t, "version: Version returns non-empty string", func(t *Test.T) {
		result := Version()
		t.Ok(result != "")
		t.End()
	})
}

func TestHelpContainsUsage(t *testing.T) {
	Test.Test(t, "help: contains usage line", func(t *Test.T) {
		t.Match(Help(), "usage: subscriber")
		t.End()
	})
}

func TestHelpContainsVersion(t *testing.T) {
	Test.Test(t, "help: contains --version flag", func(t *Test.T) {
		t.Match(Help(), "--version")
		t.End()
	})
}

func TestHelpContainsHelpFlag(t *testing.T) {
	Test.Test(t, "help: contains --help flag", func(t *Test.T) {
		t.Match(Help(), "--help")
		t.End()
	})
}

func TestHelpContainsEnvSection(t *testing.T) {
	Test.Test(t, "help: contains environment variables section", func(t *Test.T) {
		t.Match(Help(), "environment variables:")
		t.End()
	})
}

func TestHelpContainsPORT(t *testing.T) {
	Test.Test(t, "help: contains PORT env var", func(t *Test.T) {
		t.Match(Help(), "PORT")
		t.End()
	})
}

func TestHelpVersionBeforeHelp(t *testing.T) {
	Test.Test(t, "help: --version appears before --help", func(t *Test.T) {
		result := Help()
		t.Ok(strings.Index(result, "--version") < strings.Index(result, "--help"))
		t.End()
	})
}

func TestHelpFromJSONInvalid(t *testing.T) {
	Test.Test(t, "help: HelpFromJSON returns fallback on invalid JSON", func(t *Test.T) {
		result := HelpFromJSON([]byte(`{invalid`))
		t.Equal(result, "usage: subscriber [options]\n(help unavailable)\n")
		t.End()
	})
}
