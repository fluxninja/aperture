package blueprints

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	cloudv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cloud/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/blueprints"
)

func init() {
	BlueprintsApplyCmd.Flags().StringVar(&valuesFile, "values-file", "", "Values file to use for blueprint")
}

// BlueprintsApplyCmd is the command to apply a blueprint from the Cloud Controller.
var BlueprintsApplyCmd = &cobra.Command{
	Use:           "apply",
	Short:         "Cloud Blueprints Apply",
	Long:          `Apply cloud blueprint.`,
	SilenceErrors: true,
	Example:       `aperturectl cloud blueprints apply --value-file=values.yaml`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if valuesFile == "" {
			return fmt.Errorf("--values-file is required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		valuesMap, err := blueprints.Generate(valuesFile, "", "", "", false)
		if err != nil {
			return err
		}

		policy, ok := valuesMap["policy"].(map[string]any)
		if !ok {
			return fmt.Errorf("policy not found in blueprint values")
		}

		policyName, ok := policy["policy_name"].(string)
		if !ok {
			return fmt.Errorf("policy_name not found in blueprint values")
		}

		blueprintName, ok := valuesMap["blueprint"].(string)
		if !ok {
			return fmt.Errorf("blueprint not found in blueprint values")
		}

		uri, ok := valuesMap["uri"].(string)
		if !ok {
			return fmt.Errorf("uri not found in blueprint values")
		}

		var version string
		uriSlice := strings.Split(uri, "@")
		if len(uriSlice) != 2 {
			version = "latest"
		} else {
			version = uriSlice[1]
		}

		valuesFileContent, err := json.Marshal(valuesMap)
		if err != nil {
			return err
		}

		_, err = client.Apply(context.Background(), &cloudv1.ApplyRequest{
			Blueprint: &cloudv1.Blueprint{
				PolicyName:     policyName,
				Version:        version,
				Values:         valuesFileContent,
				BlueprintsName: blueprintName,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to apply blueprint: %w", err)
		}

		return nil
	},
}
