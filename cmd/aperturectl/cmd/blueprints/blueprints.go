package blueprints

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const apertureBlueprintsRepo = "github.com/fluxninja/aperture/blueprints"

var (
	blueprintsDir string

	// Args for `blueprints pull`.
	blueprintsVersion string
)

func init() {
	BlueprintsCmd.PersistentFlags().StringVar(&blueprintsVersion, "version", "main", "version of aperture blueprint")

	BlueprintsCmd.AddCommand(pullCmd)
	BlueprintsCmd.AddCommand(policyCmd)
	BlueprintsCmd.AddCommand(listCmd)
	BlueprintsCmd.AddCommand(removeCmd)
}

var BlueprintsCmd = &cobra.Command{
	Use:   "blueprints",
	Short: "Manage blueprints",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		blueprintsDir = filepath.Join(userHomeDir, ".aperturectl", "blueprints")
		err = os.MkdirAll(blueprintsDir, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	},
}
