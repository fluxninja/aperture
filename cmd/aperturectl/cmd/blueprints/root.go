package blueprints

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fluxninja/aperture/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/pkg/log"
	"github.com/gofrs/flock"
	"github.com/spf13/cobra"
)

const (
	defaultBlueprintsRepo    = "github.com/fluxninja/aperture/blueprints"
	defaultBlueprintsVersion = "latest"
	sourceFilename           = ".source"
	versionFilename          = ".version"
	relPathFilename          = ".relpath"
	lockFilename             = ".flock"
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
	BlueprintsCmd.PersistentFlags().StringVar(&blueprintsURI, "uri", "", "URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@latest. This field should not be provided when the Version is provided.")
	BlueprintsCmd.PersistentFlags().BoolVar(&skipPull, "skip-pull", false, "Skip pulling the blueprints update.")

	BlueprintsCmd.AddCommand(pullCmd)
	BlueprintsCmd.AddCommand(listCmd)
	BlueprintsCmd.AddCommand(removeCmd)
	BlueprintsCmd.AddCommand(generateCmd)
	BlueprintsCmd.AddCommand(valuesCmd)
}

// BlueprintsCmd is the root command for blueprints.
var BlueprintsCmd = &cobra.Command{
	Use:   "blueprints",
	Short: "Aperture Blueprints",
	Long: `
Use this command to pull, list, remove and generate Aperture Policy resources using the Aperture Blueprints.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
				blueprintsVersion, err = utils.ResolveLatestVersion()
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
		lock = flock.New(filepath.Join(blueprintsDir, lockFilename))

		// pull the latest blueprints based on skipPull and whether cmd is remove
		if !skipPull && cmd.Use != "remove" {
			err = pullFunc(cmd, args)
			if err != nil {
				return err
			}
		} else {
			log.Debug().Msg("skipping pulling blueprints")
		}
		return nil
	},
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		unlock()
	},
}

func writerLock() error {
	// Get writer lock
	locked, err := lock.TryLockContext(context.Background(), 10)
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("unable to acquire lock on blueprints directory")
	}
	return nil
}

func readerLock() error {
	// Get reader lock
	locked, err := lock.TryRLockContext(context.Background(), 10)
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("unable to acquire lock on blueprints directory")
	}
	return nil
}

func unlock() {
	err := lock.Unlock()
	if err != nil {
		log.Error().Err(err).Msg("unable to release lock on blueprints directory")
		// try resetting lock by removing lockfile
		os.Remove(filepath.Join(blueprintsDir, lockFilename))
	}
}
