package tape

import "testing"

// Extension defines an operator that extends *T with custom assertions.
type Extension func(*T)

// Extensions groups multiple extensions.
type Extensions []Extension

// Extend applies extensions to a test function.
func Extend(extensions Extensions) func(t *testing.T, name string, fn func(t *T)) {
	return func(t *testing.T, name string, fn func(t *T)) {
		Test(t, name, func(tt *T) {
			for _, ext := range extensions {
				ext(tt)
			}
			fn(tt)
		})
	}
}
