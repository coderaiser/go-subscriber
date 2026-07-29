//go:build no_external

package statemachine

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSourceTOMLLoadsHandRolled(t *testing.T) {
	src := FileSource{Path: "testdata/runner.toml"}
	defs, err := src.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(defs) == 0 {
		t.Fatal("expected transitions from hand-rolled TOML parser")
	}
}

func TestHandRolledTOMLSectionsAndComments(t *testing.T) {
	src := FileSource{Path: "testdata/runner_full.toml"}
	defs, err := src.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(defs) != 3 {
		t.Fatalf("expected 3 transitions, got %d", len(defs))
	}
}

func TestHandRolledTOMLMalformedLine(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "malformed.toml", `[transitions.idle]
run = "running"
borked
`)
	src := FileSource{Path: dir + "/malformed.toml"}
	defs, err := src.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(defs) != 1 {
		t.Fatalf("expected 1 transition, got %d", len(defs))
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// errReader errors on the second Read call
type errReader struct{ calls int }

func (r *errReader) Read(p []byte) (int, error) {
	r.calls++
	if r.calls > 1 {
		return 0, errors.New("read error")
	}
	content := "[transitions.idle]\n"
	n := copy(p, content)
	return n, nil
}

func TestParseTOMLReaderScannerError(t *testing.T) {
	_, err := parseTOMLReader(&errReader{}, "test")
	if err == nil {
		t.Fatal("expected error from scanner failure")
	}
}
