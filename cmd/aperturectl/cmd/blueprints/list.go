package blueprints

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listAll bool

func init() {
	listCmd.Flags().BoolVar(&listAll, "all", false, "list all versions of aperture blueprints")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Aperture Blueprints",
	Long: `
Use this command to list the Aperture Blueprints which are already pulled and available in local system.`,
	Example: `aperturectl blueprints list

aperturectl blueprints list --version v0.22.0

aperturectl blueprints list --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if listAll {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
			fmt.Fprintln(w, "VERSION\tPOLICIES")

			blueprintsList, err := getBlueprints()
			if err != nil {
				return err
			}

			for version, policies := range blueprintsList {
				fmt.Fprintf(w, "%s\t%s\n", version, strings.Join(policies, ", "))
			}

			w.Flush()
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

			policies, err := getBlueprintsByVersion(blueprintsVersion)
			if err != nil {
				return err
			}

			fmt.Printf("Blueprints: %s\n", blueprintsVersion)
			for _, policy := range policies {
				fmt.Fprintf(w, "%s\n", policy)
			}

			w.Flush()
		}

		return nil
	},
}

func getBlueprintsByVersion(v string) ([]string, error) {
	const libSubPath = "lib/1.0"

	policies := []string{}

	libPath := filepath.Join(blueprintsDir, v, apertureBlueprintsURI, libSubPath)
	err := filepath.WalkDir(libPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "config.libsonnet" {
			strippedPath := strings.TrimPrefix(path, libPath)
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

func getBlueprints() (map[string][]string, error) {
	blueprintsList := map[string][]string{}
	blueprintsVersions, err := os.ReadDir(blueprintsDir)
	if err != nil {
		return nil, err
	}

	for _, blueprintsVersion := range blueprintsVersions {
		if blueprintsVersion.IsDir() {
			v := blueprintsVersion.Name()
			policies, err := getBlueprintsByVersion(v)
			if err != nil {
				return nil, err
			}
			blueprintsList[v] = policies
		}
	}

	return blueprintsList, nil
}
