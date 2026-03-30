package task

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/seifenkehrer/seifenkehrer/internal/storage"
)

type configEntry struct {
	Disabled bool   `json:"disabled,omitempty"`
	Interval string `json:"interval,omitempty"`
}

type config struct {
	path    string
	store   storage.Storage
	Entries map[string]configEntry `json:"tasks"`
}

func loadConfig(configDir string) (*config, error) {
	path := filepath.Join(configDir, "config.json")
	store := storage.FileStorage{}
	c := &config{
		path:    path,
		store:   store,
		Entries: make(map[string]configEntry),
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
			c.Entries = make(map[string]configEntry)
		}
	}

	return c, nil
}

func (c *config) isDisabled(task string) bool {
	e, ok := c.Entries[task]
	return ok && e.Disabled
}

func (c *config) interval(task string) string {
	e, ok := c.Entries[task]
	if !ok {
		return ""
	}
	return e.Interval
}

func (c *config) setDisabled(task string, disabled bool) {
	e := c.Entries[task]
	e.Disabled = disabled
	if !e.Disabled && e.Interval == "" {
		delete(c.Entries, task)
		return
	}
	c.Entries[task] = e
}

func (c *config) setInterval(task string, interval string) {
	e := c.Entries[task]
	e.Interval = interval
	if !e.Disabled && e.Interval == "" {
		delete(c.Entries, task)
		return
	}
	c.Entries[task] = e
}

func (c *config) save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return c.store.Save(c.path, data)
}

func Enable(configDir string, task string) error {
	c, err := loadConfig(configDir)
	if err != nil {
		return err
	}
	c.setDisabled(task, false)
	return c.save()
}

func Disable(configDir string, task string) error {
	c, err := loadConfig(configDir)
	if err != nil {
		return err
	}
	c.setDisabled(task, true)
	return c.save()
}

func SetInterval(configDir string, task string, interval string) error {
	c, err := loadConfig(configDir)
	if err != nil {
		return err
	}
	c.setInterval(task, interval)
	return c.save()
}
