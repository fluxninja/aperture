---
sidebar_label: Result-Cache
hide_title: true
keywords:
  - aperturectl
  - aperturectl_result-cache
---

<!-- markdownlint-disable -->

## aperturectl result-cache

Result Cache related commands

### Synopsis

Use this command to interact with Aperture's Result Cache.

### Options

```
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for result-cache
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperture-cli/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl result-cache delete](/reference/aperture-cli/aperturectl/result-cache/delete/delete.md) - Delete a result cache entry
- [aperturectl result-cache get](/reference/aperture-cli/aperturectl/result-cache/get/get.md) - Get a result cache entry
- [aperturectl result-cache set](/reference/aperture-cli/aperturectl/result-cache/set/set.md) - Set a result cache entry
