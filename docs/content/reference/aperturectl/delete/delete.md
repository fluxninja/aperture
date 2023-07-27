---
sidebar_label: Delete
hide_title: true
keywords:
  - aperturectl
  - aperturectl_delete
---

<!-- markdownlint-disable -->

## aperturectl delete

Delete Aperture Policies

### Synopsis

Use this command to delete the Aperture Policies.

### Options

```
      --api-key string         FluxNinja API Key to be used when using Cloud Controller
      --config string          Path to the Aperture config file
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for delete
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --policy string          Name of the Policy to delete
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl delete policy](/reference/aperturectl/delete/policy/policy.md) - Delete Aperture Policy from the cluster
