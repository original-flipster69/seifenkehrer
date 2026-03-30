package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/seifenkehrer/seifenkehrer/internal/execute"
	"github.com/seifenkehrer/seifenkehrer/internal/storage"
	"github.com/seifenkehrer/seifenkehrer/internal/task"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Run cleanup tasks and delete matched files/folders",
	RunE:  run,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

type taskGroup struct {
	name  string
	paths []string
	sizes []int64
}

func run(cmd *cobra.Command, args []string) error {
	exec, err := execute.New(filepath.Dir(tasksDir), storage.FileStorage{})
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	res, errs := task.Resolve(tasksDir, exec)
	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "  [error] %v\n", e)
	}

	var groups []taskGroup
	for _, r := range res {
		if r.Error != nil {
			fmt.Fprintf(os.Stderr, "  [error] %s: %v\n", r.Name, r.Error)
			continue
		}
		if len(r.Paths) == 0 {
			continue
		}

		g := taskGroup{name: r.Name}
		for _, p := range r.Paths {
			g.paths = append(g.paths, p)
			g.sizes = append(g.sizes, pathSize(p))
		}
		groups = append(groups, g)
	}

	if len(groups) == 0 {
		fmt.Println("\n  Nothing to clean up.")
		return nil
	}

	reader := bufio.NewReader(os.Stdin)

	for _, g := range groups {
		totalSize := int64(0)
		for _, s := range g.sizes {
			totalSize += s
		}

		fmt.Printf("\n%s\n", sectionHeader(g.name, formatSize(totalSize)))

		for i, p := range g.paths {
			fmt.Printf("    %s (%s)\n", p, formatSize(g.sizes[i]))
		}

		fmt.Print("\n  Delete? [y]es all / [n]o skip / [i]ndividual: ")
		answer := readLine(reader)

		switch answer {
		case "y", "yes":
			report := exec.Delete(g.name, g.paths)
			printReport(report)
		case "i", "individual":
			selectIndividual(reader, exec, g)
		default:
			fmt.Println("  Skipped.")
		}
	}

	return nil
}

func sectionHeader(title string, detail string) string {
	if detail != "" {
		full := fmt.Sprintf("%s (%s)", title, detail)
		return fmt.Sprintf("  %s\n  %s", full, strings.Repeat("─", utf8.RuneCountInString(full)))
	}
	return fmt.Sprintf("  %s\n  %s", title, strings.Repeat("─", utf8.RuneCountInString(title)))
}

func selectIndividual(reader *bufio.Reader, executor *execute.Executor, g taskGroup) {
	var anyDeleted bool
	for i, p := range g.paths {
		fmt.Printf("    Delete %s (%s)? [y/N]: ", p, formatSize(g.sizes[i]))
		answer := readLine(reader)
		if answer == "y" || answer == "yes" {
			report := executor.DeleteOne(p)
			if len(report.Deleted) > 0 {
				anyDeleted = true
				fmt.Printf("    Deleted.\n")
			}
			for ep, e := range report.Errors {
				fmt.Fprintf(os.Stderr, "    [error] %s: %v\n", ep, e)
			}
		}
	}
	if anyDeleted {
		executor.RecordRun(g.name)
	}
}

func printReport(report execute.Report) {
	fmt.Printf("  Deleted %d path(s).\n", len(report.Deleted))
	for p, e := range report.Errors {
		fmt.Fprintf(os.Stderr, "  [error] %s: %v\n", p, e)
	}
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(line))
}

func pathSize(p string) int64 {
	info, err := os.Lstat(p)
	if err != nil {
		return 0
	}
	if !info.IsDir() {
		return info.Size()
	}

	var total int64
	filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if !d.IsDir() {
			fi, err := d.Info()
			if err != nil {
				return nil
			}
			total += fi.Size()
		}
		return nil
	})
	return total
}

func formatSize(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)

	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
