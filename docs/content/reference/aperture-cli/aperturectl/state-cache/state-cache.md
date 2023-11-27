---
sidebar_label: State-Cache
hide_title: true
keywords:
  - aperturectl
  - aperturectl_state-cache
---

<!-- markdownlint-disable -->

## aperturectl state-cache

State Cache related commands

### Synopsis

Use this command to interact with Aperture State Cache.

### Options

```
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for state-cache
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperture-cli/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl state-cache delete](/reference/aperture-cli/aperturectl/state-cache/delete/delete.md) - Delete a state cache entry
- [aperturectl state-cache get](/reference/aperture-cli/aperturectl/state-cache/get/get.md) - Get a state cache entry
- [aperturectl state-cache set](/reference/aperture-cli/aperturectl/state-cache/set/set.md) - Set a state cache entry
