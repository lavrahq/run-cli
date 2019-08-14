// Copyright © 2019 Scott Plunkett <plunkets@aeoss.io>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/lavrahq/cli/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// initialize tabwriter
		w := new(tabwriter.Writer)

		// io, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)
		defer w.Flush()

		fmt.Fprintln(w)
		fmt.Fprintln(w, "Version Summary:")
		fmt.Fprintf(w, "\n %s\t%s", "Commit:", version.GitCommit)
		fmt.Fprintf(w, "\n %s\t%s", "Branch:", version.GitBranch)
		fmt.Fprintf(w, "\n %s\t%s", "Date:", version.BuildDate)
		fmt.Fprintf(w, "\n %s\t%s\n", "Version:", version.Version)
		fmt.Fprintln(w)

		if version.IsDevelopment() {
			fmt.Println("You are running the Lavra CLI in Development mode. Catch those 🐛 .")
			fmt.Println()

			return
		}

		if version.IsLatest() {
			fmt.Println("You are currently using the latest version of the Lavra CLI. 👏")
			fmt.Println()

			return
		}

		fmt.Println("There is an update available! 🚀")
		fmt.Println("You are currently running", version.Version, ". Version", version.LatestVersion(), "is now available.")
		fmt.Println()
		fmt.Println("Use `lavra update` to update to the latest version. If auto-updating is enabled, we'll periodically run this for you.")
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
