package flags_test

import (
	"testing"

	"github.com/coderaiser/go-subscriber/internal/flags"
	Test "github.com/coderaiser/go-subscriber/internal/tape"
)

func TestNoFlags(t *testing.T) {
	Test.Test(t, "flags: no args returns exitCode -1 (continue)", func(t *Test.T) {
		result := flags.Parse([]string{})
		t.Ok(result.ExitCode == -1)
		t.End()
	})
}

func TestVersionShort(t *testing.T) {
	Test.Test(t, "flags: -v returns exitCode 0", func(t *Test.T) {
		result := flags.Parse([]string{"-v"})
		t.Ok(result.ExitCode == 0)
		t.End()
	})
}

func TestVersionLong(t *testing.T) {
	Test.Test(t, "flags: --version returns exitCode 0", func(t *Test.T) {
		result := flags.Parse([]string{"--version"})
		t.Ok(result.ExitCode == 0)
		t.End()
	})
}

func TestVersionOutput(t *testing.T) {
	Test.Test(t, "flags: --version output contains v", func(t *Test.T) {
		result := flags.Parse([]string{"--version"})
		t.Match(result.Output, "v")
		t.End()
	})
}

func TestHelpShort(t *testing.T) {
	Test.Test(t, "flags: -h returns exitCode 0", func(t *Test.T) {
		result := flags.Parse([]string{"-h"})
		t.Ok(result.ExitCode == 0)
		t.End()
	})
}

func TestHelpLong(t *testing.T) {
	Test.Test(t, "flags: --help returns exitCode 0", func(t *Test.T) {
		result := flags.Parse([]string{"--help"})
		t.Ok(result.ExitCode == 0)
		t.End()
	})
}

func TestHelpOutput(t *testing.T) {
	Test.Test(t, "flags: --help output contains usage", func(t *Test.T) {
		result := flags.Parse([]string{"--help"})
		t.Match(result.Output, "usage:")
		t.End()
	})
}

func TestUnknownFlag(t *testing.T) {
	Test.Test(t, "flags: unknown flag returns exitCode 1", func(t *Test.T) {
		result := flags.Parse([]string{"--foo"})
		t.Ok(result.ExitCode == 1)
		t.End()
	})
}

func TestUnknownFlagOutput(t *testing.T) {
	Test.Test(t, "flags: unknown flag output contains flag name", func(t *Test.T) {
		result := flags.Parse([]string{"--foo"})
		t.Match(result.Output, "--foo")
		t.End()
	})
}
