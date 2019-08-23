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
	"github.com/spf13/viper"

	"github.com/lavrahq/cli/util/cmdutil"
	"github.com/spf13/cobra"
)

// contextsSwitchCmd represents the contextsSwitch command
var contextsSwitchCmd = &cobra.Command{
	Use:     "switch <name>",
	Short:   "Switches to the names context.",
	PreRun:  cmdutil.PreRun,
	PostRun: cmdutil.PostRun,
	Aliases: []string{"use"},
	Run: func(cmd *cobra.Command, args []string) {
		var key = args[0]

		viper.Set("currentContext", key)
		viper.WriteConfig()

		cmd.Println("Set the current context to '" + key + "'!")
	},
}

func init() {
	contextsCmd.AddCommand(contextsSwitchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contextsSwitchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// contextsSwitchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
