package statemachine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// TransitionSource loads transition definitions once at startup.
// Load() is called once in New() — never per Apply().
type TransitionSource interface {
	Load() ([]TransitionDef, error)
}

// MemorySource holds transitions in memory.
// Primary use: tests — no file I/O needed.
type MemorySource struct {
	Defs []TransitionDef
}

func (s *MemorySource) Load() ([]TransitionDef, error) {
	return s.Defs, nil
}

// FileSource loads transitions from a file.
// Supported extensions: .toml, .json
// .toml: uses BurntSushi/toml by default, hand-rolled with -tags no_external
// .json: uses encoding/json (stdlib always)
type FileSource struct {
	Path string
}

func (s FileSource) Load() ([]TransitionDef, error) {
	switch filepath.Ext(s.Path) {
	case ".toml":
		return loadTOML(s.Path)
	case ".json":
		return loadJSON(s.Path)
	default:
		return nil, fmt.Errorf("FileSource: unsupported extension %q", filepath.Ext(s.Path))
	}
}

// loadJSON loads transitions from a JSON file using stdlib encoding/json.
func loadJSON(path string) ([]TransitionDef, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("FileSource: read %s: %w", path, err)
	}
	var cfg struct {
		Transitions map[string]map[string]string `json:"transitions"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("FileSource: parse %s: %w", path, err)
	}
	return toTransitionDefs(cfg.Transitions), nil
}

func toTransitionDefs(transitions map[string]map[string]string) []TransitionDef {
	var defs []TransitionDef
	for from, events := range transitions {
		for event, to := range events {
			defs = append(defs, TransitionDef{From: from, Event: event, To: to})
		}
	}
	return defs
}
