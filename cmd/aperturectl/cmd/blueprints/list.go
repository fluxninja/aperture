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

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		blueprintsCacheRoot, blueprintsURIRoot, _, err := utils.Pull(blueprintsURI, blueprintsVersion, blueprints, utils.DefaultBlueprintsRepo, skipPull, true)
		if err != nil {
			return err
		}

		if all {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

			blueprintsList, err := getCachedBlueprints(blueprintsCacheRoot)
			if err != nil {
				return err
			}

			for version, blueprints := range blueprintsList {
				fmt.Printf("%s\n", version)
				for i, blueprint := range blueprints {
					fmt.Fprintf(w, "%s\n", blueprint)
					if i == len(blueprints)-1 {
						fmt.Fprintf(w, "\n")
					}
				}
				w.Flush()
			}

			w.Flush()
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

			blueprints, err := getBlueprints(blueprintsURIRoot, utils.AllowDeprecated)
			if err != nil {
				return err
			}

			for _, blueprint := range blueprints {
				fmt.Fprintf(w, "%s\n", blueprint)
			}

			w.Flush()
		}

		return nil
	},
}

func getBlueprints(blueprintsDir string, listDeprecated bool) ([]string, error) {
	policies := []string{}

	err := symwalk.Walk(blueprintsDir, func(path string, fi fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() && fi.Name() == "bundle.libsonnet" {
			strippedPath := strings.TrimPrefix(path, blueprintsDir)
			strippedPath = strings.TrimSuffix(strippedPath, "/bundle.libsonnet")
			strippedPath = strings.TrimPrefix(strippedPath, "/")
			if !listDeprecated {
				// extract blueprint dir from path
				blueprintDir := strings.TrimSuffix(path, "/bundle.libsonnet")
				ok, _ := utils.IsBlueprintDeprecated(blueprintDir)
				if ok {
					return nil
				}
			}
			policies = append(policies, strippedPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func getCachedBlueprints(blueprintsCacheRoot string) (map[string][]string, error) {
	blueprintsList := map[string][]string{}
	blueprintsURIDirs, err := os.ReadDir(blueprintsCacheRoot)
	if err != nil {
		return nil, err
	}

	for _, blueprintsURIDir := range blueprintsURIDirs {
		if blueprintsURIDir.IsDir() {
			dir := filepath.Join(blueprintsCacheRoot, blueprintsURIDir.Name())
			policies, err := getBlueprints(dir, utils.AllowDeprecated)
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
