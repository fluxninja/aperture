package dynamicconfig

import (
	"fmt"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

// GetCmd is a command to get a policy's dynamic config.
var GetCmd = &cobra.Command{
	Use:           "get POLICY_NAME",
	Short:         "Get Aperture DynamicConfig for a Policy.",
	Long:          "Use this command to get the Aperture DynamicConfig of a Policy.",
	Example:       "aperture cloud dynamic-config get rate-limiting",
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		policyName := args[0]
		err := utils.GetDynamicConfigUsingAPI(client, policyName)
		if err != nil {
			return fmt.Errorf("failed to get dynamic config using API: %w", err)
		}

		return nil
	},
}
