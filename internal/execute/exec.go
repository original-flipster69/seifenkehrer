package execute

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/seifenkehrer/seifenkehrer/internal/storage"
)

var protectedPrefixes = []string{
	"/System",
	"/usr",
	"/bin",
	"/sbin",
	"/etc",
	"/var",
	"/private/etc",
	"/private/var",
}

type Executor struct {
	state *state
}

type Report struct {
	Deleted []string
	Errors  map[string]error
}

func New(configDir string, store storage.Storage) (*Executor, error) {
	st, err := load(configDir, store)
	if err != nil {
		return nil, err
	}
	return &Executor{state: st}, nil
}

func (e *Executor) LastRun(task string) *time.Time {
	return e.state.lastRun(task)
}

func (e *Executor) Delete(task string, paths []string) Report {
	report := deletePaths(paths)
	if len(report.Deleted) > 0 {
		e.state.RecordRun(task)
	}
	return report
}

func (e *Executor) DeleteOne(path string) Report {
	return deletePaths([]string{path})
}

func (e *Executor) RecordRun(task string) {
	e.state.RecordRun(task)
}

func deletePaths(paths []string) Report {
	report := Report{
		Errors: make(map[string]error),
	}

	for _, p := range paths {
		if err := validatePath(p); err != nil {
			report.Errors[p] = err
			continue
		}

		info, err := os.Lstat(p)
		if err != nil {
			report.Errors[p] = err
			continue
		}

		if info.Mode()&os.ModeSymlink != 0 {
			err = os.Remove(p)
		} else if info.IsDir() {
			if containsSymlinks(p) {
				report.Errors[p] = fmt.Errorf("refusing to delete %s: contains symlinks", p)
				continue
			}
			err = os.RemoveAll(p)
		} else {
			err = os.Remove(p)
		}

		if err != nil {
			report.Errors[p] = fmt.Errorf("deleting %s: %w", p, err)
		} else {
			report.Deleted = append(report.Deleted, p)
		}
	}

	return report
}

func validatePath(p string) error {
	abs, err := filepath.Abs(p)
	if err != nil {
		return fmt.Errorf("cannot resolve %s: %w", p, err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot determine home directory: %w", err)
	}

	for _, prefix := range protectedPrefixes {
		if strings.HasPrefix(abs, prefix+"/") || abs == prefix {
			return fmt.Errorf("refusing to delete protected path: %s", abs)
		}
	}

	if abs == home {
		return fmt.Errorf("refusing to delete home directory")
	}

	if abs == "/" {
		return fmt.Errorf("refusing to delete root directory")
	}

	return nil
}

func containsSymlinks(dir string) bool {
	found := false
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			found = true
			return filepath.SkipAll
		}
		return nil
	})
	return found
}
