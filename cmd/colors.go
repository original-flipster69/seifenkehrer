package cmd

import (
	"fmt"
	"os"
)

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
