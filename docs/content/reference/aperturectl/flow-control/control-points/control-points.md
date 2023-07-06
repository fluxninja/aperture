---
sidebar_label: Control-Points
hide_title: true
keywords:
  - aperturectl
  - aperturectl_flow-control_control-points
---

<!-- markdownlint-disable -->

## aperturectl flow-control control-points

List Flow Control control points

### Synopsis

List Flow Control control points

```
aperturectl flow-control control-points [flags]
```

### Examples

```
aperturectl flow-control control-points
```

### Options

```
  -h, --help   help for control-points
```

### Options inherited from parent commands

```
      --api-key string         FluxNinja ARC API Key to be used when using Cloud Controller
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl flow-control](/reference/aperturectl/flow-control/flow-control.md) - Flow Control integrations
