---
sidebar_label: List
hide_title: true
keywords:
  - aperturectl
  - aperturectl_policy_list
---

<!-- markdownlint-disable -->

## aperturectl policy list

List all Aperture Policies from the Aperture Controller

### Synopsis

Use this command to list all the Aperture Policies from the Aperture Controller.

```
aperturectl policy list [flags]
```

### Examples

```
aperturectl policy list
```

### Options

```
  -h, --help   help for list
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
