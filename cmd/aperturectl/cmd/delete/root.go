package delete

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	cmdv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cmd/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

var (
	controller   utils.ControllerConn
	client       cmdv1.ControllerClient
	controllerNs string
	policyName   string
)

func init() {
	controller.InitFlags(DeleteCmd.PersistentFlags())
	DeleteCmd.PersistentFlags().StringVar(&policyName, "policy", "", "Name of the Policy to delete")

	DeleteCmd.AddCommand(DeletePolicyCmd)
}

// DeleteCmd is the command to delete a policy from the cluster.
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Aperture Policies",
	Long: `
Use this command to delete the Aperture Policies.`,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if policyName == "" {
			return errors.New("policy name is required")
		}

		var err error
		err = controller.PreRunE(cmd, args)
		if err != nil {
			return fmt.Errorf("failed to run controller pre-run: %w", err)
		}

		controllerNs = utils.GetControllerNs()

		client, err = controller.Client()
		if err != nil {
			return fmt.Errorf("failed to get controller client: %w", err)
		}
		return nil
	},
	PersistentPostRun: controller.PostRun,
}
