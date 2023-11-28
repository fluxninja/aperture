---
sidebar_label: Set
hide_title: true
keywords:
  - aperturectl
  - aperturectl_state-cache_set
---

<!-- markdownlint-disable -->

## aperturectl state-cache set

Set a state cache entry

### Synopsis

Set a state cache entry

```
aperturectl state-cache set [flags]
```

### Options

```
  -a, --agent-group string     Agent group
  -c, --control-point string   Control point
  -h, --help                   help for set
  -k, --key string             Key
  -t, --ttl int                TTL in milliseconds (default 600000)
      --value string           Value
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

- [aperturectl state-cache](/reference/aperture-cli/aperturectl/state-cache/state-cache.md) - State Cache related commands
