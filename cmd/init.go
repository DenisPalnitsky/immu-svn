/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize remote repository",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("Repository name %v\n", getRepoName())
		createSvn().Init()

		fmt.Printf("Repository %s initialized\n", getRepoName())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
