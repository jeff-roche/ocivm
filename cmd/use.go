package cmd

import (
	"fmt"
	"ocivm/src/manager"
	"os"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Set the active version of the installer",
	Long:  "Set the active version of the installer",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("1 and only 1 version may be specified. Got %s\n", args)
			os.Exit(1)
		}

		if err := manager.UseVersion(args[0], &PrimaryManifest); err != nil {
			fmt.Printf("unable to set the active version to %s: %s\n", args[0], err)
			os.Exit(1)
		}

		fmt.Printf("Successfully activated openshift-install version %s\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
