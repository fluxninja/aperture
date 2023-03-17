package build

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	controllerCmd.Flags().StringVarP(&buildConfigFile, "config", "c", "", "path to the build configuration file (default: build-config.yaml in the main package directory)")
	controllerCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "path to the output directory (default: current directory)")
}

var controllerCmd = &cobra.Command{
	Use:     "controller",
	Short:   "Build controller binary for Aperture",
	Long:    "Build controller binary for Aperture",
	Example: fmt.Sprintf(exampleFmt, "controller", "controller"),
	RunE:    buildRunE("aperture-controller"),
}
