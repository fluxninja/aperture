package blueprints

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fluxninja/aperture/pkg/log"
	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
)

func init() {
	dashboardCmd.Flags().StringVar(&policyType, "policy_type", "", "Type of policy to generate e.g. static-rate-limiting, latency-aimd-concurrency-limiting")
	dashboardCmd.Flags().StringVar(&outputDir, "output_dir", "", "Directory path where the generated manifests will be stored. If not provided, will be printed on console")
	dashboardCmd.Flags().StringVar(&valuesFile, "values_file", "", "Path to the values file for blueprints input")
}

var dashboardCmd = &cobra.Command{
	Use:           "generate-dashboard",
	Short:         "Generate dashboard",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// check if policy or cr is provided
		if policyType == "" {
			return fmt.Errorf("--policy_type must be provided")
		}

		if err := getValuesFilePath(); err != nil {
			return err
		}
		log.Info().Msgf("Using '%s' as values file\n", valuesFile)

		if err := validatePolicyType(cmd, args); err != nil {
			return err
		}

		apertureBlueprintsDir := filepath.Join(blueprintsDir, blueprintsVersion)

		vm := jsonnet.MakeVM()

		vm.Importer(&jsonnet.FileImporter{
			JPaths: []string{apertureBlueprintsDir},
		})

		jsonStr, err := vm.EvaluateAnonymousSnippet("dashboard.libsonnet", fmt.Sprintf(`
		local dashboard = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/policies/%s/dashboard.libsonnet';
		local config = std.parseYaml(importstr '%s');
		local dashboardResource = dashboard(config.common + config.dashboard).dashboard;

		dashboardResource
		`, policyType, valuesFile))
		if err != nil {
			return err
		}

		if outputDir != "" {
			updatedOutputDir, err := getOutputDir(outputDir)
			if err != nil {
				return err
			}

			dashboardPath := filepath.Join(updatedOutputDir, "dashboard.json")
			err = os.WriteFile(dashboardPath, []byte(jsonStr), 0o600)
			if err != nil {
				return err
			}
			log.Info().Msgf("Stored dashbobard JSON at '%s'\n", dashboardPath)
		} else {
			fmt.Println(jsonStr)
		}

		return nil
	},
}
