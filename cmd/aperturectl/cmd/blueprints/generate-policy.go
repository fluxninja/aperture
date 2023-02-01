package blueprints

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var generatePolicyCmd = &cobra.Command{
	Use:           "policy",
	Short:         "Generate policy",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := blueprintExists(blueprintsVersion, blueprintName)
		if err != nil {
			return err
		}

		apertureBlueprintsDir := filepath.Join(blueprintsDir, blueprintsVersion)

		vm := jsonnet.MakeVM()
		vm.Importer(&jsonnet.FileImporter{
			JPaths: []string{apertureBlueprintsDir},
		})

		jsonStr, err := vm.EvaluateAnonymousSnippet("policy.libsonnet", fmt.Sprintf(`
		local policy = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/%s/policy.libsonnet';
		local config = std.parseYaml(importstr '%s');
		local policyResource = policy(config.common + config.policy).policyResource;

		policyResource
		`, blueprintName, valuesFile))
		if err != nil {
			return err
		}

		var yamlData map[string]interface{}
		err = json.Unmarshal([]byte(jsonStr), &yamlData)
		if err != nil {
			return err
		}

		yamlBytes, err := yaml.Marshal(yamlData)
		if err != nil {
			return err
		}

		if outputDir != "" {
			updatedOutputDir, er := getOutputDir(outputDir)
			if er != nil {
				return er
			}
			policyFilePath := filepath.Join(updatedOutputDir, "policy.yaml")
			err = os.WriteFile(policyFilePath, yamlBytes, 0o600)
			if err != nil {
				return err
			}
		} else {
			fmt.Println(string(yamlBytes))
		}

		return nil
	},
}
