---
sidebar_label: Get
hide_title: true
keywords:
  - aperturectl
  - aperturectl_policy_get
---

<!-- markdownlint-disable -->

## aperturectl policy get

Get Aperture Policy from the Aperture Controller

### Synopsis

Use this command to get the Aperture Policy from the Aperture Controller.

```
aperturectl policy get POLICY_NAME [flags]
```

### Examples

```
aperturectl policy get POLICY_NAME
```

### Options

```
  -h, --help   help for get
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

- [aperturectl policy](/reference/aperture-cli/aperturectl/policy/policy.md) - Aperture Policy related commands for the Controller
