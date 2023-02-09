---
sidebar_label: Policy
hide_title: true
keywords:
  - aperturectl
  - aperturectl_apply_policy
---

## aperturectl apply policy

Apply Aperture Policy to the cluster

### Synopsis

Use this command to apply the Aperture Policy to the cluster.

```
aperturectl apply policy [flags]
```

### Examples

```
aperturectl apply policy --file=policies/static-rate-limiting.yaml

aperturectl apply policy --dir=policies
```

### Options

```
      --dir string    Path to directory containing Aperture Policy files
      --file string   Path to Aperture Policy file
  -h, --help          help for policy
```

### SEE ALSO

- [aperturectl apply](/reference/aperturectl/apply/apply.md) - Apply Aperture Policy to the cluster
