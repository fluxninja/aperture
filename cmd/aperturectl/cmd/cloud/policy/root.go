package policy

import (
	"fmt"

	"github.com/spf13/cobra"

	cloudutils "github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/cloud/utils"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var (
	Controller  cloudutils.ControllerConn
	client      utils.PolicyClient
	cloudClient utils.CloudPolicyClient
)

func init() {
	Controller.InitFlags(PolicyCmd.PersistentFlags())

	PolicyCmd.AddCommand(ApplyCmd)
	PolicyCmd.AddCommand(GetCmd)
	PolicyCmd.AddCommand(ListCmd)
	PolicyCmd.AddCommand(DeleteCmd)
}

// PolicyCmd is the command to apply a policy to the Cloud Controller.
var PolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Aperture Policy related commands for the Cloud Controller",
	Long: `
Use this command to manage the Aperture Policies to the Cloud Controller.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := Controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to run controller pre-run: %w", err)
		}

		client, err = Controller.PolicyClient()
		if err != nil {
			return fmt.Errorf("failed to get cloud controller client: %w", err)
		}

		cloudClient, err = Controller.CloudPolicyClient()
		if err != nil {
			return fmt.Errorf("failed to get cloud controller client: %w", err)
		}

		return nil
	},
	PersistentPostRun: Controller.PostRun,
}
