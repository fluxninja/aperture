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

### Options

```
  -c, --config string       path to the build configuration file
  -h, --help                help for controller
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
