package policy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	policylangv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

// GetCmd is the command to get a policy from the Aperture Cloud Controller.
var GetCmd = &cobra.Command{
	Use:           "get POLICY_NAME",
	Short:         "Get Aperture Policy from the Aperture Cloud Controller",
	Long:          `Use this command to get the Aperture Policy from the Aperture Cloud Controller.`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	Example:       `aperturectl cloud policy get POLICY_NAME`,
	RunE: func(_ *cobra.Command, args []string) error {
		policy, err := client.GetPolicy(context.Background(), &policylangv1.GetPolicyRequest{
			Name: args[0],
		})
		if err != nil {
			return fmt.Errorf("failed to get policy: %w", err)
		}

		jsonBytes, err := policy.Policy.MarshalJSON()
		if err != nil {
			return err
		}

		yamlString, err := utils.GetYAMLString(jsonBytes)
		if err != nil {
			return err
		}

		fmt.Println(yamlString)
		return nil
	},
}
