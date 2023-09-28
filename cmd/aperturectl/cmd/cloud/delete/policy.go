package delete

import (
	"fmt"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// DeletePolicyCmd is the command to delete a policy from the Aperture Cloud Controller.
var DeletePolicyCmd = &cobra.Command{
	Use:           "policy",
	Short:         "Delete Aperture Policy from the Aperture Cloud Controller",
	Long:          `Use this command to delete the Aperture Policy from the Aperture Cloud Controller.`,
	SilenceErrors: true,
	Example:       `aperturectl cloud delete policy --policy=rate-limiting --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key PERSONAL_API_KEY`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return deletePolicy()
	},
}

// deletePolicy deletes the policy from the cluster.
func deletePolicy() error {
	err := utils.DeletePolicyUsingAPI(client, policyName)
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	log.Info().Str("policy", policyName).Msg("Deleted Policy successfully")
	return nil
}
