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

```
aperturectl apply [flags]
```

### Examples

```
aperturectl apply --file=policy.yaml

aperturectl apply --dir=policy-dir
```

### Options

```
      --dir string                   Path to directory containing Aperture Policy files
      --dynamic-config-file string   Path to the dynamic config file
      --file string                  Path to Aperture Policy file
  -h, --help                         help for apply
      --kube-config string           Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
