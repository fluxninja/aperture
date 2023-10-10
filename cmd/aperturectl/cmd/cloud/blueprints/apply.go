package blueprints

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cloudv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cloud/v1"
)

func init() {
	BlueprintsApplyCmd.Flags().StringVar(&valuesFile, "values-file", "", "Values file to use for blueprint")
}

// BlueprintsApplyCmd is the command to apply a blueprint from the Cloud Controller.
var BlueprintsApplyCmd = &cobra.Command{
	Use:           "Apply",
	Short:         "Cloud Blueprints Apply",
	Long:          `Apply cloud blueprint.`,
	SilenceErrors: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if valuesFile == "" {
			return fmt.Errorf("--values-file is required")
		}

		_, err := os.Stat(valuesFile)
		if err != nil && os.IsNotExist(err) {
			return fmt.Errorf("values file does not exist: %w", err)
		}

		content, err := os.ReadFile(valuesFile)
		if err != nil {
			return fmt.Errorf("failed to read values file: %w", err)
		}
		valuesFileContent = content

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := client.Apply(context.Background(), &cloudv1.ApplyRequest{
			Blueprint: &cloudv1.Blueprint{
				Content: valuesFileContent,
			},
		}, nil)
		if err != nil {
			return err
		}

		return nil
	},
}
