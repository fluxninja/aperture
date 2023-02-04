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

var cachedAll bool

func init() {
	listCmd.Flags().BoolVar(&cachedAll, "all", false, "show the entire cache of Aperture Blueprints")
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
		if cachedAll {
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

			fmt.Printf("Blueprints for: %s\n", getURI(blueprintsDir))
			for _, policy := range policies {
				fmt.Fprintf(w, "%s\n", policy)
			}

			w.Flush()
		}

		return nil
	},
}

func getBlueprints(blueprintsDir string) ([]string, error) {
	uri := getURI(blueprintsDir)

	if uri != "" {
		// tokenize @ sign and take the first element
		uri = strings.Split(uri, "@")[0]
	}

	policies := []string{}

	blueprintsPath := filepath.Join(blueprintsDir, uri)
	err := filepath.WalkDir(blueprintsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == "config.libsonnet" {
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
	blueprintsURIs, err := os.ReadDir(blueprintsCacheRoot)
	if err != nil {
		return nil, err
	}

	for _, blueprintsURI := range blueprintsURIs {
		if blueprintsURI.IsDir() {
			dir := filepath.Join(blueprintsCacheRoot, blueprintsURI.Name())
			uri := getURI(dir)
			policies, err := getBlueprints(dir)
			if err != nil {
				return nil, err
			}
			blueprintsList[uri] = policies
		}
	}

	return blueprintsList, nil
}
