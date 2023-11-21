---
sidebar_label: Archive
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_blueprints_archive
---

<!-- markdownlint-disable -->

## aperturectl cloud blueprints archive

Cloud Blueprints Archive for the given Policy Name

### Synopsis

Archive cloud blueprint.

```
aperturectl cloud blueprints archive POLICY_NAME [flags]
```

### Examples

```
aperturectl cloud blueprints archive rate-limiting
```

### Options

```
  -h, --help   help for archive
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
