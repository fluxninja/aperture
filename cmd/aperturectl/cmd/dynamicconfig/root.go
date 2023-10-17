package dynamicconfig

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var (
	controller   utils.ControllerConn
	client       utils.PolicyClient
	controllerNs string
)

func init() {
	controller.InitFlags(DynamicConfigCmd.PersistentFlags())

	DynamicConfigCmd.AddCommand(ApplyCmd)
	DynamicConfigCmd.AddCommand(GetCmd)
	DynamicConfigCmd.AddCommand(DelCmd)
}

// DynamicConfigCmd is the command to manage DynamicCOnfig of Policies in the Controller.
var DynamicConfigCmd = &cobra.Command{
	Use:   "dynamic-config",
	Short: "DynamicConfig of Aperture Policy related commands for the Controller",
	Long: `
Use this command to manage the DynamicConfig of the Aperture Policies to the Controller.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		err = controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to run controller pre-run: %w", err)
		}

		controllerNs = utils.GetControllerNs()

		client, err = controller.PolicyClient()
		if err != nil {
			return fmt.Errorf("failed to get controller client: %w", err)
		}

		return nil
	},
	PersistentPostRun: controller.PostRun,
}
