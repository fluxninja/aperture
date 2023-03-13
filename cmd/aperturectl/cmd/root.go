package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/apply"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/autoscale"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/blueprints"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/build"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/discovery"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/flowcontrol"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/installation"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
)

// Version shows the version of ApertureCtl.
var Version = info.Version

func init() {
	RootCmd.AddCommand(blueprints.BlueprintsCmd)
	RootCmd.AddCommand(compileCmd)
	RootCmd.AddCommand(apply.ApplyCmd)
	RootCmd.AddCommand(installation.InstallCmd)
	RootCmd.AddCommand(installation.UnInstallCmd)
	RootCmd.AddCommand(flowcontrol.FlowControlCmd)
	RootCmd.AddCommand(autoscale.AutoScaleCmd)
	RootCmd.AddCommand(discovery.DiscoveryCmd)
	RootCmd.AddCommand(build.BuildCmd)
	RootCmd.AddCommand(agentsCmd)

	RootCmd.InitDefaultCompletionCmd()
	RootCmd.SilenceUsage = true
}

// RootCmd is the root command for aperturectl.
var RootCmd = &cobra.Command{
	Use:                "aperturectl",
	Short:              "aperturectl - CLI tool to interact with Aperture",
	DisableAutoGenTag:  true,
	DisableSuggestions: false,
	Long: `
aperturectl is a CLI tool which can be used to interact with Aperture seamlessly.`,
	Version: Version,
}

// Execute is the entrypoint for the CLI. It is called from the main package.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("Error executing aperturectl")
		os.Exit(1)
	}
}
