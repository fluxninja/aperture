package policy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	policylangv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/policy/language/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

// GetCmd is the command to get a policy from the Aperture Controller.
var GetCmd = &cobra.Command{
	Use:           "get POLICY_NAME",
	Short:         "Get Aperture Policy from the Aperture Controller",
	Long:          `Use this command to get the Aperture Policy from the Aperture Controller.`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	Example:       `aperturectl policy get POLICY_NAME`,
	RunE: func(_ *cobra.Command, args []string) error {
		policy, err := client.GetPolicy(context.Background(), &policylangv1.GetPolicyRequest{
			Name: args[0],
		})
		if err != nil {
			return err
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
