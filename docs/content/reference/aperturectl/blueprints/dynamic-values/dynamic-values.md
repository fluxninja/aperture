---
sidebar_label: Dynamic-Values
hide_title: true
keywords:
  - aperturectl
  - aperturectl_blueprints_dynamic-values
---

<!-- markdownlint-disable -->

## aperturectl blueprints dynamic-values

Create dynamic values file for a given Aperture Blueprint

### Synopsis

Provides a dynamic values file for a given Aperture Blueprint that can be then used to generate policies after customization

```
aperturectl blueprints dynamic-values [flags]
```

### Examples

```
aperturectl blueprints dynamic-values --name=policies/rate-limiting --output-file=values.yaml
```

### Options

```
  -h, --help                 help for dynamic-values
      --name string          Name of the Aperture Blueprint to provide values file for
      --no-yaml-modeline     Do not add YAML language server modeline to generated YAML files
      --output-file string   Path to the output values file
      --overwrite            Overwrite existing values file
```

### Options inherited from parent commands

```
      --skip-pull        Skip pulling the blueprints update.
      --uri string       URI of Custom Blueprints, could be a local path or a remote git repository, e.g. github.com/fluxninja/aperture/blueprints@latest. This field should not be provided when the Version is provided.
      --version string   Version of official Aperture Blueprints, e.g. latest. This field should not be provided when the URI is provided (default "latest")
```

### SEE ALSO

- [aperturectl blueprints](/reference/aperturectl/blueprints/blueprints.md) - Aperture Blueprints
