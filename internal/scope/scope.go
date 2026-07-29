package scope

import "strings"

// Valid checks name matches "scope: message" format.
func Valid(name string) bool {
	if name == "" {
		return false
	}
	idx := strings.Index(name, ": ")
	if idx <= 0 {
		return false
	}
	scope := name[:idx]
	msg := name[idx+2:]
	if scope == "" || msg == "" {
		return false
	}
	return true
}
