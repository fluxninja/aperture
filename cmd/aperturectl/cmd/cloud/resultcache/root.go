package resultcache

import (
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/utils"
)

var controller utils.ControllerConn

func init() {
	controller.InitFlags(CacheCmd.PersistentFlags())

	CacheCmd.AddCommand(GetCommand)
	CacheCmd.AddCommand(SetCommand)
	CacheCmd.AddCommand(DeleteCommand)
}

// CacheCmd is the command to observe Flow Control control points.
var CacheCmd = &cobra.Command{
	Use:               "result-cache",
	Short:             "Result Cache related commands",
	Long:              `Use this command to interact with Aperture's Result Cache.`,
	SilenceErrors:     true,
	PersistentPreRunE: controller.PreRunE,
	PersistentPostRun: controller.PostRun,
}
