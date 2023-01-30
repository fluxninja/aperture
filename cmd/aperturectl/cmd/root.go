package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aperturectl",
	Short: "aperturectl - CLI tool to interact with Aperture",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute is the entrypoint for the CLI. It is called from the main package.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Err: '%s'", err)
		os.Exit(1)
	}
}
