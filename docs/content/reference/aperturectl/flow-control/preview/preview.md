---
sidebar_label: Preview
hide_title: true
keywords:
  - aperturectl
  - aperturectl_flow-control_preview
---

<!-- markdownlint-disable -->

## aperturectl flow-control preview

Preview control points

### Synopsis

Preview samples of flow labels or HTTP requests on control points

```
aperturectl flow-control preview [--http] SERVICE CONTROL_POINT [flags]
```

### Options

```
      --agent-group string   Agent group (default "default")
  -h, --help                 help for preview
      --http                 Preview HTTP requests instead of flow labels
      --samples int          Number of samples to collect (default 10)
```

### Options inherited from parent commands

```
      --controller string      Address of Aperture controller
      --controller-ns string   Namespace in which the Aperture Controller is running
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl flow-control](/reference/aperturectl/flow-control/flow-control.md) - Flow Control integrations
