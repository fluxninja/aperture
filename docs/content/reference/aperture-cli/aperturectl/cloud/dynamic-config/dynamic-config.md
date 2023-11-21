---
sidebar_label: Dynamic-Config
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_dynamic-config
---

<!-- markdownlint-disable -->

## aperturectl cloud dynamic-config

DynamicConfig of Aperture Policy related commands for the Cloud Controller

### Synopsis

Use this command to manage the DynamicConfig of the Aperture Policies to the Cloud Controller.

### Options

```
      --access-token string   User Access Token to be used while connecting to Aperture Cloud
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
  -h, --help                  help for dynamic-config
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperture-cli/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
- [aperturectl cloud dynamic-config apply](/reference/aperture-cli/aperturectl/cloud/dynamic-config/apply/apply.md) - Apply Aperture DynamicConfig to a Policy
- [aperturectl cloud dynamic-config delete](/reference/aperture-cli/aperturectl/cloud/dynamic-config/delete/delete.md) - Delete Aperture DynamicConfig of a Policy.
- [aperturectl cloud dynamic-config get](/reference/aperture-cli/aperturectl/cloud/dynamic-config/get/get.md) - Get Aperture DynamicConfig for a Policy.
