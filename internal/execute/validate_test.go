package execute

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("/tmp", "sk-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestValidatePathProtectedPrefixes(t *testing.T) {
	protected := []string{
		"/System",
		"/System/Library",
		"/usr",
		"/usr/local/bin",
		"/bin",
		"/sbin",
		"/etc",
		"/var",
		"/private/etc",
		"/private/var",
	}

	for _, p := range protected {
		if err := validatePath(p); err == nil {
			t.Errorf("validatePath(%q): expected error for protected path", p)
		}
	}
}

func TestValidatePathRoot(t *testing.T) {
	if err := validatePath("/"); err == nil {
		t.Error("validatePath(/): expected error for root directory")
	}
}

func TestValidatePathHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory")
	}
	if err := validatePath(home); err == nil {
		t.Errorf("validatePath(%q): expected error for home directory", home)
	}
}

func TestValidatePathAllowsSubdirOfHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory")
	}
	sub := filepath.Join(home, "some-temp-subdir-test")
	if err := validatePath(sub); err != nil {
		t.Errorf("validatePath(%q): unexpected error: %v", sub, err)
	}
}

func TestValidatePathAllowsTmpDir(t *testing.T) {
	dir := tmpDir(t)
	if err := validatePath(dir); err != nil {
		t.Errorf("validatePath(%q): unexpected error: %v", dir, err)
	}
}

func TestContainsSymlinksNoSymlinks(t *testing.T) {
	dir := tmpDir(t)
	os.MkdirAll(filepath.Join(dir, "sub"), 0755)
	os.WriteFile(filepath.Join(dir, "sub", "file.txt"), []byte("data"), 0644)

	if containsSymlinks(dir) {
		t.Error("expected no symlinks in directory")
	}
}

func TestContainsSymlinksWithSymlink(t *testing.T) {
	dir := tmpDir(t)
	target := filepath.Join(dir, "target.txt")
	os.WriteFile(target, []byte("data"), 0644)

	link := filepath.Join(dir, "link.txt")
	os.Symlink(target, link)

	if !containsSymlinks(dir) {
		t.Error("expected symlinks to be detected")
	}
}

func TestContainsSymlinksNested(t *testing.T) {
	dir := tmpDir(t)
	sub := filepath.Join(dir, "a", "b")
	os.MkdirAll(sub, 0755)

	target := filepath.Join(dir, "target.txt")
	os.WriteFile(target, []byte("data"), 0644)
	os.Symlink(target, filepath.Join(sub, "link.txt"))

	if !containsSymlinks(dir) {
		t.Error("expected nested symlink to be detected")
	}
}

func TestContainsSymlinksEmpty(t *testing.T) {
	dir := tmpDir(t)
	if containsSymlinks(dir) {
		t.Error("empty directory should not contain symlinks")
	}
}

func TestDeleteRefusesProtectedPath(t *testing.T) {
	report := deletePaths([]string{"/System/Library"})
	if len(report.Deleted) != 0 {
		t.Error("should not delete protected path")
	}
	if _, ok := report.Errors["/System/Library"]; !ok {
		t.Error("expected error for protected path")
	}
}

func TestDeleteRefusesDirWithSymlinks(t *testing.T) {
	dir := tmpDir(t)
	sub := filepath.Join(dir, "with-links")
	os.MkdirAll(sub, 0755)

	target := filepath.Join(dir, "target.txt")
	os.WriteFile(target, []byte("data"), 0644)
	os.Symlink(target, filepath.Join(sub, "link.txt"))

	report := deletePaths([]string{sub})
	if len(report.Deleted) != 0 {
		t.Error("should not delete directory containing symlinks")
	}
	if _, ok := report.Errors[sub]; !ok {
		t.Error("expected error for directory with symlinks")
	}
}

func TestDeleteSymlinkDirectly(t *testing.T) {
	dir := tmpDir(t)
	target := filepath.Join(dir, "target.txt")
	os.WriteFile(target, []byte("data"), 0644)

	link := filepath.Join(dir, "link.txt")
	os.Symlink(target, link)

	report := deletePaths([]string{link})
	if len(report.Deleted) != 1 {
		t.Fatal("expected symlink to be deleted")
	}

	if _, err := os.Lstat(link); !os.IsNotExist(err) {
		t.Error("symlink should have been removed")
	}
	if _, err := os.Stat(target); err != nil {
		t.Error("symlink target should still exist")
	}
}

func TestDeletePartialFailure(t *testing.T) {
	dir := tmpDir(t)
	good := filepath.Join(dir, "good.txt")
	os.WriteFile(good, []byte("data"), 0644)

	report := deletePaths([]string{good, "/System/bad"})
	if len(report.Deleted) != 1 {
		t.Errorf("expected 1 deleted, got %d", len(report.Deleted))
	}
	if len(report.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(report.Errors))
	}
}
