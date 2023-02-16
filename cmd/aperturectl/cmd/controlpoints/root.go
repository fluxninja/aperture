package controlpoints

import (
	"github.com/spf13/cobra"
)

var (
	controllerAddr string
	insecure       bool
)

func init() {
	ControlPointsCmd.PersistentFlags().StringVar(
		&controllerAddr,
		"controller",
		"aperture-controller.aperture-controller",
		"Controller's address",
	)
	ControlPointsCmd.PersistentFlags().BoolVar(
		&insecure,
		"insecure",
		false,
		"Don't check TLS certificates when connecting to controller",
	)

	ControlPointsCmd.AddCommand(ListCmd)
}

// ControlPointsCmd is the command to observe control points
var ControlPointsCmd = &cobra.Command{
	Use:           "experimental-control-points",
	Short:         "Query control points information",
	Long:          `Use this command to get information about discovered control points`,
	SilenceErrors: true,
}
