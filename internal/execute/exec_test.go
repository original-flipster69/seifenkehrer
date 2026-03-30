package execute

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteFile(t *testing.T) {
	dir := tmpDir(t)
	f := filepath.Join(dir, "deleteme.txt")
	os.WriteFile(f, []byte("data"), 0644)

	report := deletePaths([]string{f})
	if len(report.Deleted) != 1 {
		t.Fatalf("expected 1 deleted, got %d", len(report.Deleted))
	}
	if len(report.Errors) != 0 {
		t.Fatalf("expected 0 errors, got %d", len(report.Errors))
	}
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Error("file should have been deleted")
	}
}

func TestDeleteDirectory(t *testing.T) {
	dir := tmpDir(t)
	sub := filepath.Join(dir, "subdir")
	os.MkdirAll(sub, 0755)
	os.WriteFile(filepath.Join(sub, "file.txt"), []byte("data"), 0644)

	report := deletePaths([]string{sub})
	if len(report.Deleted) != 1 {
		t.Fatalf("expected 1 deleted, got %d", len(report.Deleted))
	}
	if _, err := os.Stat(sub); !os.IsNotExist(err) {
		t.Error("directory should have been deleted")
	}
}

func TestDeleteNonexistent(t *testing.T) {
	report := deletePaths([]string{"/nonexistent/path/12345"})
	if len(report.Deleted) != 0 {
		t.Error("expected 0 deleted")
	}
	if len(report.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(report.Errors))
	}
}

func TestDeleteEmpty(t *testing.T) {
	report := deletePaths(nil)
	if len(report.Deleted) != 0 {
		t.Error("expected 0 deleted")
	}
	if len(report.Errors) != 0 {
		t.Error("expected 0 errors")
	}
}
