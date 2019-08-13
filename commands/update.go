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

	"github.com/lavrahq/cli/version"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the CLI to the latest version.",
	Long: `The update command allows you to update the Lavra CLI to the latest
stable version. It is not possible to update beta versions of the Lavra CLI using
this utility.`,
	Run: func(cmd *cobra.Command, args []string) {
		if version.IsDevelopment() {
			fmt.Println()
			fmt.Println("You cannot update the CLI while running in Development mode.")
			fmt.Println()

			return
		}

		if version.IsLatest() {
			fmt.Println()
			fmt.Println("You are already running the latest version. 👏")
			fmt.Println()

			return
		}

		version.Update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
