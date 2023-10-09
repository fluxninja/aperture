package build

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

// builder builds agent or controller binaries.
// it reads a yaml configuration file that contains:
// list of extension modules to be added to the final binary
//   extension modules have path to go.mod and import path
//   bundled extensions are local and must use replace directives to point to local paths
//   each extension module implements PlatformOptions(), AgentOptions(), ControllerOptions() interfaces that are called by the main binary
//   for each module we generate code that calls the above functions and adds the returned options to the final binary
//   for each module we also add the module's go.mod to the final go.mod and replace directives
// replace directives to be added to the final go.mod
// ldflags to be added to the final binary that set version and other build-time variables

const (
	defaultApertureRepo = "github.com/fluxninja/aperture"
	latestTag           = "latest"
)

var (
	builderCacheRoot string
	builderURIRoot   string
	builderDir       string

	apertureURI     string
	apertureVersion string

	skipPull bool
	lock     *flock.Flock
)

func init() {
	BuildCmd.PersistentFlags().StringVar(&apertureVersion, "version", "", "Version of Aperture, e.g. latest. This field should not be provided when the URI is provided")
	BuildCmd.PersistentFlags().StringVar(&apertureURI, "uri", "", "URI of Aperture repository, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture@latest. This field should not be provided when the Version is provided.")
	BuildCmd.PersistentFlags().BoolVar(&skipPull, "skip-pull", false, "Skip pulling the repository update.")
	BuildCmd.AddCommand(agentCmd)
	BuildCmd.AddCommand(controllerCmd)
}

// BuildCmd is the root command for the build command.
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds the agent and controller binaries",
	Long:  "Builds the agent and controller binaries",
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		// get aperture repository and save it to aperturectl root directory
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		builderCacheRoot = filepath.Join(userHomeDir, utils.AperturectlRootDir, utils.BuilderCacheRoot)
		err = os.MkdirAll(builderCacheRoot, os.ModePerm)
		if err != nil {
			return err
		}
		// either the URI or version is set, not both
		if apertureURI != "" && apertureVersion != "" {
			return errors.New("either the URI or version should be set, not both")
		}

		// set the URI
		if apertureURI == "" {
			if apertureVersion == "" {
				apertureVersion = latestTag
			}
			apertureURI = fmt.Sprintf("%s@%s", defaultApertureRepo, apertureVersion)
		} else {
			// uri can be a file or url
			// first detect if it's a local path
			if _, err = os.Stat(apertureURI); err == nil {
				// path exists
				apertureURI, err = filepath.Abs(apertureURI)
				if err != nil {
					return err
				}
			} else {
				// try to parse as url
				var apertureURL *url.URL
				apertureURL, err = url.Parse(apertureURI)
				if err != nil {
					return err
				}
				apertureURI = apertureURL.String()
			}
		}

		// convert the URI to a local dir name which is disk friendly
		dirName := strings.ReplaceAll(apertureURI, "/", "_")
		builderURIRoot = filepath.Join(builderCacheRoot, dirName)
		err = os.MkdirAll(builderURIRoot, os.ModePerm)
		if err != nil {
			return err
		}

		// get a file lock on the builderURIRoot
		lock = flock.New(filepath.Join(builderURIRoot, "lock"))
		locked, err := lock.TryLockContext(context.Background(), time.Millisecond*100)
		if err != nil {
			return err
		}
		if !locked {
			return errors.New("could not get lock on: " + builderURIRoot)
		}

		// pull the latest blueprints based on skipPull and whether cmd is remove
		if !skipPull {
			err = utils.PullSource(builderURIRoot, apertureURI)
			if err != nil {
				return err
			}
		} else {
			log.Debug().Msg("skipping pulling aperture repository")
		}

		builderDir = filepath.Join(builderURIRoot, utils.GetRelPath(builderURIRoot))
		// if builderDir is a symlink, resolve it
		builderDir, err = filepath.EvalSymlinks(builderDir)
		if err != nil {
			return err
		}
		return nil
	},
	PersistentPostRunE: func(_ *cobra.Command, _ []string) error {
		if lock == nil {
			return nil
		}
		// release the lock
		return lock.Unlock()
	},
}
