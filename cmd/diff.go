/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "diff returns last 10 revisions of a file",
	Run: func(cmd *cobra.Command, args []string) {
		file := cmd.Flag("file").Value.String()
		if file == "" {
			fmt.Println("error: file not specified")
			return
		}

		diff, err := createSvn().Diff(file)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		for i := 0; i < len(diff); i++ {
			fmt.Printf("%s\n", "------------------------------------------")
			contDiff := diff[i].Content
			if i != len(diff)-1 {
				contDiff = cmp.Diff(diff[i].Content, diff[i+1].Content)
			}

			fmt.Printf("-- Revision: %s \t Timestamp: %s\n", diff[i].Revision, diff[i].Timestamp)
			fmt.Printf("Diff:\n %s\n", contDiff)
			fmt.Printf("%s\n", "------------------------------------------")
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringP("file", "f", "", "file to diff")
}
