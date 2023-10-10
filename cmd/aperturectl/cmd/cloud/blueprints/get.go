package blueprints

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	cloudv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cloud/v1"
)

func init() {
	BlueprintsGetCmd.Flags().StringVar(&name, "name", "", "Name of blueprint to get")
}

// BlueprintsGetCmd is the command to get a blueprint from the Cloud Controller.
var BlueprintsGetCmd = &cobra.Command{
	Use:           "get",
	Short:         "Cloud Blueprints Get",
	Long:          `Get cloud blueprint.`,
	SilenceErrors: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if name == "" {
			return fmt.Errorf("--name is required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		getResponse, err := client.Get(context.Background(), &cloudv1.GetRequest{
			Name: name,
		}, nil)
		if err != nil {
			return err
		}

		fmt.Println(getResponse.GetBlueprint().GetContent())

		return nil
	},
}
