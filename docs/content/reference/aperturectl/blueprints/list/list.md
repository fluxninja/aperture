---
sidebar_label: List
hide_title: true
keywords:
  - aperturectl
  - aperturectl_blueprints_list
---

## aperturectl blueprints list

List Aperture Blueprints

### Synopsis

Use this command to list the Aperture Blueprints which are already pulled and available in local system.

```
aperturectl blueprints list [flags]
```

### Examples

```
aperturectl blueprints list

aperturectl blueprints list --version v0.22.0

aperturectl blueprints list --all
```

### Options

```
      --all    show the entire cache of Aperture Blueprints
  -h, --help   help for list
```

### Options inherited from parent commands

```
      --skip-pull        Skip pulling the blueprints update.
      --uri string       URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@latest. This field should not be provided when the Version is provided.
      --version string   Version of official Aperture Blueprints, e.g. latest. This field should not be provided when the URI is provided (default "latest")
```

### SEE ALSO

- [aperturectl blueprints](/reference/aperturectl/blueprints/blueprints.md) - Aperture Blueprints
