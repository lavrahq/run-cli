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

// contextsAddCmd represents the contextsAdd command
var contextsAddCmd = &cobra.Command{
	Use:     "add <name>",
	Short:   "Adds a new context.",
	Args:    cobra.ExactArgs(1),
	PreRun:  cmdutil.PreRun,
	PostRun: cmdutil.PostRun,
	Run: func(cmd *cobra.Command, args []string) {
		var context = args[0]
		answers := make(map[string]interface{})

		if cmd.Flags().NFlag() >= 4 {
			// If set through flags
			platform, _ := cmd.Flags().GetString("platform")
			host, _ := cmd.Flags().GetString("host")
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")

			answers["platform"] = platform
			answers["host"] = host
			answers["username"] = username
			answers["password"] = password

			if platform == "Docker" {
				path, _ := cmd.Flags().GetString("path")
				answers["path"] = path
			}
		} else {
			// If set through prompt
			deployment := prompt.Prompt{
				Name: "Add Context",
				Questions: []prompt.Question{
					prompt.Question{
						Name: "platform",
						Type: "Select",
						Options: prompt.QuestionOptions{
							Options: []string{
								"Kubernetes",
								"Docker",
							},
						},
						Validate: prompt.QuestionValidation{
							Required: true,
						},
					},
				},
			}

			deploymentAnswers := deployment.Ask()

			if deploymentAnswers["Deployment Platform"] == "Docker" {
				options := prompt.Prompt{
					Name: "Docker Deployment",
					Questions: []prompt.Question{
						prompt.Question{
							Name: "path",
							Type: "Input",
							Options: prompt.QuestionOptions{
								Message: "What is the path to your API (Docker)?",
								Default: "/var/run/docker.sock",
							},
							Validate: prompt.QuestionValidation{
								Required: true,
							},
						},
						prompt.Question{
							Name: "host",
							Type: "Input",
							Options: prompt.QuestionOptions{
								Message: "What is the IP of the host?",
							},
							Validate: prompt.QuestionValidation{
								Required: true,
							},
						},
						prompt.Question{
							Name: "username",
							Type: "Input",
							Options: prompt.QuestionOptions{
								Message: "What is the username you'd like to use?",
							},
							Validate: prompt.QuestionValidation{
								Required: true,
							},
						},
						prompt.Question{
							Name: "password",
							Type: "Password",
							Options: prompt.QuestionOptions{
								Message: "What is the password you'd like to use?",
							},
							Validate: prompt.QuestionValidation{
								Required: true,
							},
						},
					},
				}

				answers = options.Ask()
			}

			if deploymentAnswers["Deployment Platform"] == "Kubernetes" {
				options := prompt.Prompt{
					Name: "Kubernetes Deployment",
					Questions: []prompt.Question{
						prompt.Question{
							Name: "host",
							Type: "Input",
							Options: prompt.QuestionOptions{
								Message: "What is the IP of the host?",
							},
							Validate: prompt.QuestionValidation{
								Required: true,
							},
						},
						prompt.Question{
							Name: "username",
							Type: "Input",
							Options: prompt.QuestionOptions{
								Message: "What is the username you'd like to use?",
							},
							Validate: prompt.QuestionValidation{
								Required: true,
							},
						},
						prompt.Question{
							Name: "password",
							Type: "Password",
							Options: prompt.QuestionOptions{
								Message: "What is the password you'd like to use?",
							},
							Validate: prompt.QuestionValidation{
								Required: true,
							},
						},
					},
				}

				answers = options.Ask()
			}
		}

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
