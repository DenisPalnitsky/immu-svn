/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Adds or updates files in the repository",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := createSvn().Commit()
		if err != nil {
			fmt.Printf("Error committing files %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Files added %d\n", res.FilesAdded)
		fmt.Printf("Files updated %d\n", res.FilesUpdated)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
