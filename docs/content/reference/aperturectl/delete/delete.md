---
sidebar_label: Delete
hide_title: true
keywords:
  - aperturectl
  - aperturectl_delete
---

## aperturectl delete

Delete Aperture Policies

### Synopsis

Use this command to delete the Aperture Policies.

### Options

```
      --controller string      Address of Aperture controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for delete
      --insecure               Allow insecure connection to controller
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --policy string          Name of the Policy to delete
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl delete policy](/reference/aperturectl/delete/policy/policy.md) - Delete Aperture Policy from the cluster
