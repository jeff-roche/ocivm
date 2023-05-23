/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jeff-roche/ocivm/src/installer"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall a specific version of the openshift tools",
	Long:  "Uninstall a specific version of the openshift tools",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("no version specified")
			os.Exit(1)
		}

		for _, ver := range args {
			if err := installer.UninstallIfInstalled(ver, &PrimaryManifest); err != nil {
				fmt.Printf("unable to uninstall version \"%s\": %s\n", ver, err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
