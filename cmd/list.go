/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list available versions of the openshift cli tools",
	Long:  "list available versions of the openshift cli tools",
	Run: func(cmd *cobra.Command, args []string) {
		current, _ := cmd.Flags().GetBool("current")
		remote, _ := cmd.Flags().GetBool("remote")

		PrimaryManifest.ListVersions(current, remote)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().BoolP("current", "c", false, "Highlights the currently active version")
	listCmd.PersistentFlags().BoolP("remote", "r", false, "Lists all available versions from the remote and marks installed versions with a *")
}
