package dynamicconfig

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var (
	Controller utils.ControllerConn
	client     utils.SelfHostedPolicyClient
)

func init() {
	Controller.InitFlags(DynamicConfigCmd.PersistentFlags())

	DynamicConfigCmd.AddCommand(ApplyCmd)
}

// DynamicConfigCmd is the command to manage DynamicCOnfig of Policies in the Cloud Controller.
var DynamicConfigCmd = &cobra.Command{
	Use:   "dynamic-config",
	Short: "DynamicConfig of Aperture Policy related commands for the Cloud Controller",
	Long: `
Use this command to manage the DynamicConfig of the Aperture Policies to the Cloud Controller.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		err = Controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to run controller pre-run: %w", err)
		}

		client, err = Controller.PolicyClient()
		if err != nil {
			return fmt.Errorf("failed to get controller client: %w", err)
		}

		return nil
	},
	PersistentPostRun: Controller.PostRun,
}
