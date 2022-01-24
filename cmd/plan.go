/*
Copyright © 2021 Luca Lanziani <luca@lanziani.com>

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
	"fmt"
	githuborg "gitops-toolbox/github-org-manager/lib"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan <organization>",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		organization := args[0]
		token := viper.GetString("github-token")
		org := githuborg.GetClient(token, organization)
		outOfSyncRepos, err := org.GetOutOfSyncRepos()
		if err != nil {
			log.Fatal(err)
		}

		if len(outOfSyncRepos) > 0 {
			fmt.Println(outOfSyncRepos)
		} else {
			fmt.Println("No repos to sync")
		}
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// planCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// planCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
