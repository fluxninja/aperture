---
sidebar_label: Apply
hide_title: true
keywords:
  - aperturectl
  - aperturectl_apply
---

## aperturectl apply

Apply Aperture Policy to the cluster

### Synopsis

Use this command to apply the Aperture Policy to the cluster.

### Examples

```
aperturectl apply --file=policy.yaml

aperturectl apply --dir=policy-dir
```

### Options

```
  -h, --help                 help for apply
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl apply dynamic-config](/reference/aperturectl/apply/dynamic-config/dynamic-config.md) - Apply Aperture DynamicConfig to a Policy
- [aperturectl apply policy](/reference/aperturectl/apply/policy/policy.md) - Apply Aperture Policy to the cluster
