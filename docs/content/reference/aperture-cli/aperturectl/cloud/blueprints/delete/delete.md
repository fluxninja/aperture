---
sidebar_label: Delete
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_blueprints_delete
---

<!-- markdownlint-disable -->

## aperturectl cloud blueprints delete

Cloud Blueprints Delete for the given Policy Name

### Synopsis

Delete cloud blueprint.

```
aperturectl cloud blueprints delete POLICY_NAME [flags]
```

### Examples

```
aperturectl cloud blueprints delete rate-limiting
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
      --api-key string        Aperture Cloud User API Key to be used when using Cloud Controller
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud blueprints](/reference/aperture-cli/aperturectl/cloud/blueprints/blueprints.md) - Cloud Blueprints
