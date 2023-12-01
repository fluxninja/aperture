package blueprints

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	cloudv1 "github.com/fluxninja/aperture/api/v2/gen/proto/go/aperture/cloud/v1"
	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/blueprints"
)

func init() {
	BlueprintsApplyCmd.Flags().StringVar(&valuesFile, "values-file", "", "Values file to use for blueprint")
	BlueprintsApplyCmd.Flags().StringVar(&valuesDir, "values-dir", "", "Path to directory containing blueprint values files")
}

// BlueprintsApplyCmd is the command to apply a blueprint from the Cloud Controller.
var BlueprintsApplyCmd = &cobra.Command{
	Use:           "apply",
	Short:         "Cloud Blueprints Apply",
	Long:          `Apply cloud blueprint.`,
	SilenceErrors: true,
	Example:       `aperturectl cloud blueprints apply --value-file=values.yaml`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if valuesFile == "" && valuesDir == "" {
			return fmt.Errorf("either --values-file or --dir is required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle directory with multiple values files
		if valuesDir != "" {
			err := filepath.Walk(valuesDir, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					if err := applyBlueprint(path); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("failed to apply blueprints: %w", err)
			}
			return nil
		}

		// Handle single values file
		if err := applyBlueprint(valuesFile); err != nil {
			return fmt.Errorf("failed to apply blueprint: %w", err)
		}

		return nil
	},
}

func applyBlueprint(valuesFile string) error {
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
}
