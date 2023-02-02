package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/blueprints"
)

func init() {
	RootCmd.AddCommand(blueprints.BlueprintsCmd)
	RootCmd.AddCommand(compileCmd)
	RootCmd.AddCommand(versionCmd)
}

// RootCmd is the root command for aperturectl.
var RootCmd = &cobra.Command{
	Use:   "aperturectl",
	Short: "aperturectl - CLI tool to interact with Aperture",
}

// Execute is the entrypoint for the CLI. It is called from the main package.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		os.Exit(1)
	}
}
