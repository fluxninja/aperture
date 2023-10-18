package blueprints

import (
	"context"

	"github.com/spf13/cobra"

	cloudv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cloud/v1"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

// BlueprintsArchiveCmd is the command to archive a blueprint from the Cloud Controller.
var BlueprintsArchiveCmd = &cobra.Command{
	Use:           "archive POLICY_NAME",
	Short:         "Cloud Blueprints Achieve for the given Policy Name",
	Long:          `Archive cloud blueprint.`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	Example:       `aperturectl cloud blueprints archive rate-limiting`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := client.Archive(context.Background(), &cloudv1.DeleteRequest{
			PolicyName: args[0],
		})
		if err != nil {
			return err
		}

		log.Info().Str("policy-name", args[0]).Msg("Successfully archived blueprint for Policy")
		return nil
	},
}
