package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var controller utils.ControllerConn

func init() {
	controller.InitFlags(StatusCmd.PersistentFlags())
}

// StatusCmd is the command to get a status of jobs.
var StatusCmd = &cobra.Command{
	Use:           "status",
	Short:         "Get Jobs status",
	Long:          `Use this command to get the status of internal jobs.`,
	SilenceErrors: true,
	Example: `
	aperturectl cloud status
	`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to run controller pre-run: %w", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := controller.StatusClient()
		if err != nil {
			return err
		}

		return utils.ParseStatus(client)
	},
	PersistentPostRun: controller.PostRun,
}
