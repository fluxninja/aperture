---
sidebar_label: Get
hide_title: true
keywords:
  - aperturectl
  - aperturectl_global-cache_get
---

<!-- markdownlint-disable -->

## aperturectl global-cache get

Get a global cache entry

### Synopsis

Get a global cache entry

```
aperturectl global-cache get [flags]
```

### Options

```
  -a, --agent-group string   Agent group
  -h, --help                 help for get
  -k, --key string           Key
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

- [aperturectl global-cache](/reference/aperture-cli/aperturectl/global-cache/global-cache.md) - Global Cache related commands
