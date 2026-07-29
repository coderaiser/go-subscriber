package only_test

import (
	"os"
	"testing"

	tape "github.com/coderaiser/go-subscriber"
	"github.com/coderaiser/go-subscriber/cmd/tape/only"
)

// OnlyT extends tape.T with fixture-based FindOnly operator.
type OnlyT struct{ *tape.T }

func (t *OnlyT) FindOnly(file string, expected []only.OnlyCall) {
	t.TB().Helper()
	src, err := os.ReadFile("fixture/" + file + "/" + file + ".go")
	if err != nil {
		t.TB().Fatalf("fixture not found: %v", err)
	}
	calls, err := only.FindOnlyCallsInSource(string(src))
	if err != nil {
		t.TB().Fatalf("unexpected error: %v", err)
	}
	t.DeepEqual(calls, expected)
}

func (t *OnlyT) Pattern(calls []only.OnlyCall, expected string) {
	t.TB().Helper()
	got := only.BuildRunPattern(calls)
	t.Equal(got, expected)
}

func runOnlyTest(tb *testing.T, name string, fn func(*OnlyT)) {
	tape.Test(tb, name, func(base *tape.T) {
		fn(&OnlyT{T: base})
	})
}

func TestFindNoCalls(t *testing.T) {
	runOnlyTest(t, "only: no Only calls returns empty", func(t *OnlyT) {
		t.FindOnly("no-only", nil)
		t.End()
	})
}

func TestFindOneCall(t *testing.T) {
	runOnlyTest(t, "only: one Only call returns one result", func(t *OnlyT) {
		t.FindOnly("one-only", []only.OnlyCall{
			{Parent: "TestParser", Name: "parser: run action"},
		})
		t.End()
	})
}

func TestFindMultipleCalls(t *testing.T) {
	runOnlyTest(t, "only: multiple Only calls in same func", func(t *OnlyT) {
		t.FindOnly("multi-only", []only.OnlyCall{
			{Parent: "TestParser", Name: "parser: run action"},
			{Parent: "TestParser", Name: "parser: fail action"},
		})
		t.End()
	})
}

func TestFindCrossFunc(t *testing.T) {
	runOnlyTest(t, "only: Only calls in different TestXxx functions", func(t *OnlyT) {
		t.FindOnly("cross-func", []only.OnlyCall{
			{Parent: "TestParser", Name: "parser: run action"},
			{Parent: "TestRunner", Name: "runner: starts"},
		})
		t.End()
	})
}

func TestBuildPatternEmpty(t *testing.T) {
	runOnlyTest(t, "only: empty calls returns empty pattern", func(t *OnlyT) {
		t.Pattern(nil, "")
		t.End()
	})
}

func TestBuildPatternSingle(t *testing.T) {
	runOnlyTest(t, "only: single call builds pattern", func(t *OnlyT) {
		t.Pattern([]only.OnlyCall{
			{Parent: "TestParser", Name: "parser: run action"},
		}, "TestParser/parser:_run_action")
		t.End()
	})
}

func TestBuildPatternMultiple(t *testing.T) {
	runOnlyTest(t, "only: multiple calls joined with pipe", func(t *OnlyT) {
		t.Pattern([]only.OnlyCall{
			{Parent: "TestParser", Name: "parser: run action"},
			{Parent: "TestRunner", Name: "runner: starts"},
		}, "TestParser/parser:_run_action|TestRunner/runner:_starts")
		t.End()
	})
}

func TestFindOnlyCallsInDir(t *testing.T) {
	runOnlyTest(t, "only: FindOnlyCalls reads directory", func(t *OnlyT) {
		calls, err := only.FindOnlyCalls("fixture/one-only")
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Ok(len(calls) == 1)
		t.End()
	})
}

func TestFindOnlyCallsMissingDir(t *testing.T) {
	runOnlyTest(t, "only: FindOnlyCalls errors on missing dir", func(t *OnlyT) {
		_, err := only.FindOnlyCalls("nonexistent")
		t.Error(err)
		t.End()
	})
}

func TestFindOnlyCallsInSourceInvalidGo(t *testing.T) {
	runOnlyTest(t, "only: FindOnlyCallsInSource errors on invalid Go", func(t *OnlyT) {
		_, err := only.FindOnlyCallsInSource("this is not go code {{{{\"")
		t.Error(err)
		t.End()
	})
}

func TestIsOnlyCallBareIdent(t *testing.T) {
	os.Setenv("TAPE_CHECK_ASSERTIONS_COUNT", "0")
	defer os.Unsetenv("TAPE_CHECK_ASSERTIONS_COUNT")
	runOnlyTest(t, "only: bare Only() call is detected", func(t *OnlyT) {
		src := `package p

func TestFoo(t *testing.T) {
	Only(t, "scope: msg", func(t *T) {})
}`
		calls, err := only.FindOnlyCallsInSource(src)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.DeepEqual(calls, []only.OnlyCall{
			{Parent: "TestFoo", Name: "scope: msg"},
		})
		t.End()
	})
}

func TestBuildPatternSortsByParentThenName(t *testing.T) {
	runOnlyTest(t, "only: BuildRunPattern sorts deterministically", func(t *OnlyT) {
		got := only.BuildRunPattern([]only.OnlyCall{
			{Parent: "TestB", Name: "z: last"},
			{Parent: "TestA", Name: "a: first"},
		})
		t.Equal(got, "TestA/a:_first|TestB/z:_last")
		t.End()
	})
}

func TestIsOnlyCallOtherExpr(t *testing.T) {
	runOnlyTest(t, "only: non-Only call is ignored", func(t *OnlyT) {
		src := `package p

func TestFoo(t *testing.T) {
	Other("something.something")
}`
		calls, err := only.FindOnlyCallsInSource(src)
		if err != nil {
			t.TB().Fatal(err)
		}
		t.Ok(len(calls) == 0)
		t.End()
	})
}
