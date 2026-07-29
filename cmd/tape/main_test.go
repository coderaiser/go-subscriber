package main

import (
	"strings"
	"testing"

	tape "github.com/coderaiser/go-subscriber"
)

func TestRunHelp(t *testing.T) {
	tape.Test(t, "main: -h prints usage", func(t *tape.T) {
		var out, errOut strings.Builder
		code := run([]string{"-h"}, &out, &errOut)
		t.Equal(code, 0)
		t.End()
	})
}

func TestRunVersion(t *testing.T) {
	tape.Test(t, "main: -v prints version", func(t *tape.T) {
		var out, errOut strings.Builder
		code := run([]string{"-v"}, &out, &errOut)
		t.Equal(code, 0)
		t.End()
	})
}

func TestRunVersionOutput(t *testing.T) {
	tape.Test(t, "main: -v output contains version string", func(t *tape.T) {
		var out, errOut strings.Builder
		run([]string{"-v"}, &out, &errOut)
		t.Match(out.String(), `\d+\.\d+\.\d+`)
		t.End()
	})
}

func TestRunHelpOutput(t *testing.T) {
	tape.Test(t, "main: -h output contains Usage", func(t *tape.T) {
		var out, errOut strings.Builder
		run([]string{"-h"}, &out, &errOut)
		t.Match(out.String(), `Usage`)
		t.End()
	})
}

func TestUsage(t *testing.T) {
	tape.Test(t, "main: usage contains go-subscriber", func(t *tape.T) {
		t.Match(usage, `go-subscriber`)
		t.End()
	})
}

func TestVersionConst(t *testing.T) {
	tape.Test(t, "main: version is 1.0.0", func(t *tape.T) {
		t.Equal(version, "1.0.0")
		t.End()
	})
}
