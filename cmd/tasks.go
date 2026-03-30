package cmd

import (
	"fmt"
	"os"

	"github.com/seifenkehrer/seifenkehrer/internal/task"
	"github.com/spf13/cobra"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List all installed cleanup tasks",
	RunE:  listTasks,
}

func init() {
	rootCmd.AddCommand(tasksCmd)
}

func listTasks(cmd *cobra.Command, args []string) error {
	tasks, errs := task.LoadAll(tasksDir)

	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "  [error] %v\n", e)
	}

	if len(tasks) == 0 {
		fmt.Println("\n  No tasks found.")
		return nil
	}

	fmt.Printf("\n  Installed tasks (%d):\n\n", len(tasks))
	for _, t := range tasks {
		status := "enabled"
		if t.Disabled {
			status = "disabled"
		}

		detail := status
		if t.Interval != "" {
			detail += ", every " + t.Interval
		}

		fmt.Printf("  %s [%s]\n", t.Name, detail)
		fmt.Printf("    %s\n\n", t.Description)
	}

	return nil
}
