package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coderaiser/go-subscriber/internal/store"
	Test "github.com/coderaiser/go-subscriber/internal/tape"
)

func TestSeedLoadsSubscribers(t *testing.T) {
	Test.Test(t, "store: Seed loads subscribers from JSON file", func(t *Test.T) {
		dir := t.TB().TempDir()
		path := filepath.Join(dir, "test.json")
		os.WriteFile(path, []byte(`[
			{"msisdn": "48100000001", "service_id": "demo-svc", "state": "active"},
			{"msisdn": "48100000002", "service_id": "demo-svc", "state": "trial"}
		]`), 0644)
		ss := store.NewStateStore()
		err := store.Seed(path, ss)
		t.NoError(err)
		ptr, _ := ss.Get("48100000001:demo-svc")
		t.Ok(ptr != nil)
		t.Equal(*ptr, "active")
		ptr2, _ := ss.Get("48100000002:demo-svc")
		t.Ok(ptr2 != nil)
		t.Equal(*ptr2, "trial")
		t.End()
	})
}

func TestSeedMissingFile(t *testing.T) {
	Test.Test(t, "store: Seed with missing file returns no error", func(t *Test.T) {
		ss := store.NewStateStore()
		err := store.Seed("/nonexistent/path.json", ss)
		t.NoError(err)
		t.End()
	})
}

func TestSeedInvalidJSON(t *testing.T) {
	Test.Test(t, "store: Seed with invalid JSON returns error", func(t *Test.T) {
		dir := t.TB().TempDir()
		path := filepath.Join(dir, "bad.json")
		os.WriteFile(path, []byte(`{invalid}`), 0644)
		ss := store.NewStateStore()
		err := store.Seed(path, ss)
		t.Error(err)
		t.End()
	})
}

func TestSeedUnreadableFile(t *testing.T) {
	Test.Test(t, "store: Seed with unreadable file returns error", func(t *Test.T) {
		if os.Getuid() == 0 {
			t.TB().Skip("running as root; permission checks do not apply")
		}
		dir := t.TB().TempDir()
		path := filepath.Join(dir, "noperm.json")
		os.WriteFile(path, []byte(`[]`), 0000)
		ss := store.NewStateStore()
		err := store.Seed(path, ss)
		t.Error(err)
		t.End()
	})
}
