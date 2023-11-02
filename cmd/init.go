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

		svn := createSvn()
		err := svn.Init()
		if err != nil {
			fmt.Printf("Error initializing repository %v\n", err)
		}
		fmt.Printf("Repository %s initialized\n", svn.RepoName)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
