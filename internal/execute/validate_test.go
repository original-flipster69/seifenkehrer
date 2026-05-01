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
		if err := ValidatePath(p); err == nil {
			t.Errorf("ValidatePath(%q): expected error for protected path", p)
		}
	}
}

func TestValidatePathRoot(t *testing.T) {
	if err := ValidatePath("/"); err == nil {
		t.Error("ValidatePath(/): expected error for root directory")
	}
}

func TestValidatePathHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory")
	}
	if err := ValidatePath(home); err == nil {
		t.Errorf("ValidatePath(%q): expected error for home directory", home)
	}
}

func TestValidatePathAllowsSubdirOfHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home directory")
	}
	sub := filepath.Join(home, "some-temp-subdir-test")
	if err := ValidatePath(sub); err != nil {
		t.Errorf("ValidatePath(%q): unexpected error: %v", sub, err)
	}
}

func TestValidatePathAllowsTmpDir(t *testing.T) {
	dir := tmpDir(t)
	if err := ValidatePath(dir); err != nil {
		t.Errorf("ValidatePath(%q): unexpected error: %v", dir, err)
	}
}

func TestValidatePathRejectsSymlinkTraversal(t *testing.T) {
	dir := tmpDir(t)
	link := filepath.Join(dir, "escape")
	os.Symlink("/usr", link)

	traversed := filepath.Join(link, "local")
	if err := ValidatePath(traversed); err == nil {
		t.Errorf("ValidatePath(%q): expected error for symlink traversal into /usr", traversed)
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

func TestDeleteDirWithSymlinkPreservesExternalTarget(t *testing.T) {
	outer := tmpDir(t)
	target := filepath.Join(outer, "outside.txt")
	if err := os.WriteFile(target, []byte("keep me"), 0644); err != nil {
		t.Fatal(err)
	}

	victim := filepath.Join(outer, "victim")
	if err := os.MkdirAll(victim, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(victim, "link.txt")); err != nil {
		t.Fatal(err)
	}

	report := deletePaths([]string{victim})
	if len(report.Deleted) != 1 {
		t.Fatalf("expected 1 deleted, got %d (errors: %v)", len(report.Deleted), report.Errors)
	}
	if _, err := os.Stat(victim); !os.IsNotExist(err) {
		t.Error("victim directory should have been deleted")
	}
	if _, err := os.Stat(target); err != nil {
		t.Errorf("external symlink target should still exist: %v", err)
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
