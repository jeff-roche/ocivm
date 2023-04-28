/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jeff-roche/ocivm/src/installer"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install the specified openshift-installer",
	Long:  "install the specified openshift-installer optionally passing a space delimited list of versions to install",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("no version specified")
			os.Exit(1)
		}

		// Loop over each version specified
		for _, v := range args {
			if err := installer.GetNewInstaller(strings.TrimSpace(v), &PrimaryManifest); err != nil {
				fmt.Printf("unable to install version %s: %s\n", v, err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
