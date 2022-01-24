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
	githuborg "gitops-toolbox/github-org-manager/lib"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <organization> <reponame>",
	Short: "print json representation of the repo to stdout",
	Long:  ``,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		organization, reponame := args[0], args[1]
		token := viper.GetString("github-token")
		org := githuborg.GetClient(token, organization)

		log.Println("Importing", reponame)

		repo, err := org.GetRepo(reponame)

		if err != nil {
			log.Fatal(err)
		}

		repoAsString, err := json.MarshalIndent(map[string]githuborg.Repo{
			*repo.Name: *repo,
		}, "", "  ")

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", repoAsString)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	importCmd.PersistentFlags().String("name", "", "reponame")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
