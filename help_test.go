package main

import (
	"strings"
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestHelpContainsUsage(t *testing.T) {
	tape.Test(t, "help: contains usage line", func(t *tape.T) {
		t.Match(Help(), "usage: subscriber")
		t.End()
	})
}

func TestHelpContainsVersion(t *testing.T) {
	tape.Test(t, "help: contains --version flag", func(t *tape.T) {
		t.Match(Help(), "--version")
		t.End()
	})
}

func TestHelpContainsHelpFlag(t *testing.T) {
	tape.Test(t, "help: contains --help flag", func(t *tape.T) {
		t.Match(Help(), "--help")
		t.End()
	})
}

func TestHelpContainsEnvSection(t *testing.T) {
	tape.Test(t, "help: contains environment variables section", func(t *tape.T) {
		t.Match(Help(), "environment variables:")
		t.End()
	})
}

func TestHelpContainsPORT(t *testing.T) {
	tape.Test(t, "help: contains PORT env var", func(t *tape.T) {
		t.Match(Help(), "PORT")
		t.End()
	})
}

func TestHelpVersionBeforeHelp(t *testing.T) {
	tape.Test(t, "help: --version appears before --help", func(t *tape.T) {
		result := Help()
		t.Ok(strings.Index(result, "--version") < strings.Index(result, "--help"))
		t.End()
	})
}

func TestHelpFromJSONInvalid(t *testing.T) {
	tape.Test(t, "help: HelpFromJSON returns fallback on invalid JSON", func(t *tape.T) {
		result := HelpFromJSON([]byte(`{invalid`))
		t.Equal(result, "usage: subscriber [options]\n(help unavailable)\n")
		t.End()
	})
}
