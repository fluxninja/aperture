---
sidebar_label: List
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_policy_list
---

<!-- markdownlint-disable -->

## aperturectl cloud policy list

List all Aperture Policies from the Aperture Cloud Controller

### Synopsis

Use this command to list all the Aperture Policies from the Aperture Cloud Controller.

```
aperturectl cloud policy list [flags]
```

### Examples

```
aperturectl cloud policy list
```

### Options

```
  -h, --help   help for list
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

- [aperturectl cloud policy](/reference/aperture-cli/aperturectl/cloud/policy/policy.md) - Aperture Policy related commands for the Cloud Controller
