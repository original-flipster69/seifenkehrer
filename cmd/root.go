package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var tasksDir string

var rootCmd = &cobra.Command{
	Use:   "sk",
	Short: "seifenkehrer - modular macOS cleanup tool",
	Long:  "A modular cleanup tool that runs user-defined tasks to find and delete unnecessary files and folders.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		printBanner()
	},
}

func init() {
	homeDir, _ := os.UserHomeDir()
	defaultTasksDir := filepath.Join(homeDir, ".sk", "tasks")
	rootCmd.PersistentFlags().StringVar(&tasksDir, "tasks-dir", defaultTasksDir, "directory containing cleanup task definitions")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var banner string

func init() {
	p := "\033[38;5;141m"
	o := "\033[38;5;208m"
	d := "\033[38;5;143m"
	r := "\033[0m"

	banner = "" +
		p + "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣤⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀" + r + "\n" +
		p + "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣯⡀⣹⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀" + r + "\n" +
		p + "⠀⠀⣴⠟⢷⡄⠀⠀⠀⠀⠀⠀⠀⠈⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀" + r + "\n" +
		p + "⠀⠀⢿⣤⣼⠇⠀⠀⠀⠀⠀⠀⢀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀" + r + "\n" +
		p + "⠀⠀⠀⠀⠀⠀⣀⠀⠀⠀⢠⡟⠛⠛⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀" + r + "\n" +
		p + "⠀⠀⠀⠀⠀⠘⠛⠃⠀⠀⠘⣷⣄⣠⡿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀" + r + "\n" +
		p + "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠁⠀⠀⠀⠀" + o + "⣀⣤⣄⡀" + r + "⠀⠀⠀⠀⠀⠀⠀⠀\n" +
		"⠀⠀⠀⠀⠀⠀⠀⠀" + p + "⢠⣴⣦⡄" + o + "⠀⠀⣀⣤⣶⣿⣿⣿⣿⣿⣷⣦⣄⡀" + r + "⠀⠀⠀⠀\n" +
		"⠀⠀⠀⠀⠀⠀⠀⠀" + p + "⢿⣿⣿⠇" + o + "⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠟⢻⣿⡇" + r + "⠀⠀⠀ " + d + "▄▄▄ ▗▞▀▚▖▄ ▗▞▀▀▘▗▞▀▚▖▄▄▄▄  █  ▄ ▗▞▀▚▖▐▌    ▄▄▄ ▗▞▀▚▖ ▄▄▄" + r + " \n" +
		"⠀⠀⠀⠀⠀⠀" + o + "⣠⣤⣤⣭⣥⣴⣿⣿⣿⣿⣿⣿⠿⠛⢉⣡⣴⣾⣿⣿⡇" + r + "⠀⠀⠀" + d + "▀▄▄  ▐▛▀▀▘▄ ▐▌   ▐▛▀▀▘█   █ █▄▀  ▐▛▀▀▘▐▌   █    ▐▛▀▀▘█" + r + "    \n" +
		"⠀⠀⠀⠀⠀" + o + "⣾⣿⣿⣈⠙⠻⣿⣿⡿⠟⠛⣉⣠⣴⣾⣿⣿⣿⣿⣿⣿⠇" + r + "⠀⠀⠀" + d + "▄▄▄▀ ▝▚▄▄▖█ ▐▛▀▘ ▝▚▄▄▖█   █ █ ▀▄ ▝▚▄▄▖▐▛▀▚▖█    ▝▚▄▄▖█" + r + "    \n" +
		"⠀⠀⠀⠀⠀" + o + "⣿⣿⣿⣿⣿⣶⣤⣤⣤⣶⣿⣿⣿⣿⣿⣿⣿⠿⠛⠉" + r + "⠀⠀⠀⠀⠀" + d + "          █ ▐▌              █  █      ▐▌ ▐▌" + r + "               \n" +
		"⠀⠀⠀⠀⠀" + o + "⠻⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⠛⠉" + r + "⠀⠀⠀⠀⠀⠀⠀⠀⠀\n" +
		"⠀⠀⠀⠀⠀⠀⠀" + o + "⠈⠙⠻⢿⣿⣿⡿⠟⠋⠁" + r + "⠀⠀      ⠀" + p + "Storage Cleanup CLI" + r + "\n"
}

func printBanner() {
	fmt.Print(banner)
}
