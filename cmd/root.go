package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var tasksDir string

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "seifenkehrer",
	Version: version,
	Short:   "seifenkehrer - modular macOS cleanup tool",
	Long:    "A modular cleanup tool that runs user-defined tasks to find and delete unnecessary files and folders.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		printBanner()
	},
}

func init() {
	homeDir, _ := os.UserHomeDir()
	defaultTasksDir := filepath.Join(homeDir, ".seifenkehrer", "tasks")
	rootCmd.PersistentFlags().StringVar(&tasksDir, "tasks-dir", defaultTasksDir, "directory containing cleanup task definitions")
}

func Execute() {
	if name := filepath.Base(os.Args[0]); name != "" && name != "." && name != "/" {
		rootCmd.Use = name
	}
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
		"⠀⠀⠀⠀⠀⠀⠀" + o + "⠈⠙⠻⢿⣿⣿⡿⠟⠋⠁" + r + "⠀⠀      ⠀" + p + "Storage Cleanup CLI" + r + "  " + cDim + version + cReset + "\n\n\n"
}

func printBanner() {
	fmt.Print(banner)
}

var (
	cOrange = "\033[38;5;208m"
	cPurple = "\033[38;5;141m"
	cGold   = "\033[38;5;143m"
	cGrey   = "\033[38;5;248m"
	cRed    = "\033[38;5;167m"
	cGreen  = "\033[38;5;114m"
	cDim    = "\033[2m"
	cBold   = "\033[1m"
	cReset  = "\033[0m"
)

func orange(s string) string { return cOrange + s + cReset }
func purple(s string) string { return cPurple + s + cReset }
func gold(s string) string   { return cGold + s + cReset }
func grey(s string) string   { return cGrey + s + cReset }
func red(s string) string    { return cRed + s + cReset }
func green(s string) string  { return cGreen + s + cReset }
func dim(s string) string    { return cDim + s + cReset }
func bold(s string) string   { return cBold + s + cReset }

func printError(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "  %s %s\n", red("✗"), fmt.Sprintf(format, a...))
}

func printSuccess(format string, a ...any) {
	fmt.Printf("  %s %s\n", green("✓"), fmt.Sprintf(format, a...))
}

func printInfo(format string, a ...any) {
	fmt.Printf("  %s %s\n", purple("●"), fmt.Sprintf(format, a...))
}
