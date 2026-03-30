package task

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/seifenkehrer/seifenkehrer/internal/storage"
)

type taskConfigEntry struct {
	Disabled bool   `json:"disabled,omitempty"`
	Interval string `json:"interval,omitempty"`
}

type taskConfig struct {
	path    string
	store   storage.Storage
	Entries map[string]taskConfigEntry `json:"tasks"`
}

func loadTaskConfig(configDir string) (*taskConfig, error) {
	return loadTaskConfigWithStorage(configDir, storage.FileStorage{})
}

func loadTaskConfigWithStorage(configDir string, store storage.Storage) (*taskConfig, error) {
	path := filepath.Join(configDir, "config.json")
	c := &taskConfig{
		path:    path,
		store:   store,
		Entries: make(map[string]taskConfigEntry),
	}

	data, err := store.Load(path)
	if err != nil {
		return nil, err
	}

	if data != nil {
		if err := json.Unmarshal(data, c); err != nil {
			return nil, fmt.Errorf("parsing config %s: %w", path, err)
		}
		if c.Entries == nil {
			c.Entries = make(map[string]taskConfigEntry)
		}
	}

	return c, nil
}

func (c *taskConfig) isDisabled(taskName string) bool {
	e, ok := c.Entries[taskName]
	return ok && e.Disabled
}

func (c *taskConfig) intervalOverride(taskName string) string {
	e, ok := c.Entries[taskName]
	if !ok {
		return ""
	}
	return e.Interval
}

func (c *taskConfig) setDisabled(taskName string, disabled bool) {
	e := c.Entries[taskName]
	e.Disabled = disabled
	if !e.Disabled && e.Interval == "" {
		delete(c.Entries, taskName)
		return
	}
	c.Entries[taskName] = e
}

func (c *taskConfig) setInterval(taskName string, interval string) {
	e := c.Entries[taskName]
	e.Interval = interval
	if !e.Disabled && e.Interval == "" {
		delete(c.Entries, taskName)
		return
	}
	c.Entries[taskName] = e
}

func (c *taskConfig) save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return c.store.Save(c.path, data)
}

func Enable(configDir string, taskName string) error {
	c, err := loadTaskConfig(configDir)
	if err != nil {
		return err
	}
	c.setDisabled(taskName, false)
	return c.save()
}

func Disable(configDir string, taskName string) error {
	c, err := loadTaskConfig(configDir)
	if err != nil {
		return err
	}
	c.setDisabled(taskName, true)
	return c.save()
}

func SetInterval(configDir string, taskName string, interval string) error {
	c, err := loadTaskConfig(configDir)
	if err != nil {
		return err
	}
	c.setInterval(taskName, interval)
	return c.save()
}
