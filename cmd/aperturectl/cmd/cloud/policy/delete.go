package policy

import (
	"fmt"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// DeleteCmd is the command to delete a policy from the Aperture Cloud Controller.
var DeleteCmd = &cobra.Command{
	Use:           "delete POLICY_NAME",
	Short:         "Delete Aperture Policy from the Aperture Cloud Controller",
	Long:          `Use this command to delete the Aperture Policy from the Aperture Cloud Controller.`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	Example:       `aperturectl cloud policy delete POLICY_NAME`,
	RunE: func(_ *cobra.Command, args []string) error {
		return deletePolicy(args[0])
	},
}

// deletePolicy deletes the policy from the cluster.
func deletePolicy(policyName string) error {
	err := utils.DeletePolicyUsingAPI(cloudClient, policyName)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	log.Info().Str("policy", policyName).Msg("Deleted Policy successfully")
	return nil
}
