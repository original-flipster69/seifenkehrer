package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{1073741824, "1.0 GB"},
		{1610612736, "1.5 GB"},
	}

	for _, tt := range tests {
		got := formatSize(tt.bytes)
		if got != tt.expected {
			t.Errorf("formatSize(%d) = %q, want %q", tt.bytes, got, tt.expected)
		}
	}
}

func TestPathSizeFile(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.txt")
	data := []byte("hello world")
	os.WriteFile(f, data, 0644)

	size := pathSize(f)
	if size != int64(len(data)) {
		t.Errorf("pathSize(%q) = %d, want %d", f, size, len(data))
	}
}

func TestPathSizeDirectory(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.MkdirAll(sub, 0755)
	os.WriteFile(filepath.Join(sub, "a.txt"), []byte("aaaa"), 0644)
	os.WriteFile(filepath.Join(sub, "b.txt"), []byte("bb"), 0644)

	size := pathSize(sub)
	if size != 6 {
		t.Errorf("pathSize(%q) = %d, want 6", sub, size)
	}
}

func TestPathSizeNonexistent(t *testing.T) {
	size := pathSize("/nonexistent/path/12345")
	if size != 0 {
		t.Errorf("pathSize(nonexistent) = %d, want 0", size)
	}
}

func TestPathSizeSkipsSymlinks(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.MkdirAll(sub, 0755)
	os.WriteFile(filepath.Join(sub, "real.txt"), []byte("real"), 0644)

	target := filepath.Join(dir, "external.txt")
	os.WriteFile(target, []byte("big external file"), 0644)
	os.Symlink(target, filepath.Join(sub, "link.txt"))

	size := pathSize(sub)
	if size != 4 {
		t.Errorf("pathSize(%q) = %d, want 4 (only real.txt)", sub, size)
	}
}

func TestSectionHeaderWithDetail(t *testing.T) {
	header := sectionHeader("Title", "100 MB")
	if header == "" {
		t.Error("expected non-empty header")
	}
}

func TestSectionHeaderWithoutDetail(t *testing.T) {
	header := sectionHeader("Title", "")
	if header == "" {
		t.Error("expected non-empty header")
	}
}
