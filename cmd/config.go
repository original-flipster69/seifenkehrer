package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/seifenkehrer/seifenkehrer/internal/task"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure cleanup tasks",
}

var enableCmd = &cobra.Command{
	Use:   "enable <task-name>",
	Short: "Enable a disabled task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := task.Enable(filepath.Dir(tasksDir), args[0]); err != nil {
			return err
		}
		printSuccess("Enabled task %s", orange(args[0]))
		return nil
	},
}

var disableCmd = &cobra.Command{
	Use:   "disable <task-name>",
	Short: "Disable a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := task.Disable(filepath.Dir(tasksDir), args[0]); err != nil {
			return err
		}
		printSuccess("Disabled task %s", orange(args[0]))
		return nil
	},
}

var intervalCmd = &cobra.Command{
	Use:   "interval <task-name> <duration>",
	Short: "Set the run interval for a task (e.g. 24h, 168h)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := time.ParseDuration(args[1]); err != nil {
			return fmt.Errorf("invalid duration %q: %w", args[1], err)
		}
		if err := task.SetInterval(filepath.Dir(tasksDir), args[0], args[1]); err != nil {
			return err
		}
		printSuccess("Set interval for %s to %s", orange(args[0]), gold(args[1]))
		return nil
	},
}

func init() {
	configCmd.AddCommand(enableCmd)
	configCmd.AddCommand(disableCmd)
	configCmd.AddCommand(intervalCmd)
	rootCmd.AddCommand(configCmd)
}
