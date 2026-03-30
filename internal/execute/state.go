package execute

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/seifenkehrer/seifenkehrer/internal/storage"
)

type state struct {
	path    string
	store   storage.Storage
	entries map[string]stateEntry
}

type stateEntry struct {
	LastRun time.Time `json:"last_run"`
}

func loadState(configDir string, store storage.Storage) (*state, error) {
	path := filepath.Join(configDir, "state.json")
	s := &state{
		path:    path,
		store:   store,
		entries: make(map[string]stateEntry),
	}

	data, err := store.Load(path)
	if err != nil {
		return nil, err
	}

	if data != nil {
		if err := json.Unmarshal(data, &s.entries); err != nil {
			return nil, fmt.Errorf("parsing state file %s: %w", path, err)
		}
	}

	return s, nil
}

func (s *state) ShouldSkip(taskName string, interval time.Duration) bool {
	if interval <= 0 {
		return false
	}
	entry, ok := s.entries[taskName]
	if !ok {
		return false
	}
	return time.Since(entry.LastRun) < interval
}

func (s *state) RecordRun(taskName string) {
	s.entries[taskName] = stateEntry{LastRun: time.Now()}
}

func (s *state) Save() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}
	return s.store.Save(s.path, data)
}
