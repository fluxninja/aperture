---
sidebar_label: Policy
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_delete_policy
---

<!-- markdownlint-disable -->

## aperturectl cloud delete policy

Delete Aperture Policy from the Aperture Cloud Controller

### Synopsis

Use this command to delete the Aperture Policy from the Aperture Cloud Controller.

```
aperturectl cloud delete policy [flags]
```

### Examples

```
aperturectl cloud delete policy --policy=rate-limiting --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key API_KEY
```

### Options

```
  -h, --help   help for policy
```

### Options inherited from parent commands

```
      --api-key string      Aperture Cloud API Key to be used when using Cloud Controller
      --config string       Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string   Address of Aperture Cloud Controller
      --insecure            Allow connection to controller running without TLS
      --policy string       Name of the Policy to delete
      --skip-verify         Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud delete](/reference/aperturectl/cloud/delete/delete.md) - Delete Aperture Policies from Aperture Cloud
