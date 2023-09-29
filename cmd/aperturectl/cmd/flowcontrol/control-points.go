package flowcontrol

import (
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/spf13/cobra"
)

// ControlPointsCmd is the command to list control points.
var ControlPointsCmd = &cobra.Command{
	Use:           "control-points",
	Short:         "List Flow Control control points",
	Long:          `List Flow Control control points`,
	SilenceErrors: true,
	Example:       `aperturectl flow-control control-points`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := controller.IntrospectionClient()
		if err != nil {
			return err
		}

		return utils.ParseControlPoints(client)
	},
}
