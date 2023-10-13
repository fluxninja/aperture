package blueprints

import (
	"context"

	"github.com/spf13/cobra"

	cloudv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cloud/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// BlueprintsDeleteCmd is the command to delete a blueprint from the Cloud Controller.
var BlueprintsDeleteCmd = &cobra.Command{
	Use:           "delete POLICY_NAME",
	Short:         "Cloud Blueprints Delete for the given Policy Name",
	Long:          `Delete cloud blueprint.`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	Example:       `aperturectl cloud blueprints delete rate-limiting`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := client.Delete(context.Background(), &cloudv1.DeleteRequest{
			PolicyName: args[0],
		})
		if err != nil {
			return err
		}

		log.Info().Str("policy-name", args[0]).Msg("Successfully deleted blueprint for Policy")
		return nil
	},
}
