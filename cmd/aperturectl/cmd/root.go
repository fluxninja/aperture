package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/apply"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/autoscale"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/blueprints"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/build"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/decisions"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/delete"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/discovery"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/flowcontrol"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/installation"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/status"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/config"
	"github.com/fluxninja/aperture/v2/pkg/info"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// Version shows the version of ApertureCtl.
var (
	Version = info.Version
	verbose bool

	controller utils.ControllerConn
)

func init() {
	RootCmd.AddCommand(cloud.CloudCmd)
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
	RootCmd.AddCommand(delete.DeleteCmd)
	RootCmd.AddCommand(decisions.DecisionsCmd)
	RootCmd.AddCommand(policiesCmd)
	RootCmd.AddCommand(status.StatusCmd)

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

// Execute is the entry point for the CLI. It is called from the main package.
func Execute() {
	// Process the verbose and allowDeprecated flags before executing the command.
	// Fun fact: PersistentPreRunE does not allow chaining.

	// set flags manually using pflag and parse them
	pflag.BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	pflag.BoolVar(&utils.AllowDeprecated, "allow-deprecated", false, "Allow deprecated fields in the configuration")
	// Set help flag to false by default to print help for aperturectl command instead of pflag
	pflag.BoolP("help", "h", false, "Display help for aperturectl command")
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

	err := config.RegisterDeprecatedValidator(utils.AllowDeprecated)
	if err != nil {
		log.Error().Err(err).Msg("Error registering deprecated validator")
		os.Exit(1)
	}

	if err := RootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("Error executing aperturectl")
		os.Exit(1)
	}
}
