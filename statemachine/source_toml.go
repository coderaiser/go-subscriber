//go:build !no_external

package statemachine

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

func loadTOML(path string) ([]TransitionDef, error) {
	var cfg struct {
		Transitions map[string]map[string]string `toml:"transitions"`
	}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("FileSource: load %s: %w", path, err)
	}
	return toTransitionDefs(cfg.Transitions), nil
}
