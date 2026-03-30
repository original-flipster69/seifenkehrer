package execute

import (
	"testing"
	"time"

	"github.com/seifenkehrer/seifenkehrer/internal/storage"
)

func TestNewExecutorEmptyState(t *testing.T) {
	dir := t.TempDir()
	e, err := New(dir, storage.FileStorage{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.LastRun("anything") != nil {
		t.Error("empty state should return nil for any task")
	}
}

func TestRecordAndLastRun(t *testing.T) {
	dir := t.TempDir()
	e, _ := New(dir, storage.FileStorage{})

	e.RecordRun("test-task")

	lr := e.LastRun("test-task")
	if lr == nil {
		t.Fatal("expected non-nil last run after recording")
	}
	if time.Since(*lr) > time.Second {
		t.Error("last run should be very recent")
	}
	if e.LastRun("other-task") != nil {
		t.Error("unrecorded task should return nil")
	}
}

func TestSaveAndReload(t *testing.T) {
	dir := t.TempDir()
	e, _ := New(dir, storage.FileStorage{})
	e.RecordRun("persisted-task")

	e2, err := New(dir, storage.FileStorage{})
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	if e2.LastRun("persisted-task") == nil {
		t.Error("persisted task should still be tracked after reload")
	}
}

func TestLastRunExpired(t *testing.T) {
	dir := t.TempDir()
	e, _ := New(dir, storage.FileStorage{})

	past := time.Now().Add(-48 * time.Hour)
	e.state.entries["old-task"] = stateEntry{LastRun: past}

	lr := e.LastRun("old-task")
	if lr == nil {
		t.Fatal("expected non-nil last run")
	}
	if time.Since(*lr) < 47*time.Hour {
		t.Error("last run should be ~48 hours ago")
	}
}
