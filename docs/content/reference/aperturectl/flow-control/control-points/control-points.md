---
sidebar_label: Control-Points
hide_title: true
keywords:
  - aperturectl
  - aperturectl_flow-control_control-points
---

## aperturectl flow-control control-points

List Flow Control control points

### Synopsis

List Flow Control control points

```
aperturectl flow-control control-points [flags]
```

### Examples

```
aperturectl flow-control control-points --kube
```

### Options

```
  -h, --help   help for control-points
```

### Options inherited from parent commands

```
      --controller string    Address of Aperture controller
      --insecure             Allow insecure connection to controller
      --kube                 Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
```

### SEE ALSO

- [aperturectl flow-control](/reference/aperturectl/flow-control/flow-control.md) -
  Flow Control integrations
