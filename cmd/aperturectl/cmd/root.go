package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/apply"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/autoscale"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/blueprints"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/build"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/discovery"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/flowcontrol"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/installation"
	"github.com/fluxninja/aperture/pkg/config"
	"github.com/fluxninja/aperture/pkg/info"
	"github.com/fluxninja/aperture/pkg/log"
)

// Version shows the version of ApertureCtl.
var (
	Version         = info.Version
	allowDeprecated bool
	verbose         bool
)

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
	// Process the verbose and allowDeprecated flags before executing the command.
	// Fun fact: PersistentPreRunE doesn't allow chaining.

	// set flags manually using pflag and parse them
	pflag.BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	pflag.BoolVar(&allowDeprecated, "allow-deprecated", false, "Allow deprecated fields in the configuration")
	// configure pflag to ignore unknown flags
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.Parse()
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = false

	level := "info"
	if verbose {
		level = "debug"
	}

	logger := log.NewLogger(log.GetPrettyConsoleWriter(), level)
	log.SetGlobalLogger(logger)

	err := config.RegisterDeprecatedValidator(allowDeprecated)
	if err != nil {
		log.Error().Err(err).Msg("Error registering deprecated validator")
		os.Exit(1)
	}

	if err := RootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("Error executing aperturectl")
		os.Exit(1)
	}
}
