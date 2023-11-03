package policy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var (
	Controller   utils.ControllerConn
	client       utils.SelfHostedPolicyClient
	controllerNs string
)

func init() {
	Controller.InitFlags(PolicyCmd.PersistentFlags())

	PolicyCmd.AddCommand(ApplyCmd)
	PolicyCmd.AddCommand(GetCmd)
	PolicyCmd.AddCommand(ListCmd)
	PolicyCmd.AddCommand(DeleteCmd)
}

// PolicyCmd is the command to apply a policy to the Controller.
var PolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Aperture Policy related commands for the Controller",
	Long: `
Use this command to manage the Aperture Policies to the Controller.`,
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
