---
title: Aperturectl Apply
description: Aperturectl Apply
keywords:
  - aperturectl
  - aperturectl_apply
sidebar_position: 14
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
      --dir string           Path to directory containing Aperture Policy files
      --file string          Path to Aperture Policy file
  -h, --help                 help for apply
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
```

### SEE ALSO

- [aperturectl](aperturectl.md) - aperturectl - CLI tool to interact with Aperture
