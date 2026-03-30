package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileStorageLoadMissing(t *testing.T) {
	data, err := FileStorage{}.Load("/nonexistent/path/file.json")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if data != nil {
		t.Error("expected nil data for missing file")
	}
}

func TestFileStorageSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")
	fs := FileStorage{}

	content := []byte(`{"key": "value"}`)
	if err := fs.Save(path, content); err != nil {
		t.Fatalf("save error: %v", err)
	}

	data, err := fs.Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("expected %q, got %q", content, data)
	}
}

func TestFileStorageSaveCreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "deep", "file.json")
	fs := FileStorage{}

	if err := fs.Save(path, []byte("data")); err != nil {
		t.Fatalf("save error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestFileStorageSaveIsAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "atomic.json")
	fs := FileStorage{}

	if err := fs.Save(path, []byte("first")); err != nil {
		t.Fatal(err)
	}
	if err := fs.Save(path, []byte("second")); err != nil {
		t.Fatal(err)
	}

	data, _ := fs.Load(path)
	if string(data) != "second" {
		t.Errorf("expected 'second', got %q", data)
	}

	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.Name() != "atomic.json" {
			t.Errorf("leftover temp file: %s", e.Name())
		}
	}
}
