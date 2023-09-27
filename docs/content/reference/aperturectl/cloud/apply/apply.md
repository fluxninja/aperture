---
sidebar_label: Apply
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_apply
---

<!-- markdownlint-disable -->

## aperturectl cloud apply

Apply Aperture Policies to the Cloud Controller

### Synopsis

Use this command to apply the Aperture Policies to the Cloud Controller.

### Options

```
      --api-key string        Aperture Cloud User API Key to be used when using Cloud Controller
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
  -h, --help                  help for apply
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
- [aperturectl cloud apply policy](/reference/aperturectl/cloud/apply/policy/policy.md) - Apply Aperture Policy to the Aperture Cloud Controller
