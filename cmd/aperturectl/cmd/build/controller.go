package build

import "github.com/spf13/cobra"

func init() {
	controllerCmd.Flags().StringVarP(&buildConfigFile, "config", "c", "", "path to the build configuration file")
	controllerCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "path to the output directory")
}

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Build controller binary for Aperture",
	Long:  "Build controller binary for Aperture",
	Example: `# Build controller binary for Aperture

aperturectl --uri . build controller -c build_config.yaml -o /

Where build_config.yaml can be:
---
build:
  version: 1.0.0
  git_commit_hash: 1234567890
  git_branch: branch1
  ldflags:
    - -some-flag
    - -some-other-flag
  flags:
    - -some-flag
    - -some-other-flag
bundled_extensions: # remote extensions to be bundled
  - go_mod_name: github.com/org/name
    version: v1.0.0
    pkg_name: pkg
extensions: # built-in extensions to be enabled
  - fluxninja
  - sentry
replaces:
  - old: github.com/org/name
    new: github.com/org/name2
enable_core_extensions: false # default is true`,
	RunE: buildRunE("aperture-controller"),
}
