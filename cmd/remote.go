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
	"encoding/json"
	"fmt"
	"log"

	githuborg "github.com/gitops-toolbox/github-org-manager/lib"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote <organization>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		organization := args[0]
		token := viper.GetString("github-token")
		org := githuborg.GetClient(token, organization)
		repos, err := org.GetRepos()

		if err != nil {
			fmt.Println(err)
		}

		for _, repo := range repos {
			jr, err := json.Marshal(repo)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%s\n", jr)
		}
	},
}

func init() {
	rootCmd.AddCommand(remoteCmd)

	// Here you will define your flags and configuration settings.

	//Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// remoteCmd.PersistentFlags().String("github-token", "", "Github token to authenticate")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	remoteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
