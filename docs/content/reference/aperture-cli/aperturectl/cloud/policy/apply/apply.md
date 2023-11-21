---
sidebar_label: Apply
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_policy_apply
---

<!-- markdownlint-disable -->

## aperturectl cloud policy apply

Apply Aperture Policy to the Aperture Cloud Controller

### Synopsis

Use this command to apply the Aperture Policy to the Aperture Cloud Controller.

```
aperturectl cloud policy apply [flags]
```

### Examples

```
aperturectl cloud policy apply --file=policies/rate-limiting.yaml

aperturectl cloud policy apply --dir=policies
```

### Options

```
      --dir string    Path to directory containing Aperture Policy files
      --file string   Path to Aperture Policy file
  -f, --force         Force apply policy even if it already exists
  -h, --help          help for apply
  -s, --select-all    Apply all policies in the directory
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
