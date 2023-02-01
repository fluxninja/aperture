package blueprints

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-jsonnet"
	"github.com/spf13/cobra"
)

var generateDashboardCmd = &cobra.Command{
	Use:           "dashboard",
	Short:         "Generate dashboard",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if blueprintName == "" {
			return fmt.Errorf("--name must be provided")
		}

		if err := blueprintExists(blueprintsVersion, blueprintName); err != nil {
			return err
		}

		apertureBlueprintsDir := filepath.Join(blueprintsDir, blueprintsVersion)

		vm := jsonnet.MakeVM()

		vm.Importer(&jsonnet.FileImporter{
			JPaths: []string{apertureBlueprintsDir},
		})

		jsonStr, err := vm.EvaluateAnonymousSnippet("dashboard.libsonnet", fmt.Sprintf(`
		local dashboard = import 'github.com/fluxninja/aperture/blueprints/lib/1.0/%s/dashboard.libsonnet';
		local config = std.parseYaml(importstr '%s');
		local dashboardResource = dashboard(config.common + config.dashboard).dashboard;

		dashboardResource
		`, blueprintName, valuesFile))
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
		} else {
			fmt.Println(jsonStr)
		}

		return nil
	},
}
