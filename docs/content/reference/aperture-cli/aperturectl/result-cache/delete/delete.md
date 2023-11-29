---
sidebar_label: Delete
hide_title: true
keywords:
  - aperturectl
  - aperturectl_result-cache_delete
---

<!-- markdownlint-disable -->

## aperturectl result-cache delete

Delete a result cache entry

### Synopsis

Delete a result cache entry

```
aperturectl result-cache delete [flags]
```

### Options

```
  -a, --agent-group string     Agent group
  -c, --control-point string   Control point
  -h, --help                   help for delete
  -k, --key string             Key
```

### Options inherited from parent commands

```
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl result-cache](/reference/aperture-cli/aperturectl/result-cache/result-cache.md) - Result Cache related commands
