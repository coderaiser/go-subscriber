package tape_test

import (
	"errors"
	"testing"

	Test "github.com/coderaiser/go-subscriber/internal/tape"
)

func TestTapeTestRuns(t *testing.T) {
	Test.Test(t, "tape: test runs callback", func(t *Test.T) {
		t.Ok(true)
		t.End()
	})
}

func TestTapeOkPass(t *testing.T) {
	Test.Test(t, "tape: Ok passes for true", func(t *Test.T) {
		t.Ok(true)
		t.End()
	})
}

func TestTapeNotOkPass(t *testing.T) {
	Test.Test(t, "tape: NotOk passes for false", func(t *Test.T) {
		t.NotOk(false)
		t.End()
	})
}

func TestTapeNoErrorPass(t *testing.T) {
	Test.Test(t, "tape: NoError passes for nil", func(t *Test.T) {
		t.NoError(nil)
		t.End()
	})
}

func TestTapeErrorPass(t *testing.T) {
	Test.Test(t, "tape: Error passes for non-nil", func(t *Test.T) {
		t.Error(errors.New("err"))
		t.End()
	})
}

func TestTapeEqualPass(t *testing.T) {
	Test.Test(t, "tape: Equal passes for equal values", func(t *Test.T) {
		t.Equal(1, 1)
		t.End()
	})
}

func TestTapeNotEqualPass(t *testing.T) {
	Test.Test(t, "tape: NotEqual passes for different values", func(t *Test.T) {
		t.NotEqual(1, 2)
		t.End()
	})
}

func TestTapeDeepEqual(t *testing.T) {
	Test.Test(t, "tape: DeepEqual passes for equal slices", func(t *Test.T) {
		t.DeepEqual([]int{1, 2}, []int{1, 2})
		t.End()
	})
}

func TestTapeNotDeepEqual(t *testing.T) {
	Test.Test(t, "tape: NotDeepEqual passes for different slices", func(t *Test.T) {
		t.NotDeepEqual([]int{1, 2}, []int{3, 4})
		t.End()
	})
}

func TestTapeMatchPass(t *testing.T) {
	Test.Test(t, "tape: Match passes for matching string", func(t *Test.T) {
		t.Match("hello world", "hello")
		t.End()
	})
}

func TestTapeNotMatchPass(t *testing.T) {
	Test.Test(t, "tape: NotMatch passes for non-matching string", func(t *Test.T) {
		t.NotMatch("hello world", "goodbye")
		t.End()
	})
}

func TestTapeTB(t *testing.T) {
	Test.Test(t, "tape: TB returns *testing.T", func(t *Test.T) {
		tb := t.TB()
		t.Ok(tb != nil)
		t.End()
	})
}


func TestTapeComment(t *testing.T) {
	Test.Test(t, "tape: Comment does not fail", func(t *Test.T) {
		t.Comment("hello")
		t.End()
	})
}
