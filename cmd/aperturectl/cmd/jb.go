package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/jsonnet-bundler/jsonnet-bundler/pkg"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	specv1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"
	"github.com/spf13/cobra"
)

var (
	vendor     string
	uris       []string
	single     bool
	legacyName string
	workDir    string
)

func init() {
	compileCmd.PersistentFlags().StringVar(&vendor, "vendor", "vendor", "path to the directory used to cache packages in")
	rootCmd.AddCommand(jbCmd)

	jbInstallCmd.Flags().StringSliceVar(&uris, "uri", []string{}, "URI of the package to install (URLs or file paths)")
	jbInstallCmd.Flags().BoolVar(&single, "single", false, "install package without dependencies")
	jbInstallCmd.Flags().StringVar(&legacyName, "legacy-name", "", "set legacy name")
	jbCmd.AddCommand(jbInstallCmd)
}

var jbCmd = &cobra.Command{
	Use:   "jb",
	Short: "Jsonnet bundler",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		workDir = wd

		vendor = filepath.Clean(vendor)

		return nil
	},
}

var jbInstallCmd = &cobra.Command{
	Use:           "install",
	Short:         "Install a new dependency",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if workDir == "" {
			workDir = "."
		}

		jbFileBytes, err := os.ReadFile(filepath.Join(workDir, jsonnetfile.File))
		if err != nil {
			return err
		}

		jsonnetFile, err := jsonnetfile.Unmarshal(jbFileBytes)
		if err != nil {
			return err
		}

		jbLockFileBytes, err := os.ReadFile(filepath.Join(workDir, jsonnetfile.LockFile))
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		lockFile, err := jsonnetfile.Unmarshal(jbLockFileBytes)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Join(workDir, vendor, ".tmp"), os.ModePerm)
		if err != nil {
			return err
		}

		if len(uris) > 0 && legacyName != "" {
			return errors.New("cannot specify both --uri and --legacy-name")
		}

		for _, u := range uris {
			d := deps.Parse(workDir, u)
			if d == nil {
				return fmt.Errorf("unable to parse package URI %s", u)
			}

			if single {
				d.Single = true
			}

			if legacyName != "" {
				d.LegacyNameCompat = legacyName
			}

			if !depEqual(jsonnetFile.Dependencies[d.Name()], *d) {
				// the dep passed on the cli is different from the jsonnetFile
				jsonnetFile.Dependencies[d.Name()] = *d

				// we want to install the passed version (ignore the lock)
				delete(lockFile.Dependencies, d.Name())
			}
		}

		jsonnetPkgHomeDir := filepath.Join(workDir, vendor)
		locked, err := pkg.Ensure(jsonnetFile, jsonnetPkgHomeDir, lockFile.Dependencies)
		if err != nil {
			return err
		}

		pkg.CleanLegacyName(jsonnetFile.Dependencies)
		err = writeChangedJsonnetFile(jbFileBytes, &jsonnetFile, filepath.Join(workDir, jsonnetfile.File))
		if err != nil {
			return err
		}
		err = writeChangedJsonnetFile(jbLockFileBytes, &specv1.JsonnetFile{Dependencies: locked}, filepath.Join(workDir, jsonnetfile.LockFile))
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
