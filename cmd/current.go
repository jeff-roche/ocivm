/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "List the currently active openshift cli version",
	Long:  "List the currently active openshift cli version",
	Run: func(cmd *cobra.Command, args []string) {
		if PrimaryManifest.CurrentVersion == "" {
			fmt.Println("No version of openshift is currently active")
		} else {
			fmt.Println(PrimaryManifest.CurrentVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
