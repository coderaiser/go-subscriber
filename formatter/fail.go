package formatter

import "fmt"

// FailFormatter is short formatter that prefixes fail output with test name.
// Matches supertape formatter-fail which re-exports tap but adds test name.
type FailFormatter struct {
	ShortFormatter
	currentTest string
}

func NewFail() *FailFormatter { return &FailFormatter{} }

func (f *FailFormatter) Test(name string) string {
	f.currentTest = name
	return ""
}

func (f *FailFormatter) Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string {
	prefix := ""
	if f.currentTest != "" {
		prefix = fmt.Sprintf("# %s\n", f.currentTest)
	}
	return prefix + f.ShortFormatter.Fail(count, message, operator, result, expected, output, at, errorStack)
}
