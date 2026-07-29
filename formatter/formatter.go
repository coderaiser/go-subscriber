package formatter

import (
	"io"
	"os"

	"github.com/coderaiser/go-subscriber/model"
)

// Formatter matches supertape's formatter event API exactly.
// Each method returns the string to write to output (empty = no output).
type Formatter interface {
	Start(total int) string
	Test(name string) string
	TestEnd(count, total, failed int, name string) string
	Success(count int, message string) string
	Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string
	Comment(message string) string
	End(count, passed, failed, skipped int) string
}

// State tracks streaming state for FromEvent.
type State struct {
	count     int
	failed    int
	total     int
	outputs   map[string][]string // test name → buffered output lines
	formatter Formatter
	w         io.Writer
}

// New returns the formatter for the given format string.
// Auto-detects CI: CI=true → tap.
// Defaults to progress-bar in terminal.
func New(format string, w io.Writer, total int) *State {
	if os.Getenv("CI") == "true" {
		format = "tap"
	} else if format == "" {
		format = "progress-bar"
	}
	var f Formatter
	switch format {
	case "tap":
		f = NewTAP()
	case "short":
		f = NewShort()
	case "fail":
		f = NewFail()
	case "time":
		f = NewTime(total)
	case "json-lines":
		f = NewJSONLines(total)
	default: // progress-bar
		f = NewProgressBar(total)
	}
	s := &State{
		total:     total,
		outputs:   make(map[string][]string),
		formatter: f,
		w:         w,
	}
	write(w, f.Start(total))
	return s
}

// FromEvent routes a model.Event to the appropriate formatter method.
// Writes output immediately — streaming.
func (s *State) FromEvent(e model.Event) {
	if e.Test == "" {
		return
	}
	switch e.Action {
	case "run":
		write(s.w, s.formatter.Test(e.Test))
	case "output":
		s.outputs[e.Test] = append(s.outputs[e.Test], e.Output)
	case "pass":
		s.count++
		write(s.w, s.formatter.Success(s.count, e.Test))
		write(s.w, s.formatter.TestEnd(s.count, s.total, s.failed, e.Test))
	case "fail":
		s.count++
		s.failed++
		lines := s.outputs[e.Test]
		fields := ParseOutput(lines)
		write(s.w, s.formatter.Fail(
			s.count, e.Test,
			fields.Operator, fields.Result, fields.Expected,
			fields.Raw, fields.At, fields.ErrorStack,
		))
		write(s.w, s.formatter.TestEnd(s.count, s.total, s.failed, e.Test))
	case "skip":
		s.count++
		write(s.w, s.formatter.TestEnd(s.count, s.total, s.failed, e.Test))
	}
}

// End writes the final summary.
func (s *State) End(passed, failed, skipped int) {
	write(s.w, s.formatter.End(s.count, passed, failed, skipped))
}

func write(w io.Writer, s string) {
	if s != "" {
		w.Write([]byte(s))
	}
}
