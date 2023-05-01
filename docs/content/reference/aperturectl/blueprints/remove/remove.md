---
sidebar_label: Remove
hide_title: true
keywords:
  - aperturectl
  - aperturectl_blueprints_remove
---

<!-- markdownlint-disable -->

## aperturectl blueprints remove

Remove a Blueprint

### Synopsis

Use this command to remove a pulled Aperture Blueprint from local system.

```
aperturectl blueprints remove [flags]
```

### Examples

```
aperturectl blueprints remove

aperturectl blueprints remove --version latest

aperturectl blueprints remove --all
```

### Options

```
      --all    remove all versions of Aperture Blueprints
  -h, --help   help for remove
```

### Options inherited from parent commands

```
      --skip-pull        Skip pulling the blueprints update.
      --uri string       URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@latest. This field should not be provided when the Version is provided.
      --version string   Version of official Aperture Blueprints, e.g. latest. This field should not be provided when the URI is provided (default "latest")
```

### SEE ALSO

- [aperturectl blueprints](/reference/aperturectl/blueprints/blueprints.md) - Aperture Blueprints
