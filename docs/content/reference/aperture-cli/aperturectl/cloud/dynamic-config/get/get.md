---
sidebar_label: Get
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_dynamic-config_get
---

<!-- markdownlint-disable -->

## aperturectl cloud dynamic-config get

Get Aperture DynamicConfig for a Policy.

### Synopsis

Use this command to get the Aperture DynamicConfig of a Policy.

```
aperturectl cloud dynamic-config get POLICY_NAME [flags]
```

### Examples

```
aperture cloud dynamic-config get rate-limiting
```

### Options

```
  -h, --help   help for get
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
