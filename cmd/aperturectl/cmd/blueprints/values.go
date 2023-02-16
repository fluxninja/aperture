package blueprints

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/pkg/log"
)

const (
	valuesFileName                      = "values.yaml"
	requiredValuesFileName              = "values-required.yaml"
	dynamicConfigValuesFileName         = "dynamic-config-values.yaml"
	requiredDynamicConfigValuesFileName = "dynamic-config-values-required.yaml"
	fallbackEditor                      = "vi"
)

func init() {
	valuesCmd.Flags().StringVar(&blueprintName, "name", "", "Name of the Aperture Blueprint to provide values file for")
	valuesCmd.Flags().StringVar(&valuesFile, "output-file", "", "Path to the output values file")
	valuesCmd.Flags().BoolVar(&onlyRequired, "only-required", false, "Show only required values")
	valuesCmd.Flags().BoolVar(&dynamicConfig, "dynamic-config", false, "Show dynamic config values instead")
	valuesCmd.Flags().BoolVar(&noYAMLModeline, "no-yaml-modeline", false, "Do not add YAML language server modeline to generated YAML files")
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
	Short: "Provide values file for a given Aperture Blueprint",
	Long: `
Provides a values file for a given Aperture Blueprint that can be then used to generate policies after customization`,
	SilenceErrors: true,
	Example: `aperturectl blueprints values --name=policies/static-rate-limiting --output-file=values.yaml

aperturectl blueprints values --name=policies/static-rate-limiting --output-file=values.yaml --only-required`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if blueprintName == "" {
			return fmt.Errorf("--name must be provided")
		}
		if valuesFile == "" {
			return fmt.Errorf("--output-file must be provided")
		}
		blueprintsDir = filepath.Join(blueprintsURIRoot, getRelPath(blueprintsURIRoot))

		var valFileName string

		if !dynamicConfig {
			if !onlyRequired {
				valFileName = valuesFileName
			} else {
				valFileName = requiredValuesFileName
			}
		} else {
			if !onlyRequired {
				valFileName = dynamicConfigValuesFileName
			} else {
				valFileName = requiredDynamicConfigValuesFileName
			}
		}

		if onlyRequired {
			valFileName = requiredValuesFileName
		}

		srcValuesFile := filepath.Join(blueprintsDir, blueprintName, valFileName)
		if _, err := os.Stat(srcValuesFile); err != nil {
			return fmt.Errorf("values file not found for the blueprint at: %s", srcValuesFile)
		}

		in, err := os.Open(srcValuesFile)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(valuesFile)
		if err != nil {
			return err
		}
		defer func() {
			cerr := out.Close()
			if err == nil {
				err = cerr
			}
		}()
		// prepend YAML modeline to the file
		if !noYAMLModeline {
			var schemaURL string
			if !dynamicConfig {
				schemaURL = fmt.Sprintf("file:%s", filepath.Join(blueprintsDir, blueprintName, "definitions.json"))
			} else {
				schemaURL = fmt.Sprintf("file:%s", filepath.Join(blueprintsDir, blueprintName, "dynamic-config-definitions.json"))
			}
			_, err = out.WriteString("# yaml-language-server: $schema=" + schemaURL + "\n")
			if err != nil {
				return err
			}
		}

		if _, err = io.Copy(out, in); err != nil {
			return err
		}
		err = out.Sync()
		if err != nil {
			return err
		}

		command, args := getEnvEditorWithFallback()
		args = append(args, valuesFile)
		cmd := exec.Command(command, args...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("error opening file with editor: %s", err)
		}

		log.Info().Msgf("values file for the blueprint %s is available at: %s", blueprintName, valuesFile)
		return nil
	},
}
