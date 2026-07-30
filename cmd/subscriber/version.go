package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed package.json
var packageJSONBytes []byte

// VersionFromJSON parses the version field from package.json bytes.
func VersionFromJSON(data []byte) string {
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "unknown"
	}
	if pkg.Version == "" {
		return "unknown"
	}
	return pkg.Version
}

// Version returns the version embedded from package.json at build time.
func Version() string {
	return VersionFromJSON(packageJSONBytes)
}

// VersionLine returns "vX.Y.Z" for -v/--version output.
func VersionLine() string {
	return fmt.Sprintf("v%s", Version())
}
