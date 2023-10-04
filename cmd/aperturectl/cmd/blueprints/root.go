package blueprints

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/info"
)

const (
	defaultBlueprintsRepo = "github.com/fluxninja/aperture/blueprints"
	latestTag             = "latest"
)

var (
	// Location of cache for blueprints. E.g. ~/.aperturectl/blueprints.
	blueprintsCacheRoot string
	// Location of blueprints uri within cache. E.g. ~/.aperturectl/blueprints/github.com/fluxninja/aperture/blueprints@latest.
	blueprintsURIRoot string

	// Location of blueprints directory within URI directory. E.g. ~/.aperturectl/blueprints/github.com_fluxninja_aperture_blueprints@v0.26.1/github.com/fluxninja/aperture/blueprints/.
	blueprintsDir string

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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		newer, err := utils.IsCurrentVersionNewer(info.Version)
		if err != nil {
			return err
		}
		if newer {
			err = removeCmd.RunE(cmd, args)
			if err != nil {
				return err
			}
			err = removeCmd.PostRunE(cmd, args)
			if err != nil {
				return err
			}
			err = pullCmd.RunE(cmd, args)
			if err != nil {
				return err
			}
			err = pullCmd.PostRunE(cmd, args)
			if err != nil {
				return err
			}
			return utils.UpdateVersionFile(info.Version)
		}
		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		return pullCmd.PostRunE(cmd, args)
	},
}
