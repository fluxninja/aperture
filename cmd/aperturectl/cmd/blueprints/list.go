package blueprints

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/facebookgo/symwalk"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().BoolVar(&all, "all", false, "show the entire cache of Aperture Blueprints")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Aperture Blueprints",
	Long: `
Use this command to list the Aperture Blueprints which are already pulled and available in local system.`,
	Example: `aperturectl blueprints list

aperturectl blueprints list --version v0.22.0

aperturectl blueprints list --all`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if all {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
			fmt.Fprintln(w, "URI\tPOLICIES")

			blueprintsList, err := getCachedBlueprints()
			if err != nil {
				return err
			}

			for version, policies := range blueprintsList {
				fmt.Fprintf(w, "%s\t%s\n", version, strings.Join(policies, ", "))
			}

			w.Flush()
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

			policies, err := getBlueprints(blueprintsDir)
			if err != nil {
				return err
			}

			for _, policy := range policies {
				fmt.Fprintf(w, "%s\n", policy)
			}

			w.Flush()
		}

		return nil
	},
}

func getBlueprints(blueprintsDir string) ([]string, error) {
	relPath := getRelPath(blueprintsDir)

	policies := []string{}

	blueprintsPath := filepath.Join(blueprintsDir, relPath)
	err := symwalk.Walk(blueprintsPath, func(path string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() && fi.Name() == "config.libsonnet" {
			strippedPath := strings.TrimPrefix(path, blueprintsPath)
			strippedPath = strings.TrimSuffix(strippedPath, "/config.libsonnet")
			strippedPath = strings.TrimPrefix(strippedPath, "/")
			policies = append(policies, strippedPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func getCachedBlueprints() (map[string][]string, error) {
	blueprintsList := map[string][]string{}
	blueprintsDirs, err := os.ReadDir(blueprintsCacheRoot)
	if err != nil {
		return nil, err
	}

	for _, blueprintsDir := range blueprintsDirs {
		if blueprintsDir.IsDir() {
			dir := filepath.Join(blueprintsCacheRoot, blueprintsDir.Name())
			policies, err := getBlueprints(dir)
			if err != nil {
				return nil, err
			}
			source := getSource(dir)
			blueprintsList[source] = policies
		}
	}

	return blueprintsList, nil
}
