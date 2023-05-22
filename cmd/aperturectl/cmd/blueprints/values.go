package blueprints

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	valuesFileName         = "values.yaml"
	requiredValuesFileName = "values-required.yaml"
	fallbackEditor         = "vi"
)

func init() {
	valuesCmd.Flags().StringVar(&blueprintName, "name", "", "Name of the Aperture Blueprint to provide values file for")
	valuesCmd.Flags().StringVar(&valuesFile, "output-file", "", "Path to the output values file")
	valuesCmd.Flags().BoolVar(&onlyRequired, "only-required", false, "Show only required values")
	valuesCmd.Flags().BoolVar(&noYAMLModeline, "no-yaml-modeline", false, "Do not add YAML language server modeline to generated YAML files")
	valuesCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing values file")
}

func getEnvEditorWithFallback() (string, []string) {
	visual := os.Getenv("VISUAL")
	var found string
	if visual != "" {
		found = visual
	}

	if found == "" {
		editor := os.Getenv("EDITOR")
		if editor != "" {
			found = editor
		} else {
			found = fallbackEditor
		}
	}

	parts := strings.Split(found, " ")
	return parts[0], parts[1:]
}

var valuesCmd = &cobra.Command{
	Use:   "values",
	Short: "Create values file for a given Aperture Blueprint",
	Long: `
Provides a values file for a given Aperture Blueprint that can be then used to generate policies after customization`,
	SilenceErrors: true,
	Example: `aperturectl blueprints values --name=policies/rate-limiting --output-file=values.yaml

aperturectl blueprints values --name=policies/rate-limiting --output-file=values.yaml --only-required`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if blueprintName == "" {
			return fmt.Errorf("--name must be provided")
		}
		if valuesFile == "" {
			return fmt.Errorf("--output-file must be provided")
		}
		return createValuesFile(blueprintName, valuesFile, false)
	},
}
