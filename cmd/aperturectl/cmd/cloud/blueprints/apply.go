package blueprints

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	cloudv1 "github.com/fluxninja/aperture/v2/api/gen/proto/go/aperture/cloud/v1"
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
		valuesMap := make(map[string]any)
		err := yaml.Unmarshal(valuesFileContent, &valuesMap)
		if err != nil {
			return fmt.Errorf("failed to unmarshal values file: %w", err)
		}

		policy, ok := valuesMap["policy"].(map[string]any)
		if !ok {
			return fmt.Errorf("policy not found in blueprint values")
		}

		blueprintsName, ok := valuesMap["blueprint"].(string)
		if !ok {
			return fmt.Errorf("blueprint not found in blueprint values")
		}

		policyName, ok := policy["policy_name"].(string)
		if !ok {
			return fmt.Errorf("policy_name not found in blueprint values")
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

		_, err = client.Apply(context.Background(), &cloudv1.ApplyRequest{
			Blueprint: &cloudv1.Blueprint{
				PolicyName:     policyName,
				Version:        version,
				Values:         valuesFileContent,
				BlueprintsName: blueprintsName,
			},
		})
		if err != nil {
			return err
		}

		return nil
	},
}
