package blueprints

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

const (
	defaultBlueprintsRepo    = "github.com/fluxninja/aperture/blueprints"
	defaultBlueprintsVersion = "latest"
)

var (
	// Location of cache for blueprints. E.g. ~/.aperturectl/blueprints.
	blueprintsCacheRoot string
	// Location of blueprints uri within cache. E.g. ~/.aperturectl/blueprints/github.com/fluxninja/aperture/blueprints@latest.
	blueprintsURIRoot string

	// Location of blueprints directory within URI directory. E.g. ~/.aperturectl/blueprints/github.com_fluxninja_aperture_blueprints@v0.26.1/github.com/fluxninja/aperture/blueprints/.
	blueprintsDir string

	// Args for `blueprints`.
	blueprintsURI     string
	blueprintsVersion string
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
	BlueprintsCmd.AddCommand(dynamicValuesCmd)
}

// BlueprintsCmd is the root command for blueprints.
var BlueprintsCmd = &cobra.Command{
	Use:   "blueprints",
	Short: "Aperture Blueprints",
	Long: `
Use this command to pull, list, remove and generate Aperture Policy resources using the Aperture Blueprints.`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
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
		} else {
			var blueprintsURL *url.URL
			blueprintsURL, err = url.Parse(blueprintsURI)
			if err != nil {
				blueprintsURI, err = filepath.Abs(blueprintsURI)
				if err != nil {
					return err
				}
			} else {
				blueprintsURI = blueprintsURL.String()
			}
		}

		// convert the URI to a local dir name which is disk friendly
		dirName := strings.ReplaceAll(blueprintsURI, "/", "_")
		blueprintsURIRoot = filepath.Join(blueprintsCacheRoot, dirName)
		err = os.MkdirAll(blueprintsURIRoot, os.ModePerm)
		if err != nil {
			return err
		}

		// pull the latest blueprints based on skipPull and whether cmd is remove
		if !skipPull && cmd.Use != "remove" {
			err = utils.PullSource(blueprintsURIRoot, blueprintsURI)
			if err != nil {
				return err
			}
		} else {
			log.Debug().Msg("skipping pulling blueprints")
		}

		blueprintsDir = filepath.Join(blueprintsURIRoot, utils.GetRelPath(blueprintsURIRoot))
		// resolve symlink
		blueprintsDir, err = filepath.EvalSymlinks(blueprintsDir)
		if err != nil {
			return err
		}
		return nil
	},
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		utils.Unlock(blueprintsURIRoot)
	},
}
