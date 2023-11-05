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

- [aperturectl](/reference/aperture-cli/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl blueprints dynamic-values](/reference/aperture-cli/aperturectl/blueprints/dynamic-values/dynamic-values.md) - Create dynamic values file for a given Aperture Blueprint
- [aperturectl blueprints generate](/reference/aperture-cli/aperturectl/blueprints/generate/generate.md) - Generate Aperture Policy related resources from Aperture Blueprints
- [aperturectl blueprints list](/reference/aperture-cli/aperturectl/blueprints/list/list.md) - List Aperture Blueprints
- [aperturectl blueprints pull](/reference/aperture-cli/aperturectl/blueprints/pull/pull.md) - Pull Aperture Blueprints
- [aperturectl blueprints remove](/reference/aperture-cli/aperturectl/blueprints/remove/remove.md) - Remove a Blueprint
- [aperturectl blueprints values](/reference/aperture-cli/aperturectl/blueprints/values/values.md) - Create values file for a given Aperture Blueprint
