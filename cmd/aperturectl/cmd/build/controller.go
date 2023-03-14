package build

import "github.com/spf13/cobra"

func init() {
	controllerCmd.Flags().StringVarP(&buildConfigFile, "config", "c", "", "path to the build configuration file")
	controllerCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "path to the output directory")
}

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Build controller binary for Aperture",
	Long:  "Build controller binary for Aperture",
	RunE:  buildRunE("aperture-controller"),
}
