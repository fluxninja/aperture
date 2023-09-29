package autoscale

import (
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/spf13/cobra"
)

// ControlPointsCmd is the command to list control points.
var ControlPointsCmd = &cobra.Command{
	Use:           "control-points",
	Short:         "List AutoScale control points",
	Long:          `List AutoScale control points`,
	SilenceErrors: true,
	Example:       `aperturectl cloud auto-scale control-points`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		return utils.ParseControlPoints(client)
	},
}
