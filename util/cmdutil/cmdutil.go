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

// ExitWithMessageStep allows exiting the command execution with a specific
// message and step.
func ExitWithMessageStep(message string, step string) {
	fmt.Println()
	fmt.Printf(" failed: %s\n when: %s\n", aurora.Red(message), step)
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

// ExitWithErrorStep allows exiting the command execution with a specific
// error
func ExitWithErrorStep(err error, step string) {
	fmt.Println()
	fmt.Printf(" failed: %s\n when %s\n", aurora.Red(err.Error()), step)
	fmt.Println()

	os.Exit(1)
}

// CheckCommandError checks if the error is nil, if not, it returns the error
// and stops command execution.
func CheckCommandError(err error, step string, message ...string) {
	if err != nil {
		if len(message) > 0 {
			ExitWithMessageStep(strings.Join(message, " "), step)

			return
		}

		ExitWithErrorStep(err, step)
	}

	return
}
