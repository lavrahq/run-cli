/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/lavrahq/cli/packages/prompt"
	"github.com/lavrahq/cli/util/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var platform = prompt.Question{
	Name: "Platform",
	Type: "Select",
	Options: prompt.QuestionOptions{
		Options: []string{
			"Kubernetes",
			"Docker",
		},
	},
}

var path = prompt.Question{
	Name: "Path",
	Type: "Input",
	Options: prompt.QuestionOptions{
		Message: "What is the path to your API?",
		Default: "/var/run/docker.sock",
		Help:    "The path to the Docker instance.",
	},
	Validate: prompt.QuestionValidation{
		Required: true,
	},
	When: `(Answers.Platform.Value == "Docker")`,
}

var host = prompt.Question{
	Name: "Host",
	Type: "Input",
	Options: prompt.QuestionOptions{
		Message: "Hostname or IP",
		Help:    "The hostname or IP to connect to the platform.",
	},
	Validate: prompt.QuestionValidation{
		Required: true,
	},
}

var username = prompt.Question{
	Name: "Username",
	Type: "Input",
	Options: prompt.QuestionOptions{
		Message: "Basic Auth Username",
		Help:    "The username used to connect to the platform.",
	},
	Validate: prompt.QuestionValidation{
		Required: true,
	},
}

var password = prompt.Question{
	Name: "Password",
	Type: "Password",
	Options: prompt.QuestionOptions{
		Message: "Basic Auth Password",
		Help:    "The password used to connect to the platform",
	},
	Validate: prompt.QuestionValidation{
		Required: true,
	},
}

// contextsAddCmd represents the contextsAdd command
var contextsAddCmd = &cobra.Command{
	Use:     "add <name>",
	Short:   "Adds a new context.",
	Args:    cobra.ExactArgs(1),
	PreRun:  cmdutil.PreRun,
	PostRun: cmdutil.PostRun,
	Run: func(cmd *cobra.Command, args []string) {
		var context = args[0]
		questions := []prompt.Question{}
		answers := make(map[string]interface{})

		answers["Platform"], _ = cmd.Flags().GetString("platform")
		if answers["Platform"] == "" || (answers["Platform"] != "Kubernetes" && answers["Platform"] != "Docker") {
			questions = append(questions, platform)
		}

		answers["Path"], _ = cmd.Flags().GetString("path")
		if answers["Path"] == "" {
			questions = append(questions, path)
		}

		answers["Host"], _ = cmd.Flags().GetString("host")
		if answers["Host"] == "" {
			questions = append(questions, host)
		}

		answers["Username"], _ = cmd.Flags().GetString("username")
		if answers["Username"] == "" {
			questions = append(questions, username)
		}

		answers["Password"], _ = cmd.Flags().GetString("password")
		if answers["Password"] == "" {
			questions = append(questions, password)
		}

		asker := prompt.Prompt{
			Name:      "create-contexts",
			Questions: questions,
		}

		answers = asker.Ask()

		viper.Set("contexts."+context, answers)
		viper.WriteConfig()
	},
}

func init() {
	contextsCmd.AddCommand(contextsAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contextsAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	contextsAddCmd.Flags().String("platform", "", "Deploy to Docker or Kubernetes platform.")
	contextsAddCmd.Flags().String("path", "", "Path.")
	contextsAddCmd.Flags().String("host", "", "Host.")
	contextsAddCmd.Flags().String("username", "", "Username.")
	contextsAddCmd.Flags().String("password", "", "Password.")
}
