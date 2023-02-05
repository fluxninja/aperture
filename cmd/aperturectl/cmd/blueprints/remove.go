package blueprints

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	removeCmd.Flags().BoolVar(&all, "all", false, "remove all versions of Aperture Blueprints")
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

		if all {
			pathToRemove = blueprintsCacheRoot
		} else {
			pathToRemove = blueprintsDir
		}

		err := os.RemoveAll(pathToRemove)
		if err != nil {
			return err
		}

		return nil
	},
}
