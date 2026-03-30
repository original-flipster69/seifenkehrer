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
	if e.ShouldSkip("anything", time.Hour) {
		t.Error("empty state should not skip any task")
	}
}

func TestRecordAndSkip(t *testing.T) {
	dir := t.TempDir()
	e, _ := New(dir, storage.FileStorage{})

	e.RecordRun("test-task")

	if !e.ShouldSkip("test-task", 24*time.Hour) {
		t.Error("should skip task that just ran with 24h interval")
	}
	if e.ShouldSkip("test-task", 0) {
		t.Error("should not skip when interval is 0")
	}
	if e.ShouldSkip("other-task", 24*time.Hour) {
		t.Error("should not skip unrecorded task")
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
	if !e2.ShouldSkip("persisted-task", 24*time.Hour) {
		t.Error("persisted task should still be tracked after reload")
	}
}

func TestShouldNotSkipExpiredInterval(t *testing.T) {
	dir := t.TempDir()
	e, _ := New(dir, storage.FileStorage{})

	e.state.entries["old-task"] = stateEntry{LastRun: time.Now().Add(-48 * time.Hour)}

	if e.ShouldSkip("old-task", 24*time.Hour) {
		t.Error("should not skip task whose interval has expired")
	}
}
