package blueprints

import (
	"github.com/spf13/cobra"
)

const (
	// DefaultBlueprintsRepo is the default repository for blueprints.
	DefaultBlueprintsRepo = "github.com/fluxninja/aperture/blueprints"
	// LatestTag is the tag for the latest version of blueprints.
	LatestTag = "latest"
)

var (
	// Args for `blueprints`.
	blueprintsURI     string
	blueprintsVersion string
)

func init() {
	BlueprintsCmd.PersistentFlags().StringVar(&blueprintsVersion, "version", "", "Version of official Aperture Blueprints, e.g. latest. This field should not be provided when the URI is provided")
	BlueprintsCmd.PersistentFlags().StringVar(&blueprintsURI, "uri", "", "URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@latest. This field should not be provided when the Version is provided.")
	BlueprintsCmd.PersistentFlags().BoolVar(&skipPull, "skip-pull", false, "Skip pulling the blueprints update.")

	BlueprintsCmd.AddCommand(pullCmd)
	BlueprintsCmd.AddCommand(listCmd)
	BlueprintsCmd.AddCommand(removeCmd)
	BlueprintsCmd.AddCommand(generateCmd)
	BlueprintsCmd.AddCommand(valuesCmd)
	BlueprintsCmd.AddCommand(dynamicValuesCmd)
}

// BlueprintsCmd is the root command for blueprints.
var BlueprintsCmd = &cobra.Command{
	Use:   "blueprints",
	Short: "Aperture Blueprints",
	Long:  `Use this command to pull, list, remove and generate Aperture Policy resources using the Aperture Blueprints.`,
}
