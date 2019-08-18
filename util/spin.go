package util

import (
	"fmt"
	"time"

	spinner "github.com/briandowns/spinner"
)

// SpinnerConfig holds spinner configurations.
type SpinnerConfig struct {
	Message string
	Spinner *spinner.Spinner
}

// Spin starts the Spinner.
func Spin(message string) SpinnerConfig {
	spin := SpinnerConfig{
		Message: message,
	}

	processString := fmt.Sprintf(message)
	spin.Spinner = spinner.New(spinner.CharSets[4], 100*time.Millisecond)
	spin.Spinner.Suffix = processString

	return spin
}

// Done completes the Spinner.
func (spin SpinnerConfig) Done() {
	spin.Spinner.Stop()

	fmt.Printf("✓ %s\n", spin.Message)
}

// Failed completes the spinner with an error.
func (spin SpinnerConfig) Failed(err error) {
	spin.Spinner.Stop()

	fmt.Printf("! %s: %s\n", spin.Message, err.Error())
}
