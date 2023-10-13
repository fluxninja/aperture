---
sidebar_label: Get
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_blueprints_get
---

<!-- markdownlint-disable -->

## aperturectl cloud blueprints get

Cloud Blueprints Get

### Synopsis

Get cloud blueprint.

```
aperturectl cloud blueprints get POLICY_NAME [flags]
```

### Examples

```
aperturectl cloud blueprints get rate-limiting
```

### Options

```
  -h, --help   help for get
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

- [aperturectl cloud blueprints](/reference/aperturectl/cloud/blueprints/blueprints.md) - Cloud Blueprints
