---
sidebar_label: Control-Points
hide_title: true
keywords:
  - aperturectl
  - aperturectl_auto-scale_control-points
---

<!-- markdownlint-disable -->

## aperturectl auto-scale control-points

List AutoScale control points

### Synopsis

List AutoScale control points

```
aperturectl auto-scale control-points [flags]
```

### Examples

```
aperturectl auto-scale control-points
```

### Options

```
  -h, --help   help for control-points
```

### Options inherited from parent commands

```
      --api-key string         FluxNinja API Key to be used when using Cloud Controller
      --config string          Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl auto-scale](/reference/aperturectl/auto-scale/auto-scale.md) - AutoScale integrations
