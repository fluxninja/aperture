---
sidebar_label: Policy
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_policy
---

<!-- markdownlint-disable -->

## aperturectl cloud policy

Apply Aperture Policy to the Aperture Cloud Controller

### Synopsis

Use this command to apply the Aperture Policy to the Aperture Cloud Controller.

```
aperturectl cloud policy [flags]
```

### Examples

```
aperturectl cloud apply policy --file=policies/rate-limiting.yaml --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key API_KEY

aperturectl cloud apply policy --dir=policies --controller ORGANIZATION_NAME.app.fluxninja.com:443 --api-key API_KEY
```

### Options

```
      --dir string    Path to directory containing Aperture Policy files
      --file string   Path to Aperture Policy file
  -f, --force         Force apply policy even if it already exists
  -h, --help          help for policy
  -s, --select-all    Apply all policies in the directory
```

### SEE ALSO

- [aperturectl cloud](/reference/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
