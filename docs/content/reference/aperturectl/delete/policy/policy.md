---
sidebar_label: Policy
hide_title: true
keywords:
  - aperturectl
  - aperturectl_delete_policy
---

<!-- markdownlint-disable -->

## aperturectl delete policy

Delete Aperture Policy from the cluster

### Synopsis

Use this command to delete the Aperture Policy from the cluster.

```
aperturectl delete policy [flags]
```

### Examples

```
aperturectl delete policy --policy=rate-limiting
```

### Options

```
  -h, --help   help for policy
```

### Options inherited from parent commands

```
      --api-key string         FluxNinja Cloud API Key to be used when using Cloud Controller
      --config string          Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --policy string          Name of the Policy to delete
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl delete](/reference/aperturectl/delete/delete.md) - Delete Aperture Policies
