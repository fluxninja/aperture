---
title: Aperturectl Blueprints Remove
description: Aperturectl Blueprints Remove
keywords:
  - aperturectl
  - aperturectl_blueprints_remove
sidebar_position: 10
---

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

aperturectl blueprints remove --version v0.22.0

aperturectl blueprints remove --all
```

### Options

```
      --all    remove all versions of Aperture Blueprints
  -h, --help   help for remove
```

### Options inherited from parent commands

```
      --uri string       URI of Aperture Blueprints, could be a local path or a remote git repository (default "github.com/fluxninja/aperture/blueprints")
      --version string   Version of Aperture Blueprints (default "latest")
```

### SEE ALSO

- [aperturectl blueprints](aperturectl_blueprints.md) - Aperture Blueprints
