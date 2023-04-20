package blueprints

import (
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull Aperture Blueprints",
	Long: `
Use this command to pull the Aperture Blueprints in local system to use for generating Aperture Policies and Grafana Dashboards.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Example: `aperturectl blueprints pull

aperturectl blueprints pull --version latest`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
