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
			PolicyName: name,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Name: %s\n", getResponse.GetBlueprint().GetBlueprintsName())
		fmt.Printf("Version: %s\n", getResponse.GetBlueprint().GetVersion())
		fmt.Printf("Policy Name: %s\n", getResponse.GetBlueprint().GetPolicyName())
		fmt.Printf("Values: \n%s\n", getResponse.GetBlueprint().GetValues())

		return nil
	},
}
