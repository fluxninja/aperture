---
sidebar_label: Delete
hide_title: true
keywords:
  - aperturectl
  - aperturectl_policy_delete
---

<!-- markdownlint-disable -->

## aperturectl policy delete

Delete Aperture Policy from the Aperture Controller

### Synopsis

Use this command to delete the Aperture Policy from the Aperture Controller.

```
aperturectl policy delete POLICY_NAME [flags]
```

### Examples

```
aperturectl policy delete POLICY_NAME
```

### Options

```
  -h, --help   help for delete
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

- [aperturectl policy](/reference/aperturectl/policy/policy.md) - Aperture Policy related commands for the Controller
