---
sidebar_label: Apply
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_blueprints_apply
---

<!-- markdownlint-disable -->

## aperturectl cloud blueprints apply

Cloud Blueprints Apply

### Synopsis

Apply cloud blueprint.

```
aperturectl cloud blueprints apply [flags]
```

### Examples

```
aperturectl cloud blueprints apply --value-file=values.yaml
```

### Options

```
  -h, --help                 help for apply
      --values-dir string    Path to directory containing blueprint values files
      --values-file string   Values file to use for blueprint
```

### Options inherited from parent commands

```
      --access-token string   User Access Token to be used while connecting to Aperture Cloud
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud blueprints](/reference/aperture-cli/aperturectl/cloud/blueprints/blueprints.md) - Cloud Blueprints
