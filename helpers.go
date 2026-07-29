package tape

import (
	"sync"
	"testing"

	"github.com/coderaiser/go-subscriber/assert"
	"github.com/coderaiser/go-subscriber/internal/config"
)

var (
	mu    sync.Mutex
	count = make(map[*testing.T]int)
)

func assertOne(t *testing.T) {
	t.Helper()
	mu.Lock()
	count[t] = 0
	mu.Unlock()
	// cleanup prevents memory leak — removes entry after test finishes
	t.Cleanup(func() {
		mu.Lock()
		delete(count, t)
		mu.Unlock()
	})
}

func hit(t *testing.T) {
	t.Helper()
	mu.Lock()
	count[t]++
	c := count[t]
	mu.Unlock()
	if c <= 1 {
		return
	}
	if !config.CheckAssertionsCount() {
		return
	}
	assert.HitCheck(t, c)
}

func truthy(v any) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case int:
		return val != 0
	case string:
		return val != ""
	default:
		return true
	}
}
