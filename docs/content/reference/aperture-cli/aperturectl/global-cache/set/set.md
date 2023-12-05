---
sidebar_label: Set
hide_title: true
keywords:
  - aperturectl
  - aperturectl_global-cache_set
---

<!-- markdownlint-disable -->

## aperturectl global-cache set

Set a global cache entry

### Synopsis

Set a global cache entry

```
aperturectl global-cache set [flags]
```

### Options

```
  -a, --agent-group string   Agent group
  -h, --help                 help for set
  -k, --key string           Key
  -t, --ttl int              TTL in seconds (default 600)
      --value string         Value
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
