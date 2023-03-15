package build

import (
	"github.com/spf13/cobra"
)

func init() {
	agentCmd.Flags().StringVarP(&buildConfigFile, "config", "c", "", "path to the build configuration file")
	agentCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "path to the output directory")
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Build agent binary for Aperture",
	Long:  "Build agent binary for Aperture",
	RunE:  buildRunE("aperture-agent"),
}
