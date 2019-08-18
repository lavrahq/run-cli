package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

// PreRun holds the standard things that should run before each command
// is executed.
func PreRun(cmd *cobra.Command, args []string) {
	fmt.Println()
}

// PostRun holds the standard things that should run after each command
// is executed.
func PostRun(cmd *cobra.Command, args []string) {
	fmt.Println()
}

// ExitWithMessage allows exiting the command execution with a specific
// message
func ExitWithMessage(message string) {
	fmt.Println(message)
}

// ExitWithError allows exiting the command execution with a specific
// error
func ExitWithError(err error) {
	fmt.Println(err.Error())
}
