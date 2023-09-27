package autoscale

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

// ControlPointsCmd is the command to list control points.
var ControlPointsCmd = &cobra.Command{
	Use:           "control-points",
	Short:         "List AutoScale control points",
	Long:          `List AutoScale control points`,
	SilenceErrors: true,
	Example:       `aperturectl auto-scale control-points`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		return utils.ParseAutoScaleControlPoints(client)
	},
}
