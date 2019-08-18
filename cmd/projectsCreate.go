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

package cmd

import (
	"github.com/AlecAivazis/survey"
	"github.com/lavrahq/cli/util"
	"github.com/lavrahq/cli/util/cmdutil"
	"github.com/lavrahq/cli/util/dir"
	"github.com/lavrahq/cli/util/tmpl"
	"github.com/spf13/cobra"
)

// Stores the --track, -t flag
var flagNoTrack bool

// Stores the --template, -a flag
var flagTemplate string

// The survey questions to ask when creating the project.
var surveyQuestions = []*survey.Question{
	{
		Name:      "name",
		Prompt:    &survey.Input{Message: "Project Name"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
}

// projectsCreateCmd represents the projectsCreate command
var projectsCreateCmd = &cobra.Command{
	Use:   "create <dir=.>",
	Short: "Creates a new project at the specified directory. Defaults to current dir.",
	Long: `The create command initializes a new project in the given directory, defaulting
to the current directory if a directory is not provided. By default, the new project is
tracked and managed via the CLI.`,
	Args:    cobra.MaximumNArgs(1),
	PreRun:  cmdutil.PreRun,
	PostRun: cmdutil.PostRun,
	Run: func(cmd *cobra.Command, args []string) {
		var rawDir = "."
		if len(args) != 0 {
			rawDir = args[0]
		}

		setupProjDir := util.Spin("Configuring project directory")
		projDir, _ := dir.Make(rawDir)
		if projDir.IsProject() {
			cmdutil.ExitWithMessage("You cannot create a new project within a project root.")

			return
		}
		setupProjDir.Done()

		configureTemplate := util.Spin("Configuring project template")
		template := tmpl.Make(projDir, flagTemplate)
		configureTemplate.Done()

		fetchingTemplate := util.Spin("Fetching project template")
		isFetched, err := template.EnsureTemplateIsFetched()
		if !isFetched {
			cmdutil.ExitWithError(err)

			return
		}

		fetchingTemplate.Done()

		// projDir.TemplateFrom("/users/scott/")
	},
}

func init() {
	projectsCmd.AddCommand(projectsCreateCmd)

	// Allows disabling tracking for the new project
	projectsCreateCmd.Flags().BoolVarP(&flagNoTrack, "no-track", "t", true, "Disable tracking for the new project")

	// Allows specifying the template for the new project.
	projectsCreateCmd.Flags().StringVarP(&flagTemplate, "template", "m", "empty", "Specifies the template to use for the new project")
}
