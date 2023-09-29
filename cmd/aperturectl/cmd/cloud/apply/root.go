package apply

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/utils"
)

// Controller is the controller connection object.
var Controller utils.ControllerConn

func init() {
	Controller.InitFlags(ApplyCmd.PersistentFlags())

	ApplyCmd.AddCommand(ApplyPolicyCmd)
	ApplyCmd.AddCommand(ApplyDynamicConfigCmd)
}

// ApplyCmd is the command to apply a policy to the Cloud Controller.
var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Aperture Policies to the Cloud Controller",
	Long: `
Use this command to apply the Aperture Policies to the Cloud Controller.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := Controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to run controller pre-run: %w", err)
		}

		return nil
	},
	PersistentPostRun: Controller.PostRun,
}
