---
sidebar_label: Apply
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_dynamic-config_apply
---

<!-- markdownlint-disable -->

## aperturectl cloud dynamic-config apply

Apply Aperture DynamicConfig to a Policy

### Synopsis

Use this command to apply the Aperture DynamicConfig to a Policy.

```
aperturectl cloud dynamic-config apply [flags]
```

### Examples

```
aperturectl cloud dynamic-config apply --policy=rate-limiting --file=dynamic-config.yaml
```

### Options

```
      --file string     Path to the dynamic config file
  -h, --help            help for apply
      --policy string   Name of the Policy to apply the DynamicConfig to
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

- [aperturectl cloud dynamic-config](/reference/aperture-cli/aperturectl/cloud/dynamic-config/dynamic-config.md) - DynamicConfig of Aperture Policy related commands for the Cloud Controller
