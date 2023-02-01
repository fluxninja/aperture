package blueprints

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	BlueprintsCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List blueprints",
	RunE: func(cmd *cobra.Command, args []string) error {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Version\tPolicies")

		blueprintsList, err := getBlueprintsList()
		if err != nil {
			return err
		}

		for version, policies := range blueprintsList {
			fmt.Fprintf(w, "%s\t%s\n", version, strings.Join(policies, ", "))
		}

		w.Flush()
		return nil
	},
}

func getBlueprintsList() (map[string][]string, error) {
	blueprints := map[string][]string{}
	blueprintsContents, err := os.ReadDir(blueprintsDir)
	if err != nil {
		return nil, err
	}
	for _, content := range blueprintsContents {
		if content.IsDir() {
			version := content.Name()
			policies := []string{}

			blueprintsVersionContents, err := os.ReadDir(filepath.Join(blueprintsDir, version, apertureBlueprintsURI, "lib", "1.0", "policies"))
			if err != nil && !os.IsNotExist(err) {
				return nil, err
			}
			for _, versionContent := range blueprintsVersionContents {
				if versionContent.IsDir() {
					policies = append(policies, versionContent.Name())
				}
			}

			blueprints[version] = policies
		}
	}

	return blueprints, nil
}
