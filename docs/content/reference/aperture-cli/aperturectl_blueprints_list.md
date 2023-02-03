---
title: Aperturectl Blueprints List
description: Aperturectl Blueprints List
keywords:
  - aperturectl
  - aperturectl_blueprints_list
sidebar_position: 12
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
      --all    list all versions of aperture blueprints
  -h, --help   help for list
```

### Options inherited from parent commands

```
      --uri string       URI of Aperture Blueprints, could be a local path or a remote git repository (default "github.com/fluxninja/aperture/blueprints")
      --version string   Version of Aperture Blueprints (default "latest")
```

### SEE ALSO

- [aperturectl blueprints](aperturectl_blueprints.md) - Aperture Blueprints
