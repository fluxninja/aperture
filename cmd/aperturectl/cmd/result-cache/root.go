package statecache

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var controller utils.ControllerConn

func init() {
	controller.InitFlags(FlowControlCmd.PersistentFlags())

	FlowControlCmd.AddCommand(GetCommand)
	FlowControlCmd.AddCommand(SetCommand)
	FlowControlCmd.AddCommand(DeleteCommand)
}

// FlowControlCmd is the command to observe Flow Control control points.
var FlowControlCmd = &cobra.Command{
	Use:               "state-cache",
	Short:             "State Cache related commands",
	Long:              `Use this command to interact with Aperture State Cache.`,
	SilenceErrors:     true,
	PersistentPreRunE: controller.PreRunE,
	PersistentPostRun: controller.PostRun,
}
