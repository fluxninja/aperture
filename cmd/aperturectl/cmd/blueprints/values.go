package blueprints

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/pkg/log"
)

const (
	valuesFileName                      = "values.yaml"
	requiredValuesFileName              = "values-required.yaml"
	dynamicConfigValuesFileName         = "dynamic-config-values.yaml"
	requiredDynamicConfigValuesFileName = "dynamic-config-values-required.yaml"
)

func init() {
	valuesCmd.Flags().StringVar(&blueprintName, "name", "", "Name of the Aperture Blueprint to provide values file for")
	valuesCmd.Flags().StringVar(&valuesFile, "output-file", "", "Path to the output values file")
	valuesCmd.Flags().BoolVar(&onlyRequired, "only-required", false, "Show only required values")
	valuesCmd.Flags().BoolVar(&dynamicConfig, "dynamic-config", false, "Show dynamic config values instead")
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
		blueprintDir := filepath.Join(blueprintsDir, getRelPath(blueprintsDir))

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

		file := filepath.Join(blueprintDir, blueprintName, valFileName)
		if _, err := os.Stat(file); err != nil {
			return fmt.Errorf("values file not found for the blueprint at: %s", file)
		}
		// make a new copy values.yaml to the output file
		if err := copyFile(file, valuesFile); err != nil {
			return err
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
		cmd := exec.Command(editor, valuesFile)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("error opening file with editor: %s", err)
		}

		log.Info().Msgf("values file for the blueprint %s is available at: %s", blueprintName, valuesFile)
		return nil
	},
}

// copyFile copies a file from src to dst. Any existing file will be overwritten and will not copy file attributes.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return err
}
