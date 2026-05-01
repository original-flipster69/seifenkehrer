package execute

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/original-flipster69/seifenkehrer/internal/storage"
)

type state struct {
	path    string
	store   storage.Storage
	entries map[string]stateEntry
}

type stateEntry struct {
	LastRun time.Time `json:"last_run"`
}

func load(configDir string, store storage.Storage) (*state, error) {
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

func (s *state) lastRun(task string) *time.Time {
	entry, ok := s.entries[task]
	if !ok {
		return nil
	}
	return &entry.LastRun
}

func (s *state) RecordRun(task string) {
	s.entries[task] = stateEntry{LastRun: time.Now()}
	s.save()
}

func (s *state) save() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling state: %w", err)
	}
	return s.store.Save(s.path, data)
}
