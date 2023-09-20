package apply

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var (
	// Controller is the controller connection object.
	Controller utils.ControllerConn

	client       utils.PolicyClient
	controllerNs string
)

func init() {
	Controller.InitFlags(ApplyCmd.PersistentFlags())

	ApplyCmd.AddCommand(ApplyPolicyCmd)
	ApplyCmd.AddCommand(ApplyDynamicConfigCmd)
}

// ApplyCmd is the command to apply a policy to the cluster.
var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply Aperture Policies",
	Long: `
Use this command to apply the Aperture Policies.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		err = Controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to run controller pre-run: %w", err)
		}

		controllerNs = utils.GetControllerNs()

		client, err = Controller.PolicyClient()
		if err != nil {
			return fmt.Errorf("failed to get controller client: %w", err)
		}

		return nil
	},
	PersistentPostRun: Controller.PostRun,
}
