package blueprints

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	dynamicConfigValuesFileName = "dynamic-config-values.yaml"
)

func init() {
	dynamicValuesCmd.Flags().StringVar(&blueprintName, "name", "", "Name of the Aperture Blueprint to provide values file for")
	dynamicValuesCmd.Flags().StringVar(&valuesFile, "output-file", "", "Path to the output values file")
	dynamicValuesCmd.Flags().BoolVar(&noYAMLModeline, "no-yaml-modeline", false, "Do not add YAML language server modeline to generated YAML files")
	dynamicValuesCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing values file")
}

var dynamicValuesCmd = &cobra.Command{
	Use:   "dynamic-values",
	Short: "Create dynamic values file for a given Aperture Blueprint",
	Long: `
Provides a dynamic values file for a given Aperture Blueprint that can be then used to generate policies after customization`,
	SilenceErrors: true,
	Example:       `aperturectl blueprints dynamic-values --name=policies/rate-limiting --output-file=values.yaml`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if blueprintName == "" {
			return fmt.Errorf("--name must be provided")
		}
		if valuesFile == "" {
			return fmt.Errorf("--output-file must be provided")
		}
		return createValuesFile(blueprintName, valuesFile, true)
	},
}
