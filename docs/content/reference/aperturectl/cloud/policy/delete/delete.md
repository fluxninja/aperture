---
sidebar_label: Delete
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_policy_delete
---

<!-- markdownlint-disable -->

## aperturectl cloud policy delete

Delete Aperture Policy from the Aperture Cloud Controller

### Synopsis

Use this command to delete the Aperture Policy from the Aperture Cloud Controller.

```
aperturectl cloud policy delete POLICY_NAME [flags]
```

### Examples

```
aperturectl cloud policy delete POLICY_NAME
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

- [aperturectl cloud policy](/reference/aperturectl/cloud/policy/policy.md) - Aperture Policy related commands for the Cloud Controller
