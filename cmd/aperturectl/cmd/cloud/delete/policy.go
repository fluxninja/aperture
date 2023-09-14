package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	languagev1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// DeletePolicyCmd is the command to delete a policy from the Aperture Cloud Controller.
var DeletePolicyCmd = &cobra.Command{
	Use:           "policy",
	Short:         "Delete Aperture Policy from the Aperture Cloud Controller",
	Long:          `Use this command to delete the Aperture Policy from the Aperture Cloud Controller.`,
	SilenceErrors: true,
	Example:       `aperturectl cloud delete policy --policy=rate-limiting --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key API_KEY`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return deletePolicy()
	},
}

// deletePolicy deletes the policy from the cluster.
func deletePolicy() error {
	err := deletePolicyUsingAPI()
	if err != nil {
		return fmt.Errorf("failed to delete policy: %w", err)
	}

	log.Info().Str("policy", policyName).Msg("Deleted Policy successfully")
	return nil
}

// deletePolicyUsingAPI deletes the policy using the API.
func deletePolicyUsingAPI() error {
	policyRequest := languagev1.DeletePolicyRequest{
		Name: policyName,
	}
	_, err := client.DeletePolicy(context.Background(), &policyRequest)
	if err != nil {
		log.Warn().Err(err).Str("policy", policyName).Msg("failed to delete Policy")
	}

	return nil
}
