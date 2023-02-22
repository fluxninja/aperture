package controlpoints

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
)

var controller utils.ControllerConn

func init() {
	controller.InitFlags(ControlPointsCmd.PersistentFlags())

	ControlPointsCmd.AddCommand(ListCmd)
}

// ControlPointsCmd is the command to observe control points.
var ControlPointsCmd = &cobra.Command{
	Use:               "experimental-control-points {--kube | --controller ADDRESS}",
	Short:             "Query control points information",
	Long:              `Use this command to get information about discovered control points`,
	SilenceErrors:     true,
	PersistentPreRunE: controller.PreRunE,
	PersistentPostRun: controller.PostRun,
}
