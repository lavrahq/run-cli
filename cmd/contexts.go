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
	"github.com/lavrahq/cli/util/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// contextsCmd represents the contexts command
var contextsCmd = &cobra.Command{
	Use:   "contexts",
	Short: "Manage configuration options for contexts.",
	Long: `Contexts provide a means for configuring where the Lavra CLI tool connects in order to deploy, configure, and administer
	existing Lavra products, or new products that have not been deployed. By default, when start is ran, a local context is
created that points to the local Docker engine or Kubernetes cluster.`,
	Aliases: []string{"contexts ls"},
	PreRun:  cmdutil.PreRun,
	PostRun: cmdutil.PostRun,
	Run: func(cmd *cobra.Command, args []string) {
		if help, _ := cmd.Flags().GetBool("help"); help {
			cmd.Println(cmd.HelpTemplate())

			return
		}

		currentContext := viper.Get("currentContext")

		for key := range viper.GetStringMap("contexts") {
			if key == currentContext {
				cmd.Print("✔️  ")
			} else {
				cmd.Print("   ")
			}

			cmd.Println(string(key))
		}
	},
}

func init() {
	rootCmd.AddCommand(contextsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contextsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	contextsCmd.Flags().BoolP("help", "h", false, "Lists the available contexts, with an indicator on which is the current context.")
}
