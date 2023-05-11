package discovery

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var controller utils.ControllerConn

func init() {
	controller.InitFlags(DiscoveryCmd.PersistentFlags())

	DiscoveryCmd.AddCommand(EntitiesCmd)
}

// DiscoveryCmd is the command to observe AutoScale control points.
var DiscoveryCmd = &cobra.Command{
	Use:               "discovery",
	Short:             "Discovery integrations",
	Long:              `Use this command to query information about active Discovery integrations`,
	SilenceErrors:     true,
	PersistentPreRunE: controller.PreRunE,
	PersistentPostRun: controller.PostRun,
}
