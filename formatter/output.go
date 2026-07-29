package formatter

import (
	"regexp"
	"strings"
)

// OutputFields extracted from go test -json output lines.
type OutputFields struct {
	Operator   string
	Result     string
	Expected   string
	At         string
	ErrorStack string
	Raw        string
}

var (
	reAt       = regexp.MustCompile(`^\s+(\S+\.go:\d+:)\s*`)
	reResult   = regexp.MustCompile(`^\s+result:\s+(.+)`)
	reExpected = regexp.MustCompile(`^\s+expected:\s+(.+)`)
	reOperator = regexp.MustCompile(`^\s+(Equal|NotEqual|Ok|NotOk|DeepEqual|NotDeepEqual|Match|NotMatch|Error|NoError|Pass|Fail)`)
)

// ParseOutput parses buffered output lines from go test -json.
// Returns extracted fields for use in Fail formatter event.
func ParseOutput(lines []string) OutputFields {
	raw := strings.Join(lines, "")
	fields := OutputFields{Raw: raw}

	for _, line := range lines {
		if m := reAt.FindStringSubmatch(line); m != nil {
			fields.At = m[1]
			continue
		}
		if m := reResult.FindStringSubmatch(line); m != nil {
			fields.Result = strings.TrimSpace(m[1])
			continue
		}
		if m := reExpected.FindStringSubmatch(line); m != nil {
			fields.Expected = strings.TrimSpace(m[1])
			continue
		}
		if m := reOperator.FindStringSubmatch(line); m != nil {
			fields.Operator = m[1]
			continue
		}
	}

	return fields
}
