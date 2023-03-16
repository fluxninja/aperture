---
sidebar_label: Agent
hide_title: true
keywords:
  - aperturectl
  - aperturectl_build_agent
---

## aperturectl build agent

Build agent binary for Aperture

### Synopsis

Build agent binary for Aperture

```
aperturectl build agent [flags]
```

### Examples

```
# Build agent binary for Aperture

aperturectl --uri . build agent -c build_config.yaml -o /

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
enable_core_extensions: false # default is true
```

### Options

```
  -c, --config string       path to the build configuration file
  -h, --help                help for agent
  -o, --output-dir string   path to the output directory
```

### Options inherited from parent commands

```
      --skip-pull        Skip pulling the repository update.
      --uri string       URI of Aperture repository, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture@latest. This field should not be provided when the Version is provided.
      --version string   Version of Aperture, e.g. latest. This field should not be provided when the URI is provided (default "latest")
```

### SEE ALSO

- [aperturectl build](/reference/aperturectl/build/build.md) - Builds the agent and controller binaries
