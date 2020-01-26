/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"log"

	"github.com/jnathanh/recipe/lib"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list recipes",
	Long: `recursively lists all recipes in the recipe home folder`,
	Run: func(cmd *cobra.Command, args []string) {
		usePath, err := cmd.Flags().GetBool("path")
		if err != nil {
			log.Fatal(err)
		}

		recipes, err := lib.List()
		if err != nil {
			log.Fatal(err)
		}

		for i := range recipes {
			if usePath {
				fmt.Println(recipes[i].Path)
				continue
			}

			// default to recipe name
			fmt.Println(recipes[i].Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	lsCmd.Flags().BoolP("path", "p", false, "print the path of the recipe file, rather than the recipe name")
}
