package dynamicconfig

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	languagev1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/policy/language/v1"
)

// DelCmd is a command to delete a policy's dynamic config.
var DelCmd = &cobra.Command{
	Use:           "delete POLICY_NAME",
	Short:         "Delete Aperture DynamicConfig of a Policy.",
	Long:          "Use this command to delete the Aperture DynamicConfig of a Policy.",
	Example:       "aperture cloud dynamic-config delete rate-limiting",
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		_, err := client.DeleteDynamicConfig(context.Background(), &languagev1.DeleteDynamicConfigRequest{
			PolicyName: args[0],
		})
		if err != nil {
			return err
		}
		fmt.Println("Deleted DynamicConfig successfully")
		return nil
	},
}
