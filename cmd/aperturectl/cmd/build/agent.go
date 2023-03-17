package build

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	agentCmd.Flags().StringVarP(&buildConfigFile, "config", "c", "", "path to the build configuration file (default: build-config.yaml in the main package directory)")
	agentCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "path to the output directory (default: current directory)")
}

var agentCmd = &cobra.Command{
	Use:     "agent",
	Short:   "Build agent binary for Aperture",
	Long:    "Build agent binary for Aperture",
	Example: fmt.Sprintf(exampleFmt, "agent", "agent"),
	RunE:    buildRunE("aperture-agent"),
}
