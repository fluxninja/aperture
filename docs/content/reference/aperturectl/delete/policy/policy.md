---
sidebar_label: Policy
hide_title: true
keywords:
  - aperturectl
  - aperturectl_delete_policy
---

## aperturectl delete policy

Delete Aperture Policy from the cluster

### Synopsis

Use this command to delete the Aperture Policy from the cluster.

```
aperturectl delete policy [flags]
```

### Examples

```
aperturectl delete policy --policy=static-rate-limiting
```

### Options

```
  -h, --help   help for policy
```

### Options inherited from parent commands

```
      --controller string      Address of Aperture controller
      --controller-ns string   Namespace in which the Aperture Controller is running
      --insecure               Allow insecure connection to controller
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --policy string          Name of the Policy to delete
```

### SEE ALSO

- [aperturectl delete](/reference/aperturectl/delete/delete.md) - Delete Aperture Policies
