package formatter

import (
	"encoding/json"
)

// JSONLinesFormatter outputs streaming JSON.
type JSONLinesFormatter struct {
	total int
}

func NewJSONLines(total int) *JSONLinesFormatter {
	return &JSONLinesFormatter{total: total}
}

func (f *JSONLinesFormatter) Start(total int) string  { return "" }
func (f *JSONLinesFormatter) Test(name string) string { return "" }
func (f *JSONLinesFormatter) Comment(msg string) string { return "" }

func (f *JSONLinesFormatter) TestEnd(count, total, failed int, name string) string {
	b, _ := json.Marshal(map[string]any{
		"count":  count,
		"total":  total,
		"failed": failed,
		"test":   name,
	})
	return string(b) + "\n"
}

func (f *JSONLinesFormatter) Success(count int, message string) string { return "" }

func (f *JSONLinesFormatter) Fail(count int, message, operator string, result, expected any, output, at, errorStack string) string {
	b, _ := json.Marshal(map[string]any{
		"test":     message,
		"at":       at,
		"count":    count,
		"message":  message,
		"operator": operator,
		"result":   result,
		"expected": expected,
	})
	return string(b) + "\n"
}

func (f *JSONLinesFormatter) End(count, passed, failed, skipped int) string {
	b, _ := json.Marshal(map[string]any{
		"count":   count,
		"passed":  passed,
		"failed":  failed,
		"skipped": skipped,
	})
	return string(b) + "\n"
}
