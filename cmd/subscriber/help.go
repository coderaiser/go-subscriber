package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed help.json
var helpJSONBytes []byte

type helpConfig struct {
	Flags     map[string]string `json:"flags"`
	Env       map[string]string `json:"env"`
	FlagOrder []string          `json:"flagOrder"`
	EnvOrder  []string          `json:"envOrder"`
}

// Help returns the formatted help string for the subscriber binary.
func Help() string {
	return HelpFromJSON(helpJSONBytes)
}

// HelpFromJSON parses help config from JSON bytes and formats it.
func HelpFromJSON(data []byte) string {
	var cfg helpConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return "usage: subscriber [options]\n(help unavailable)\n"
	}

	var b strings.Builder
	b.WriteString("usage: subscriber [options]\n\nflags:\n")

	for _, flag := range cfg.FlagOrder {
		if desc, ok := cfg.Flags[flag]; ok {
			fmt.Fprintf(&b, "  %-22s %s\n", flag, desc)
		}
	}

	if len(cfg.Env) > 0 {
		b.WriteString("\nenvironment variables:\n")
		for _, key := range cfg.EnvOrder {
			if desc, ok := cfg.Env[key]; ok {
				fmt.Fprintf(&b, "  %-22s %s\n", key, desc)
			}
		}
	}

	return b.String()
}
