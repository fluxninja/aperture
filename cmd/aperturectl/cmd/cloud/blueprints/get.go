package blueprints

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	cloudv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/cloud/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
)

// BlueprintsGetCmd is the command to get a blueprint from the Cloud Controller.
var BlueprintsGetCmd = &cobra.Command{
	Use:           "get POLICY_NAME",
	Short:         "Cloud Blueprints Get",
	Long:          `Get cloud blueprint.`,
	SilenceErrors: true,
	Args:          cobra.ExactArgs(1),
	Example:       `aperturectl cloud blueprints get rate-limiting`,
	RunE: func(cmd *cobra.Command, args []string) error {
		getResponse, err := client.Get(context.Background(), &cloudv1.GetRequest{
			PolicyName: args[0],
		})
		if err != nil {
			return err
		}

		fmt.Printf("Name: %s\n", getResponse.GetBlueprint().GetBlueprintsName())
		fmt.Printf("Version: %s\n", getResponse.GetBlueprint().GetVersion())
		fmt.Printf("Policy Name: %s\n", getResponse.GetBlueprint().GetPolicyName())

		yamlString, err := utils.GetYAMLString(getResponse.GetBlueprint().GetValues())
		if err != nil {
			return err
		}
		yamlString = strings.ReplaceAll(yamlString, "\n", "\n\t")
		fmt.Printf("Values: \n\t%s\n", yamlString)

		return nil
	},
}
