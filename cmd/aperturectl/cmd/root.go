package cmd

import (
	"fmt"
	"os"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/blueprints"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "aperturectl",
	Short: "aperturectl - CLI tool to interact with Aperture",
}

func init() {
	RootCmd.AddCommand(blueprints.BlueprintsCmd)
}

// Execute is the entrypoint for the CLI. It is called from the main package.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Err: '%s'", err)
		os.Exit(1)
	}
}
