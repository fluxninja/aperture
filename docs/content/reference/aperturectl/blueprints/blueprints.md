---
sidebar_label: Blueprints
hide_title: true
keywords:
  - aperturectl
  - aperturectl_blueprints
---

<!-- markdownlint-disable -->

## aperturectl blueprints

Aperture Blueprints

### Synopsis

Use this command to pull, list, remove and generate Aperture Policy resources using the Aperture Blueprints.

### Options

```
  -h, --help             help for blueprints
      --skip-pull        Skip pulling the blueprints update.
      --uri string       URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@latest. This field should not be provided when the Version is provided.
      --version string   Version of official Aperture Blueprints, e.g. latest. This field should not be provided when the URI is provided
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl blueprints dynamic-values](/reference/aperturectl/blueprints/dynamic-values/dynamic-values.md) - Create dynamic values file for a given Aperture Blueprint
- [aperturectl blueprints generate](/reference/aperturectl/blueprints/generate/generate.md) - Generate Aperture Policy related resources from Aperture Blueprints
- [aperturectl blueprints list](/reference/aperturectl/blueprints/list/list.md) - List Aperture Blueprints
- [aperturectl blueprints pull](/reference/aperturectl/blueprints/pull/pull.md) - Pull Aperture Blueprints
- [aperturectl blueprints remove](/reference/aperturectl/blueprints/remove/remove.md) - Remove a Blueprint
- [aperturectl blueprints values](/reference/aperturectl/blueprints/values/values.md) - Create values file for a given Aperture Blueprint
