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

// configSetCmd represents the configSet command
var configSetCmd = &cobra.Command{
	Use:     "set <key> [value]",
	Short:   "Set a configuration option with its value.",
	Long:    ``,
	Args:    cobra.RangeArgs(1, 2),
	Aliases: []string{"s"},
	PreRun:  cmdutil.PreRun,
	PostRun: cmdutil.PostRun,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var value interface{}

		forced, _ := cmd.Flags().GetBool("forced")

		if len(args) != 2 && !forced {
			cmd.PrintErrln("Running set without a new value will force the value to nil, you must specify -f or --force to continue.")

			return
		}

		if len(args) == 1 {
			viper.Set(key, nil)

			value = "nil"
		} else {
			viper.Set(key, args[1])

			value = args[1]
		}

		viper.WriteConfig()

		cmd.Println("Set" + key + " => " + value.(string))
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configSetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	configSetCmd.Flags().BoolP("forced", "f", false, "Force setting the value to nil.")
}
