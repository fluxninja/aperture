---
sidebar_label: Control-Points
hide_title: true
keywords:
  - aperturectl
  - aperturectl_auto-scale_control-points
---

## aperturectl auto-scale control-points

List AutoScale control points

### Synopsis

List AutoScale control points

```
aperturectl auto-scale control-points [flags]
```

### Examples

```
aperturectl auto-scale control-points --kube
```

### Options

```
  -h, --help   help for control-points
```

### Options inherited from parent commands

```
      --controller string      Address of Aperture controller
      --controller-ns string   Namespace in which the Aperture Controller is running
      --insecure               Allow insecure connection to controller
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
```

### SEE ALSO

- [aperturectl auto-scale](/reference/aperturectl/auto-scale/auto-scale.md) - AutoScale integrations
