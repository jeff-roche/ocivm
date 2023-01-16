/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var RELEASE_VERSION string = "0.0.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Release version information",
	Long:  `Display the release version information for this tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Openshift Installer Version Manager")
		fmt.Printf("Version %s %s/%s\n", RELEASE_VERSION, runtime.GOOS, runtime.GOARCH)
		fmt.Printf("Current Installer Version %s\n", PrimaryManifest.CurrentVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
