//go:build !no_external

package statemachine

import "testing"

func TestFileSourceTOMLLoadsBurntSushi(t *testing.T) {
	src := FileSource{Path: "testdata/runner.toml"}
	defs, err := src.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(defs) == 0 {
		t.Fatal("expected transitions from TOML file")
	}
}
