package blueprints

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var removeAll bool

func init() {
	removeCmd.Flags().BoolVar(&removeAll, "all", false, "remove all versions of aperture blueprints")
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a blueprint",
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
