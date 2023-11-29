---
sidebar_label: Global-Cache
hide_title: true
keywords:
  - aperturectl
  - aperturectl_global-cache
---

<!-- markdownlint-disable -->

## aperturectl global-cache

Global Cache related commands

### Synopsis

Use this command to interact with Aperture's Global Cache.

### Options

```
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for global-cache
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperture-cli/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl global-cache delete](/reference/aperture-cli/aperturectl/global-cache/delete/delete.md) - Delete a global cache entry
- [aperturectl global-cache get](/reference/aperture-cli/aperturectl/global-cache/get/get.md) - Get a global cache entry
- [aperturectl global-cache set](/reference/aperture-cli/aperturectl/global-cache/set/set.md) - Set a global cache entry
