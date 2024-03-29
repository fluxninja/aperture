---
sidebar_label: Policy
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_policy
---

<!-- markdownlint-disable -->

## aperturectl cloud policy

Aperture Policy related commands for the Cloud Controller

### Synopsis

Use this command to manage the Aperture Policies to the Cloud Controller.

### Options

```
      --access-token string   User Access Token to be used while connecting to Aperture Cloud
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
  -h, --help                  help for policy
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperture-cli/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
- [aperturectl cloud policy apply](/reference/aperture-cli/aperturectl/cloud/policy/apply/apply.md) - Apply Aperture Policy to the Aperture Cloud Controller
- [aperturectl cloud policy archive](/reference/aperture-cli/aperturectl/cloud/policy/archive/archive.md) - Archive Aperture Policy from the Aperture Cloud Controller
- [aperturectl cloud policy delete](/reference/aperture-cli/aperturectl/cloud/policy/delete/delete.md) - Delete Aperture Policy from the Aperture Cloud Controller
- [aperturectl cloud policy get](/reference/aperture-cli/aperturectl/cloud/policy/get/get.md) - Get Aperture Policy from the Aperture Cloud Controller
- [aperturectl cloud policy list](/reference/aperture-cli/aperturectl/cloud/policy/list/list.md) - List all Aperture Policies from the Aperture Cloud Controller
