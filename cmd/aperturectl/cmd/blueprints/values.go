package blueprints

import (
	"fmt"
	"os"
	"strings"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/spf13/cobra"
)

const (
	valuesFileName = "values.yaml"
	fallbackEditor = "vi"
)

func init() {
	valuesCmd.Flags().StringVar(&blueprintName, "name", "", "Name of the Aperture Blueprint to provide values file for")
	valuesCmd.Flags().StringVar(&valuesFile, "output-file", "", "Path to the output values file")
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
	Example:       `aperturectl blueprints values --name=rate-limiting/base --output-file=values.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if blueprintName == "" {
			return fmt.Errorf("--name must be provided")
		}
		if valuesFile == "" {
			return fmt.Errorf("--output-file must be provided")
		}
		_, _, blueprintsDir, err := utils.Pull(blueprintsURI, blueprintsVersion, blueprints, utils.DefaultBlueprintsRepo, skipPull, true)
		if err != nil {
			return err
		}
		return createValuesFile(blueprintsDir, blueprintName, valuesFile, false)
	},
}
