package blueprints

import (
	"fmt"
	"os"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
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
	Example: fmt.Sprintf(`aperturectl blueprints remove

aperturectl blueprints remove --version v%s

aperturectl blueprints remove --all`, utils.Version),
	RunE: func(_ *cobra.Command, _ []string) error {
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
