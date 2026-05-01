package execute

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/original-flipster69/seifenkehrer/internal/storage"
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

func ValidatePath(p string) error {
	abs, err := filepath.Abs(p)
	if err != nil {
		return fmt.Errorf("cannot resolve %s: %w", p, err)
	}

	resolved, err := filepath.EvalSymlinks(abs)
	if err == nil {
		abs = resolved
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

func deletePaths(paths []string) Report {
	report := Report{
		Errors: make(map[string]error),
	}

	for _, p := range paths {
		if err := ValidatePath(p); err != nil {
			report.Errors[p] = err
			continue
		}

		info, err := os.Lstat(p)
		if err != nil {
			report.Errors[p] = err
			continue
		}

		if info.IsDir() {
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
