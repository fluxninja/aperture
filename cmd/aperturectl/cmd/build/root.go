package build

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	defaultApertureRepo    = "github.com/fluxninja/aperture"
	defaultApertureVersion = "latest"
)

var (
	builderCacheRoot string
	builderURIRoot   string
	builderDir       string

	apertureURI     string
	apertureVersion string

	skipPull bool
)

func init() {
	BuildCmd.PersistentFlags().StringVar(&apertureVersion, "version", defaultApertureVersion, "Version of Aperture, e.g. latest. This field should not be provided when the URI is provided")
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
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// get aperture repository and save it to aperturectl root directory
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		builderCacheRoot = filepath.Join(userHomeDir, ".aperturectl", "build")
		err = os.MkdirAll(builderCacheRoot, os.ModePerm)
		if err != nil {
			return err
		}
		// either the URI or version is set, not both
		if apertureURI != "" && apertureVersion != defaultApertureVersion {
			return errors.New("either the URI or version should be set, not both")
		}

		// set the URI
		if apertureURI == "" {
			if apertureVersion == defaultApertureVersion {
				apertureVersion, err = utils.ResolveLatestVersion()
				if err != nil {
					return err
				}
			}
			apertureURI = fmt.Sprintf("%s@%s", defaultApertureRepo, apertureVersion)
		} else {
			apertureURI, err = filepath.Abs(apertureURI)
			if err != nil {
				return err
			}
		}

		// convert the URI to a local dir name which is disk friendly
		dirName := strings.ReplaceAll(apertureURI, "/", "_")
		builderURIRoot = filepath.Join(builderCacheRoot, dirName)
		err = os.MkdirAll(builderURIRoot, os.ModePerm)
		if err != nil {
			return err
		}

		// pull the latest blueprints based on skipPull and whether cmd is remove
		if !skipPull && cmd.Use != "remove" {
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
}
