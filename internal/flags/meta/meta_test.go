package meta_test

import (
	"testing"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-subscriber/internal/flags/meta"
)

func TestVersionLineDefault(t *testing.T) {
	tape.Test(t, "meta: VersionLine returns default when not set", func(t *tape.T) {
		result := meta.VersionLine()
		t.Match(result, "v")
		t.End()
	})
}

func TestHelpDefault(t *testing.T) {
	tape.Test(t, "meta: Help returns default when not set", func(t *tape.T) {
		result := meta.Help()
		t.Match(result, "usage:")
		t.End()
	})
}

func TestSetAndVersionLine(t *testing.T) {
	tape.Test(t, "meta: Set overrides VersionLine", func(t *tape.T) {
		meta.Set(func() string { return "v9.9.9" }, nil)
		t.Equal(meta.VersionLine(), "v9.9.9")
		t.End()
	})
}

func TestSetAndHelp(t *testing.T) {
	tape.Test(t, "meta: Set overrides Help", func(t *tape.T) {
		meta.Set(nil, func() string { return "custom help" })
		t.Equal(meta.Help(), "custom help")
		t.End()
	})
}
