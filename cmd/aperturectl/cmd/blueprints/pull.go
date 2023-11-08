package blueprints

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var pullCmd = &cobra.Command{
	Use:           "pull",
	Short:         "Pull Aperture Blueprints",
	Long:          `Use this command to pull the Aperture Blueprints in local system to use for generating Aperture Policies and Grafana Dashboards.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Example: `aperturectl blueprints pull

aperturectl blueprints pull --version latest`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _, _, err := utils.Pull(blueprintsURI, blueprintsVersion, blueprints, utils.DefaultBlueprintsRepo, skipPull, true)
		if err != nil {
			return err
		}
		return nil
	},
}
