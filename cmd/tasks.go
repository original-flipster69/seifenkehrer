package cmd

import (
	"fmt"

	"github.com/original-flipster69/seifenkehrer/internal/task"
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
		printError("%v", e)
	}

	if len(tasks) == 0 {
		printInfo("No tasks found.")
		return nil
	}

	fmt.Printf("\n  %s %s\n\n", orange("Installed tasks"), dim(fmt.Sprintf("(%d)", len(tasks))))
	for _, t := range tasks {
		status := green("enabled")
		if t.Disabled {
			status = red("disabled")
		}

		detail := status
		if t.Interval != "" {
			detail += dim(", every ") + gold(t.Interval)
		}

		fmt.Printf("  %s %s\n", orange(t.Name), dim("[") + detail + dim("]"))
		fmt.Printf("    %s\n\n", grey(t.Description))
	}

	return nil
}
