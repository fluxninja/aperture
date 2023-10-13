package blueprints

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/spf13/cobra"

	"github.com/fluxninja/aperture/v2/cmd/aperturectl/cmd/utils"
	"github.com/fluxninja/aperture/v2/pkg/log"
)

var lock *flock.Flock

var pullCmd = &cobra.Command{
	Use:           "pull",
	Short:         "Pull Aperture Blueprints",
	Long:          `Use this command to pull the Aperture Blueprints in local system to use for generating Aperture Policies and Grafana Dashboards.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Example: `aperturectl blueprints pull

aperturectl blueprints pull --version latest`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _, _, err := pull(blueprintsURI, blueprintsVersion, true)
		if err != nil {
			return err
		}
		return nil
	},
}

// Pull pulls the blueprints from the given URI and version.
func pull(blueprintsURI string, blueprintsVersion string, localAllowed bool) (string, string, string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", "", "", err
	}

	blueprintsCacheRoot := filepath.Join(userHomeDir, utils.AperturectlRootDir, "blueprints")
	err = os.MkdirAll(blueprintsCacheRoot, os.ModePerm)
	if err != nil {
		return "", "", "", err
	}

	// either the URI or version is set, not both
	if blueprintsURI != "" && blueprintsVersion != "" {
		return blueprintsCacheRoot, "", "", errors.New("either the URI or version should be set, not both, uri: " + blueprintsURI + ", version: " + blueprintsVersion)
	}

	// set the URI
	if blueprintsURI == "" {
		if blueprintsVersion == "" {
			blueprintsVersion = LatestTag
		}
		blueprintsURI = fmt.Sprintf("%s@%s", DefaultBlueprintsRepo, blueprintsVersion)
	} else {
		// uri can be a file or url
		// first detect if it's a local path
		if _, err = os.Stat(blueprintsURI); err == nil {
			if !localAllowed {
				return blueprintsCacheRoot, "", "", errors.New("local paths are not allowed as blueprints URI")
			}
			// path exists
			blueprintsURI, err = filepath.Abs(blueprintsURI)
			if err != nil {
				return blueprintsCacheRoot, "", "", err
			}
		} else {
			// try to parse as url
			var blueprintsURL *url.URL
			blueprintsURL, err = url.Parse(blueprintsURI)
			if err != nil {
				return blueprintsCacheRoot, "", "", err
			}
			blueprintsURI = blueprintsURL.String()
		}
	}

	// convert the URI to a local dir name which is disk friendly
	dirName := strings.ReplaceAll(blueprintsURI, "/", "_")
	blueprintsURIRoot := filepath.Join(blueprintsCacheRoot, dirName)
	err = os.MkdirAll(blueprintsURIRoot, os.ModePerm)
	if err != nil {
		return blueprintsCacheRoot, "", "", err
	}

	// get a file lock on the blueprintsURIRoot
	lock = flock.New(filepath.Join(blueprintsURIRoot, "lock"))
	locked, err := lock.TryLockContext(context.Background(), time.Millisecond*100)
	if err != nil {
		return blueprintsCacheRoot, blueprintsURIRoot, "", err
	}
	if !locked {
		return blueprintsCacheRoot, blueprintsURIRoot, "", errors.New("could not get lock on: " + blueprintsURIRoot)
	}

	// pull the latest blueprints based on skipPull and whether child command is remove
	if !skipPull {
		err = utils.PullSource(blueprintsURIRoot, blueprintsURI)
		if err != nil {
			return blueprintsCacheRoot, blueprintsURIRoot, "", err
		}
	} else {
		log.Trace().Msg("Skipping pulling blueprints")
	}

	blueprintsDir := filepath.Join(blueprintsURIRoot, utils.GetRelPath(blueprintsURIRoot))
	// resolve symlink
	blueprintsDir, err = filepath.EvalSymlinks(blueprintsDir)
	if err != nil {
		return blueprintsCacheRoot, blueprintsURIRoot, "", err
	}

	// unlock the file lock on the blueprintsURIRoot
	err = lock.Unlock()
	if err != nil {
		return blueprintsCacheRoot, blueprintsURIRoot, blueprintsDir, err
	}

	return blueprintsCacheRoot, blueprintsURIRoot, blueprintsDir, nil
}
