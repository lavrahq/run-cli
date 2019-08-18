package cmdutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// PreRun holds the standard things that should run before each command
// is executed.
func PreRun(cmd *cobra.Command, args []string) {
	fmt.Println()

	fmt.Printf("%s\n", aurora.Yellow(strings.ReplaceAll(cmd.CommandPath(), "lavra", "")))
	fmt.Println()
}

// PostRun holds the standard things that should run after each command
// is executed.
func PostRun(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Printf(" %s\n", aurora.Green("success."))
	fmt.Println()
}

// ExitWithMessage allows exiting the command execution with a specific
// message
func ExitWithMessage(message string) {
	fmt.Println()
	fmt.Printf(" failed: %s\n", aurora.Red(message))
	fmt.Println()

	os.Exit(1)
}

// ExitWithError allows exiting the command execution with a specific
// error
func ExitWithError(err error) {
	fmt.Println()
	fmt.Printf(" failed: %s\n", aurora.Red(err.Error()))
	fmt.Println()

	os.Exit(1)
}
