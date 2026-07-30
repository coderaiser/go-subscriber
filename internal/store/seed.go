package store

import (
	"encoding/json"
	"os"
)

// seedEntry represents a single subscriber in the seed file.
type seedEntry struct {
	Msisdn    string `json:"msisdn"`
	ServiceID string `json:"service_id"`
	State     string `json:"state"`
}

// Seed loads a JSON file of pre-existing subscribers into StateStore.
// Format: [{"msisdn": "123", "service_id": "svc1", "state": "active"}, ...]
// Missing or empty file is not an error — service starts with empty state.
func Seed(path string, ss *StateStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var entries []seedEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	for _, e := range entries {
		key := e.Msisdn + ":" + e.ServiceID
		ss.Set(key, e.State)
	}

	return nil
}
