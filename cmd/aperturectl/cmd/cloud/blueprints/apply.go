package blueprints

import (
	"context"

	"github.com/spf13/cobra"

	cloudv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cloud/v1"
)

// BlueprintsApplyCmd is the command to apply a blueprint from the Cloud Controller.
var BlueprintsApplyCmd = &cobra.Command{
	Use:           "Apply",
	Short:         "Cloud Blueprints Apply",
	Long:          `Apply cloud blueprint.`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := client.Apply(context.Background(), &cloudv1.ApplyRequest{
			Blueprint: &cloudv1.Blueprint{
				Content: []byte("test"),
			},
		}, nil)
		if err != nil {
			return err
		}

		return nil
	},
}
