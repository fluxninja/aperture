---
sidebar_label: Auto-Scale
hide_title: true
keywords:
  - aperturectl
  - aperturectl_auto-scale
---

<!-- markdownlint-disable -->

## aperturectl auto-scale

AutoScale integrations

### Synopsis

Use this command to query information about active AutoScale integrations

### Options

```
      --api-key string         FluxNinja ARC API Key to be used when using Cloud Controller
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for auto-scale
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl auto-scale control-points](/reference/aperturectl/auto-scale/control-points/control-points.md) - List AutoScale control points
