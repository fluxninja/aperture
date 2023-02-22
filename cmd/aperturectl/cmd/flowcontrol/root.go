package flowcontrol

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
)

var controller utils.ControllerConn

func init() {
	controller.InitFlags(FlowControlCmd.PersistentFlags())

	FlowControlCmd.AddCommand(ListCmd)
}

// FlowControlCmd is the command to observe control points.
var FlowControlCmd = &cobra.Command{
	Use:               "flow-control {--kube | --controller ADDRESS}",
	Short:             "Flow Control integrations",
	Long:              `Use this command to query information about active Flow Control integrations`,
	SilenceErrors:     true,
	PersistentPreRunE: controller.PreRunE,
	PersistentPostRun: controller.PostRun,
}
