package autoscale

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var controller utils.ControllerConn

func init() {
	controller.InitFlags(AutoScaleCmd.PersistentFlags())

	AutoScaleCmd.AddCommand(ControlPointsCmd)
}

// AutoScaleCmd is the command to observe AutoScale control points.
var AutoScaleCmd = &cobra.Command{
	Use:               "auto-scale",
	Short:             "AutoScale integrations",
	Long:              `Use this command to query information about active AutoScale integrations`,
	SilenceErrors:     true,
	PersistentPreRunE: controller.PreRunE,
	PersistentPostRun: controller.PostRun,
}
