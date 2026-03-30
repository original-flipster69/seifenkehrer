package task

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type taskDef struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Globs       []string `yaml:"globs"`
	Exclude     []string `yaml:"exclude,omitempty"`
	KeepNewest  int      `yaml:"keep_newest,omitempty"`
	Interval    string   `yaml:"interval,omitempty"`
}

type Task struct {
	Name        string
	Description string
	Interval    string
	Disabled    bool
}

func (t Task) EffectiveInterval() (time.Duration, error) {
	if t.Interval == "" {
		return 0, nil
	}
	return time.ParseDuration(t.Interval)
}

type Result struct {
	Name    string
	Paths   []string
	Skipped string
	Error   error
}

type IntervalChecker interface {
	ShouldSkip(taskName string, interval time.Duration) bool
}

type loaded struct {
	task Task
	def  taskDef
}

func LoadAll(tasksDir string) ([]Task, []error) {
	all, errs := loadAll(tasksDir)
	tasks := make([]Task, len(all))
	for i, l := range all {
		tasks[i] = l.task
	}
	return tasks, errs
}

func Resolve(tasksDir string, checker IntervalChecker) ([]Result, []error) {
	all, errs := loadAll(tasksDir)

	var results []Result
	for _, l := range all {
		if l.task.Disabled {
			results = append(results, Result{Name: l.task.Name, Skipped: "disabled"})
			continue
		}

		interval, err := l.task.EffectiveInterval()
		if err != nil {
			results = append(results, Result{Name: l.task.Name, Error: fmt.Errorf("invalid interval %q: %w", l.task.Interval, err)})
			continue
		}

		if checker != nil && checker.ShouldSkip(l.task.Name, interval) {
			results = append(results, Result{Name: l.task.Name, Skipped: "interval not elapsed"})
			continue
		}

		paths, err := resolveGlobs(l.def.Globs, l.def.Exclude, l.def.KeepNewest)
		results = append(results, Result{
			Name:  l.task.Name,
			Paths: paths,
			Error: err,
		})
	}

	return results, errs
}

func loadAll(tasksDir string) ([]loaded, []error) {
	configDir := filepath.Dir(tasksDir)
	cfg, _ := loadTaskConfig(configDir)

	entries, err := os.ReadDir(tasksDir)
	if err != nil {
		return nil, []error{fmt.Errorf("reading tasks directory %s: %w", tasksDir, err)}
	}

	var all []loaded
	var errs []error
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".yml" && ext != ".yaml" {
			continue
		}

		taskPath := filepath.Join(tasksDir, entry.Name())
		def, err := loadDef(taskPath)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		t := Task{
			Name:        def.Name,
			Description: def.Description,
			Interval:    def.Interval,
		}

		if cfg != nil {
			if cfg.isDisabled(t.Name) {
				t.Disabled = true
			}
			if override := cfg.intervalOverride(t.Name); override != "" {
				t.Interval = override
			}
		}

		all = append(all, loaded{task: t, def: def})
	}

	return all, errs
}

func loadDef(path string) (taskDef, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return taskDef{}, fmt.Errorf("reading %s: %w", path, err)
	}

	var def taskDef
	if err := yaml.Unmarshal(data, &def); err != nil {
		return taskDef{}, fmt.Errorf("parsing %s: %w", path, err)
	}

	if def.Name == "" {
		def.Name = filepath.Base(path)
	}

	if def.Description == "" {
		return taskDef{}, fmt.Errorf("parsing %s: missing required field 'description'", path)
	}

	return def, nil
}

func resolveGlobs(globs []string, exclude []string, keepNewest int) ([]string, error) {
	seen := make(map[string]bool)
	var paths []string

	excludeSet := make(map[string]bool)
	for _, e := range exclude {
		excludeSet[e] = true
	}

	for _, pattern := range globs {
		expanded := expandHome(pattern)
		matches, err := filepath.Glob(expanded)
		if err != nil {
			return nil, fmt.Errorf("invalid glob %q: %w", pattern, err)
		}

		for _, m := range matches {
			abs, err := filepath.Abs(m)
			if err != nil {
				continue
			}
			if seen[abs] {
				continue
			}
			seen[abs] = true

			if excludeSet[filepath.Base(abs)] {
				continue
			}

			paths = append(paths, abs)
		}
	}

	if keepNewest > 0 {
		var statable []string
		for _, p := range paths {
			if _, err := os.Lstat(p); err == nil {
				statable = append(statable, p)
			}
		}
		if len(statable) <= keepNewest {
			return nil, nil
		}
		sort.Slice(statable, func(i, j int) bool {
			return modTime(statable[i]) > modTime(statable[j])
		})
		paths = statable[keepNewest:]
	}

	return paths, nil
}

func modTime(path string) int64 {
	info, err := os.Lstat(path)
	if err != nil {
		return 0
	}
	return info.ModTime().UnixNano()
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
