package blueprints

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/facebookgo/symwalk"
	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
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

aperturectl blueprints list --version latest

aperturectl blueprints list --all`,
	RunE: func(_ *cobra.Command, _ []string) error {
		err := utils.ReaderLock(blueprintsURIRoot)
		if err != nil {
			return err
		}
		defer utils.Unlock(blueprintsURIRoot)
		if all {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

			blueprintsList, err := getCachedBlueprints()
			if err != nil {
				return err
			}

			for version, policies := range blueprintsList {
				fmt.Printf("%s\n", version)
				for i, policy := range policies {
					fmt.Fprintf(w, "%s\n", policy)
					if i == len(policies)-1 {
						fmt.Fprintf(w, "\n")
					}
				}
				w.Flush()
			}

			w.Flush()
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

			policies, err := getBlueprints(blueprintsURIRoot)
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

func getBlueprints(blURIRoot string) ([]string, error) {
	relPath := utils.GetRelPath(blURIRoot)
	policies := []string{}

	blDir := filepath.Join(blURIRoot, relPath)
	err := symwalk.Walk(blDir, func(path string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() && fi.Name() == "config.libsonnet" {
			strippedPath := strings.TrimPrefix(path, blDir)
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
	blueprintsURIDirs, err := os.ReadDir(blueprintsCacheRoot)
	if err != nil {
		return nil, err
	}

	for _, blueprintsURIDir := range blueprintsURIDirs {
		if blueprintsURIDir.IsDir() {
			dir := filepath.Join(blueprintsCacheRoot, blueprintsURIDir.Name())
			policies, err := getBlueprints(dir)
			if err != nil {
				return nil, err
			}
			source := utils.GetSource(dir)
			version := utils.GetVersion(dir)
			if version != "" {
				source = fmt.Sprintf("%s@%s", source, version)
			}
			blueprintsList[source] = policies
		}
	}

	return blueprintsList, nil
}
