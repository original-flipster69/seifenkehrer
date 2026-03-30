package task

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigEmpty(t *testing.T) {
	dir := t.TempDir()
	cfg, err := loadConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.isDisabled("anything") {
		t.Error("empty config should not disable any task")
	}
	if cfg.interval("anything") != "" {
		t.Error("empty config should have no interval override")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "config.json", `{
  "tasks": {
    "cleanup": {"disabled": true, "interval": "12h"}
  }
}`)

	cfg, err := loadConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.isDisabled("cleanup") {
		t.Error("expected cleanup to be disabled")
	}
	if cfg.interval("cleanup") != "12h" {
		t.Errorf("expected interval 12h, got %q", cfg.interval("cleanup"))
	}
}

func TestLoadConfigMalformedJSON(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "config.json", `not json`)

	_, err := loadConfig(dir)
	if err == nil {
		t.Error("expected error for malformed JSON")
	}
}

func TestSetDisabledTrue(t *testing.T) {
	dir := t.TempDir()
	cfg, _ := loadConfig(dir)

	cfg.setDisabled("mytask", true)

	if !cfg.isDisabled("mytask") {
		t.Error("expected task to be disabled")
	}
}

func TestSetDisabledFalseRemovesEntry(t *testing.T) {
	dir := t.TempDir()
	cfg, _ := loadConfig(dir)

	cfg.setDisabled("mytask", true)
	cfg.setDisabled("mytask", false)

	if cfg.isDisabled("mytask") {
		t.Error("expected task to be enabled")
	}
	if _, ok := cfg.Entries["mytask"]; ok {
		t.Error("expected entry to be removed when disabled=false and no interval")
	}
}

func TestSetDisabledFalseKeepsEntryWithInterval(t *testing.T) {
	dir := t.TempDir()
	cfg, _ := loadConfig(dir)

	cfg.setDisabled("mytask", true)
	cfg.setInterval("mytask", "6h")
	cfg.setDisabled("mytask", false)

	if _, ok := cfg.Entries["mytask"]; !ok {
		t.Error("expected entry to remain when interval is set")
	}
	if cfg.interval("mytask") != "6h" {
		t.Error("expected interval to be preserved")
	}
}

func TestSetIntervalClearsEntry(t *testing.T) {
	dir := t.TempDir()
	cfg, _ := loadConfig(dir)

	cfg.setInterval("mytask", "6h")
	cfg.setInterval("mytask", "")

	if _, ok := cfg.Entries["mytask"]; ok {
		t.Error("expected entry to be removed when interval cleared and not disabled")
	}
}

func TestSaveAndReloadConfig(t *testing.T) {
	dir := t.TempDir()
	cfg, _ := loadConfig(dir)

	cfg.setDisabled("task1", true)
	cfg.setInterval("task2", "48h")
	if err := cfg.save(); err != nil {
		t.Fatalf("save error: %v", err)
	}

	cfg2, err := loadConfig(dir)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if !cfg2.isDisabled("task1") {
		t.Error("expected task1 to be disabled after reload")
	}
	if cfg2.interval("task2") != "48h" {
		t.Errorf("expected task2 interval 48h after reload, got %q", cfg2.interval("task2"))
	}
}

func TestEnableDisableRoundTrip(t *testing.T) {
	dir := t.TempDir()

	if err := Disable(dir, "test-task"); err != nil {
		t.Fatalf("Disable error: %v", err)
	}

	cfg, _ := loadConfig(dir)
	if !cfg.isDisabled("test-task") {
		t.Error("expected task to be disabled")
	}

	if err := Enable(dir, "test-task"); err != nil {
		t.Fatalf("Enable error: %v", err)
	}

	cfg, _ = loadConfig(dir)
	if cfg.isDisabled("test-task") {
		t.Error("expected task to be enabled")
	}
}

func TestSetIntervalPersists(t *testing.T) {
	dir := t.TempDir()

	if err := SetInterval(dir, "test-task", "72h"); err != nil {
		t.Fatalf("SetInterval error: %v", err)
	}

	cfg, _ := loadConfig(dir)
	if cfg.interval("test-task") != "72h" {
		t.Errorf("expected interval 72h, got %q", cfg.interval("test-task"))
	}
}

func TestSaveProducesValidJSON(t *testing.T) {
	dir := t.TempDir()
	cfg, _ := loadConfig(dir)

	cfg.setDisabled("t1", true)
	cfg.setInterval("t2", "1h")
	cfg.save()

	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatalf("reading saved config: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("saved config is not valid JSON: %v", err)
	}
	tasks, ok := parsed["tasks"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'tasks' key in saved config")
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 task entries, got %d", len(tasks))
	}
}
