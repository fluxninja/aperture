package blueprints

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/flock"
	"github.com/spf13/cobra"
)

const (
	apertureRepo             = "https://github.com/fluxninja/aperture"
	defaultBlueprintsRepo    = "github.com/fluxninja/aperture/blueprints"
	defaultBlueprintsVersion = "latest"
	sourceFilename           = ".source"
	versionFilename          = ".version"
	relPathFilename          = ".relpath"
)

var (
	// Location of cache for blueprints.
	blueprintsCacheRoot string
	// Location of blueprints in disk (e.g. within cache or custom location).
	blueprintsDir string

	// Args for `blueprints`.
	blueprintsURI     string
	blueprintsVersion string
	lock              *flock.Flock
)

func init() {
	BlueprintsCmd.PersistentFlags().StringVar(&blueprintsVersion, "version", defaultBlueprintsVersion, "Version of official Aperture Blueprints, e.g. latest. This field should not be provided when the URI is provided")
	BlueprintsCmd.PersistentFlags().StringVar(&blueprintsURI, "uri", "", "URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@main. This field should not be provided when the Version is provided.")

	BlueprintsCmd.AddCommand(pullCmd)
	BlueprintsCmd.AddCommand(listCmd)
	BlueprintsCmd.AddCommand(removeCmd)
	BlueprintsCmd.AddCommand(generateCmd)
}

// BlueprintsCmd is the root command for blueprints.
var BlueprintsCmd = &cobra.Command{
	Use:   "blueprints",
	Short: "Aperture Blueprints",
	Long: `
Use this command to pull, list, remove and generate Aperture Policy resources using the Aperture Blueprints.`,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		blueprintsCacheRoot = filepath.Join(userHomeDir, ".aperturectl", "blueprints")
		err = os.MkdirAll(blueprintsCacheRoot, os.ModePerm)
		if err != nil {
			return err
		}
		// either the URI or version is set, not both
		if blueprintsURI != "" && blueprintsVersion != defaultBlueprintsVersion {
			return errors.New("either the URI or version should be set, not both")
		}

		// set the URI
		if blueprintsURI == "" {
			if blueprintsVersion == defaultBlueprintsVersion {
				blueprintsVersion, err = resolveLatestVersion()
				if err != nil {
					return err
				}
			}
			blueprintsURI = fmt.Sprintf("%s@%s", defaultBlueprintsRepo, blueprintsVersion)
		}

		// convert the URI to a local dir name which is disk friendly
		dirName := strings.ReplaceAll(blueprintsURI, "/", "_")
		blueprintsDir = filepath.Join(blueprintsCacheRoot, dirName)
		err = os.MkdirAll(blueprintsDir, os.ModePerm)
		if err != nil {
			return err
		}
		// lock blueprintsDir to prevent concurrent access using flock package
		lock = flock.New(filepath.Join(blueprintsDir, ".flock"))
		// use TryLockContext to try locking every 10sec
		locked, err := lock.TryLockContext(context.Background(), 10)
		if err != nil {
			return err
		}
		if !locked {
			return errors.New("unable to acquire lock on blueprints directory")
		}
		return nil
	},
	PersistentPostRunE: func(_ *cobra.Command, _ []string) error {
		if lock != nil {
			return lock.Unlock()
		}
		return nil
	},
}
