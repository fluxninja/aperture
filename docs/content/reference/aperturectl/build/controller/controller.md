---
sidebar_label: Controller
hide_title: true
keywords:
  - aperturectl
  - aperturectl_build_controller
---

## aperturectl build controller

Build controller binary for Aperture

### Synopsis

Build controller binary for Aperture

```
aperturectl build controller [flags]
```

### Examples

```
# Build controller binary for Aperture

aperturectl --uri . build controller -c build-config.yaml -o /

Where build-config.yaml can be:
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
enable_core_extensions: false # default is true
```

### Options

```
  -c, --config string       path to the build configuration file (default: build-config.yaml in the main package directory)
  -h, --help                help for controller
  -o, --output-dir string   path to the output directory (default: current directory)
```

### Options inherited from parent commands

```
      --skip-pull        Skip pulling the repository update.
      --uri string       URI of Aperture repository, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture@latest. This field should not be provided when the Version is provided.
      --version string   Version of Aperture, e.g. latest. This field should not be provided when the URI is provided (default "latest")
```

### SEE ALSO

- [aperturectl build](/reference/aperturectl/build/build.md) - Builds the agent and controller binaries
