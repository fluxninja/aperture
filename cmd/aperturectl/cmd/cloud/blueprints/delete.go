package blueprints

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	cloudv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cloud/v1"
)

func init() {
	BlueprintsDeleteCmd.Flags().StringVar(&name, "name", "", "Name of blueprint to get")
}

// BlueprintsDeleteCmd is the command to delete a blueprint from the Cloud Controller.
var BlueprintsDeleteCmd = &cobra.Command{
	Use:           "delete",
	Short:         "Cloud Blueprints Delete",
	Long:          `Delete cloud blueprint.`,
	SilenceErrors: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if name == "" {
			return fmt.Errorf("--name is required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := client.Delete(context.Background(), &cloudv1.DeleteRequest{
			PolicyName: name,
		})
		if err != nil {
			return err
		}

		return nil
	},
}
