package blueprints

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var removeAll bool

func init() {
	removeCmd.Flags().BoolVar(&removeAll, "all", false, "remove all versions of Aperture Blueprints")
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a Blueprint",
	Long: `
Use this command to remove a pulled Aperture Blueprint from local system.`,
	Example: `aperturectl blueprints remove

aperturectl blueprints remove --version v0.22.0

aperturectl blueprints remove --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pathToRemove := ""

		if removeAll {
			pathToRemove = blueprintsDir
		} else {
			pathToRemove = filepath.Join(blueprintsDir, blueprintsVersion)
		}

		err := os.RemoveAll(pathToRemove)
		if err != nil {
			return err
		}

		return nil
	},
}
