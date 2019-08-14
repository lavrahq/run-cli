package version

import (
	"fmt"
	"os"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

// BuildDate provides the date that the CLI was built.
var BuildDate string

// GitCommit provides the current commit at the time of build.
var GitCommit string

// GitBranch provides the branch that was checked out at the time of build.
var GitBranch string

// GitState provides the state of the Git project at the time of build.
var GitState string

// GitSummary provides the summary of the Git project at the time of build.
var GitSummary string

// Version provides the human-friendly version of the build.
var Version string

// IsDevelopment returns true if the Version is not set.
func IsDevelopment() bool {
	if Version == "" {
		return true
	}

	return false
}

// IsLatest returns true if the current version is the latest and false if it is not.
func IsLatest() bool {
	if Version == "" {
		return true
	}

	latest, found, err := selfupdate.DetectLatest("lavrahq/cli")
	if err != nil {
		return true
	}

	currentVersion := semver.MustParse(Version)
	if !found || latest.Version.LTE(currentVersion) {
		return true
	}

	return false
}

// LatestVersion determines the latest version and returns it as a string.
func LatestVersion() string {
	latest, found, err := selfupdate.DetectLatest("lavrahq/cli")
	if err != nil {
		return ""
	}

	if found {
		return latest.Version.String()
	}

	return ""
}

// Update tries to update the Lavra CLI and provides direct output to the terminal.
func Update() {
	cli, err := os.Executable()
	if err != nil {
		fmt.Println("Could not locate the executable:", err.Error())

		return
	}

	latest, found, err := selfupdate.DetectLatest("lavrahq/cli")
	if err != nil {
		fmt.Println("Could not update. Encountered an error:", err.Error())

		return
	}

	if found {
		if err := selfupdate.UpdateTo(latest.AssetURL, cli); err != nil {
			fmt.Println("Update failed:", err.Error())

			return
		}

		fmt.Println()
		fmt.Println("The update has completed successfully, you are now running", latest.Version)
		fmt.Println()
		fmt.Println("Here's what changed in this version...")
		fmt.Println(latest.ReleaseNotes)
		fmt.Println()

		return
	}

	fmt.Println("Could not find a suitable update. This is probably an issue on our end.")
}
