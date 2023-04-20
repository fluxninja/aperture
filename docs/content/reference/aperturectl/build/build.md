---
sidebar_label: Build
hide_title: true
keywords:
  - aperturectl
  - aperturectl_build
---

## aperturectl build

Builds the agent and controller binaries

### Synopsis

Builds the agent and controller binaries

### Options

```
  -h, --help             help for build
      --skip-pull        Skip pulling the repository update.
      --uri string       URI of Aperture repository, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture@latest. This field should not be provided when the Version is provided.
      --version string   Version of Aperture, e.g. latest. This field should not be provided when the URI is provided (default "latest")
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl build agent](/reference/aperturectl/build/agent/agent.md) - Build agent binary for Aperture
- [aperturectl build controller](/reference/aperturectl/build/controller/controller.md) - Build controller binary for Aperture
