---
sidebar_label: Pull
hide_title: true
keywords:
  - aperturectl
  - aperturectl_blueprints_pull
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
      --skip-pull        Skip pulling the latest blueprints.
      --uri string       URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@main. This field should not be provided when the Version is provided.
      --version string   Version of official Aperture Blueprints, e.g. latest. This field should not be provided when the URI is provided (default "latest")
```

### SEE ALSO

- [aperturectl blueprints](/reference/aperturectl/blueprints/blueprints.md) - Aperture Blueprints
