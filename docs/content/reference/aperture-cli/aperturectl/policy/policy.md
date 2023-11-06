---
sidebar_label: Policy
hide_title: true
keywords:
  - aperturectl
  - aperturectl_policy
---

<!-- markdownlint-disable -->

## aperturectl policy

Aperture Policy related commands for the Controller

### Synopsis

Use this command to manage the Aperture Policies to the Controller.

### Options

```
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for policy
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperture-cli/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl policy apply](/reference/aperture-cli/aperturectl/policy/apply/apply.md) - Apply Aperture Policy to the cluster
- [aperturectl policy delete](/reference/aperture-cli/aperturectl/policy/delete/delete.md) - Delete Aperture Policy from the Aperture Controller
- [aperturectl policy get](/reference/aperture-cli/aperturectl/policy/get/get.md) - Get Aperture Policy from the Aperture Controller
- [aperturectl policy list](/reference/aperture-cli/aperturectl/policy/list/list.md) - List all Aperture Policies from the Aperture Controller
