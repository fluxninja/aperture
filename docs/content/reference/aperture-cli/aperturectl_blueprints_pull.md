---
title: Aperturectl Blueprints Pull
description: Aperturectl Blueprints Pull
keywords:
  - aperturectl
  - aperturectl_blueprints_pull
sidebar_position: 11
---

## aperturectl blueprints pull

Pull Aperture Blueprints

### Synopsis

Use this command to pull the Aperture Blueprints in local system to use for generating Aperture Policies and Grafana Dashboards.

```
aperturectl blueprints pull [flags]
```

### Examples

```
aperturectl blueprints pull

aperturectl blueprints pull --version v0.22.0
```

### Options

```
  -h, --help   help for pull
```

### Options inherited from parent commands

```
      --uri string       URI of Aperture Blueprints, could be a local path or a remote git repository (default "github.com/fluxninja/aperture/blueprints")
      --version string   Version of Aperture Blueprints (default "main")
```

### SEE ALSO

- [aperturectl blueprints](aperturectl_blueprints.md) - Aperture Blueprints
