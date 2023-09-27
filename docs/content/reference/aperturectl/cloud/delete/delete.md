---
sidebar_label: Delete
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_delete
---

<!-- markdownlint-disable -->

## aperturectl cloud delete

Delete Aperture Policies from Aperture Cloud

### Synopsis

Use this command to delete the Aperture Policies from Aperture Cloud.

### Options

```
      --api-key string        Aperture Cloud User API Key to be used when using Cloud Controller
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
  -h, --help                  help for delete
      --insecure              Allow connection to controller running without TLS
      --policy string         Name of the Policy to delete
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
- [aperturectl cloud delete policy](/reference/aperturectl/cloud/delete/policy/policy.md) - Delete Aperture Policy from the Aperture Cloud Controller
