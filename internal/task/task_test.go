package task

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadDefRequiresDescription(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "no-desc.yaml", `
name: Test Task
globs:
  - /tmp/test/*
`)

	_, err := loadDef(filepath.Join(dir, "no-desc.yaml"))
	if err == nil {
		t.Fatal("expected error for missing description")
	}
}

func TestLoadDefDefaults(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "test.yaml", `
description: A test task
globs:
  - /tmp/test/*
`)

	def, err := loadDef(filepath.Join(dir, "test.yaml"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if def.Name != "test.yaml" {
		t.Errorf("expected name to default to filename, got %q", def.Name)
	}
}

func TestLoadAll(t *testing.T) {
	dir := t.TempDir()
	tasksDir := filepath.Join(dir, "tasks")
	os.MkdirAll(tasksDir, 0755)

	writeFile(t, tasksDir, "a.yaml", `
name: Task A
description: First task
globs:
  - /tmp/a/*
`)
	writeFile(t, tasksDir, "b.yml", `
name: Task B
description: Second task
globs:
  - /tmp/b/*
`)
	writeFile(t, tasksDir, "skip.txt", "not a task")

	tasks, errs := LoadAll(tasksDir)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestLoadAllAppliesOverrides(t *testing.T) {
	dir := t.TempDir()
	tasksDir := filepath.Join(dir, "tasks")
	os.MkdirAll(tasksDir, 0755)

	writeFile(t, tasksDir, "task.yaml", `
name: My Task
description: A task
interval: 1h
globs:
  - /tmp/*
`)
	writeFile(t, dir, "config.json", `{
  "tasks": {
    "My Task": {
      "disabled": true,
      "interval": "48h"
    }
  }
}`)

	tasks, errs := LoadAll(tasksDir)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
	if !tasks[0].Disabled {
		t.Error("expected task to be disabled")
	}
	if tasks[0].Interval != "48h" {
		t.Errorf("expected interval 48h, got %q", tasks[0].Interval)
	}
}

func TestResolveGlobsExclude(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "keep"), 0755)
	os.MkdirAll(filepath.Join(dir, "delete"), 0755)

	paths, err := resolveGlobs([]string{filepath.Join(dir, "*")}, []string{"keep"}, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(paths))
	}
	if filepath.Base(paths[0]) != "delete" {
		t.Errorf("expected 'delete', got %q", filepath.Base(paths[0]))
	}
}

func TestResolveGlobsKeepNewest(t *testing.T) {
	dir := t.TempDir()

	old := filepath.Join(dir, "old")
	os.MkdirAll(old, 0755)
	os.Chtimes(old, time.Now().Add(-2*time.Hour), time.Now().Add(-2*time.Hour))

	mid := filepath.Join(dir, "mid")
	os.MkdirAll(mid, 0755)
	os.Chtimes(mid, time.Now().Add(-1*time.Hour), time.Now().Add(-1*time.Hour))

	newest := filepath.Join(dir, "newest")
	os.MkdirAll(newest, 0755)

	paths, err := resolveGlobs([]string{filepath.Join(dir, "*")}, nil, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths for deletion, got %d", len(paths))
	}
	for _, p := range paths {
		if filepath.Base(p) == "newest" {
			t.Error("newest should have been kept")
		}
	}
}

func TestResolveGlobsKeepNewestReturnsNilWhenAllKept(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "only"), 0755)

	paths, err := resolveGlobs([]string{filepath.Join(dir, "*")}, nil, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if paths != nil {
		t.Errorf("expected nil paths when all kept, got %v", paths)
	}
}

func TestEffectiveInterval(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
		wantErr  bool
	}{
		{"", 0, false},
		{"24h", 24 * time.Hour, false},
		{"30m", 30 * time.Minute, false},
		{"garbage", 0, true},
	}

	for _, tt := range tests {
		tk := Task{Interval: tt.input}
		d, err := tk.EffectiveInterval()
		if tt.wantErr && err == nil {
			t.Errorf("EffectiveInterval(%q): expected error", tt.input)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("EffectiveInterval(%q): unexpected error: %v", tt.input, err)
		}
		if d != tt.expected {
			t.Errorf("EffectiveInterval(%q) = %v, want %v", tt.input, d, tt.expected)
		}
	}
}

type noSkip struct{}

func (noSkip) ShouldSkip(string, time.Duration) bool { return false }

type skipAll struct{}

func (skipAll) ShouldSkip(string, time.Duration) bool { return true }

func TestResolveSkipsDisabled(t *testing.T) {
	dir := t.TempDir()
	tasksDir := filepath.Join(dir, "tasks")
	os.MkdirAll(tasksDir, 0755)

	writeFile(t, tasksDir, "active.yaml", `
name: active
description: An active task
globs:
  - /tmp/nonexistent/*
`)
	writeFile(t, tasksDir, "disabled.yaml", `
name: disabled
description: A disabled task
globs:
  - /tmp/nonexistent/*
`)
	writeFile(t, dir, "config.json", `{
  "tasks": {
    "disabled": {
      "disabled": true
    }
  }
}`)

	results, _ := Resolve(tasksDir, noSkip{})
	for _, r := range results {
		if r.Name == "disabled" {
			t.Error("expected disabled task to be skipped")
		}
	}
}

func TestResolveSkipsOnInterval(t *testing.T) {
	dir := t.TempDir()
	tasksDir := filepath.Join(dir, "tasks")
	os.MkdirAll(tasksDir, 0755)

	writeFile(t, tasksDir, "task.yaml", `
name: Skippable
description: A task with interval
interval: 24h
globs:
  - /tmp/nonexistent/*
`)

	results, _ := Resolve(tasksDir, skipAll{})
	for _, r := range results {
		if r.Name == "Skippable" {
			t.Error("expected task to be skipped")
		}
	}
}

func TestResolveRunsWhenNotSkipped(t *testing.T) {
	dir := t.TempDir()
	tasksDir := filepath.Join(dir, "tasks")
	os.MkdirAll(tasksDir, 0755)

	writeFile(t, tasksDir, "task.yaml", `
name: Runnable
description: A task with interval
interval: 24h
globs:
  - /tmp/nonexistent/*
`)

	results, _ := Resolve(tasksDir, noSkip{})
	found := false
	for _, r := range results {
		if r.Name == "Runnable" {
			found = true
		}
	}
	if !found {
		t.Error("expected task to be included")
	}
}

func TestExpandHome(t *testing.T) {
	home, _ := os.UserHomeDir()
	result := expandHome("~/test")
	expected := filepath.Join(home, "test")
	if result != expected {
		t.Errorf("expandHome(~/test) = %q, want %q", result, expected)
	}

	abs := "/absolute/path"
	if expandHome(abs) != abs {
		t.Error("expandHome should not modify absolute paths")
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
