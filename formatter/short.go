package formatter

// ShortFormatter is tap minus stack trace in Fail.
type ShortFormatter struct {
	TAPFormatter
}

func NewShort() *ShortFormatter { return &ShortFormatter{} }

func (f *ShortFormatter) Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string {
	// same as TAP but errorStack is always omitted
	return f.TAPFormatter.Fail(count, message, operator, result, expected, output, at, "")
}
