package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/jsonnet-bundler/jsonnet-bundler/pkg"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	specv1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"
	"github.com/spf13/cobra"
)

const apertureBlueprintsRepo = "github.com/fluxninja/aperture/blueprints"

var (
	blueprintsDir string

	// Args for `blueprints`.
	blueprintsVersion string

	// Args for `blueprints clean`.
	removeAll bool
)

func init() {
	blueprintsCmd.PersistentFlags().StringVar(&blueprintsVersion, "version", "main", "version of aperture blueprint")
	rootCmd.AddCommand(blueprintsCmd)

	// blueprints pull
	blueprintsCmd.AddCommand(blueprintsPullCmd)

	// blueprints list
	blueprintsCmd.AddCommand(blueprintsListCmd)

	// blueprints remove
	blueprintsRemoveCmd.Flags().BoolVar(&removeAll, "all", false, "remove all versions of aperture blueprints")
	blueprintsCmd.AddCommand(blueprintsRemoveCmd)
}

var blueprintsCmd = &cobra.Command{
	Use:   "blueprints",
	Short: "Manage blueprints",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		blueprintsDir = filepath.Join(userHomeDir, ".aperturectl", "blueprints")
		err = os.MkdirAll(blueprintsDir, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	},
}

var blueprintsPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull a blueprint",
	RunE: func(cmd *cobra.Command, args []string) error {
		apertureBlueprintsDir := filepath.Join(blueprintsDir, blueprintsVersion)
		err := os.MkdirAll(apertureBlueprintsDir, os.ModePerm)
		if err != nil {
			return err
		}

		spec := specv1.New()
		contents, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			return err
		}
		spec.LegacyImports = false
		contents = append(contents, []byte("\n")...)

		filename := filepath.Join(apertureBlueprintsDir, jsonnetfile.File)
		err = os.WriteFile(filename, contents, os.ModePerm)
		if err != nil {
			return err
		}

		jbLockFileBytes, err := os.ReadFile(filepath.Join(apertureBlueprintsDir, jsonnetfile.LockFile))
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		lockFile, err := jsonnetfile.Unmarshal(jbLockFileBytes)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Join(apertureBlueprintsDir, ".tmp"), os.ModePerm)
		if err != nil {
			return err
		}

		uri := fmt.Sprintf("%s@%s", apertureBlueprintsRepo, blueprintsVersion)
		d := deps.Parse(apertureBlueprintsDir, uri)
		if !depEqual(spec.Dependencies[d.Name()], *d) {
			spec.Dependencies[d.Name()] = *d
			delete(lockFile.Dependencies, d.Name())
		}

		locked, err := pkg.Ensure(spec, apertureBlueprintsDir, lockFile.Dependencies)
		if err != nil {
			return err
		}

		err = writeChangedJsonnetFile(contents, &spec, filename)
		if err != nil {
			return err
		}
		err = writeChangedJsonnetFile(jbLockFileBytes, &specv1.JsonnetFile{Dependencies: locked}, filepath.Join(apertureBlueprintsDir, jsonnetfile.LockFile))
		if err != nil {
			return err
		}

		return nil
	},
}

func depEqual(d1, d2 deps.Dependency) bool {
	name := d1.Name() == d2.Name()
	version := d1.Version == d2.Version
	source := reflect.DeepEqual(d1.Source, d2.Source)

	return name && version && source
}

func writeChangedJsonnetFile(originalBytes []byte, modified *specv1.JsonnetFile, path string) error {
	origJsonnetFile, err := jsonnetfile.Unmarshal(originalBytes)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(origJsonnetFile, *modified) {
		return nil
	}

	return writeJSONFile(path, *modified)
}

func writeJSONFile(name string, d interface{}) error {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	b = append(b, []byte("\n")...)

	// nolint: gosec
	return os.WriteFile(name, b, 0o644)
}

var blueprintsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List blueprints",
	RunE: func(cmd *cobra.Command, args []string) error {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Version\tPolicies")
		w.Flush()

		blueprintsContents, err := os.ReadDir(blueprintsDir)
		if err != nil {
			return err
		}
		for _, content := range blueprintsContents {
			if content.IsDir() {
				version := content.Name()
				policies := []string{}

				blueprintsVersionContents, err := os.ReadDir(filepath.Join(blueprintsDir, version, apertureBlueprintsRepo, "lib", "1.0", "policies"))
				if err != nil && !os.IsNotExist(err) {
					return err
				}
				for _, versionContent := range blueprintsVersionContents {
					if versionContent.IsDir() {
						policies = append(policies, versionContent.Name())
					}
				}

				fmt.Fprintf(w, "%s\t%s\n", version, strings.Join(policies, ", "))
			}
		}

		w.Flush()
		return nil
	},
}

var blueprintsRemoveCmd = &cobra.Command{
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
