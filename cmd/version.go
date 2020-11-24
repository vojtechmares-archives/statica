package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string = ""

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version of statica",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run: func(c *cobra.Command, args []string) {
		fmt.Printf("Statica version: %s\n", version)
	},
}
