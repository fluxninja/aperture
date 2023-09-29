---
sidebar_label: Policy
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_apply_policy
---

<!-- markdownlint-disable -->

## aperturectl cloud apply policy

Apply Aperture Policy to the Aperture Cloud Controller

### Synopsis

Use this command to apply the Aperture Policy to the Aperture Cloud Controller.

```
aperturectl cloud apply policy [flags]
```

### Examples

```
aperturectl cloud apply policy --file=policies/rate-limiting.yaml --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key PERSONAL_API_KEY

aperturectl cloud apply policy --dir=policies --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key PERSONAL_API_KEY
```

### Options

```
      --dir string    Path to directory containing Aperture Policy files
      --file string   Path to Aperture Policy file
  -f, --force         Force apply policy even if it already exists
  -h, --help          help for policy
  -s, --select-all    Apply all policies in the directory
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

- [aperturectl cloud apply](/reference/aperturectl/cloud/apply/apply.md) - Apply Aperture Policies to the Cloud Controller
