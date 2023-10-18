---
sidebar_label: Dynamic-Config
hide_title: true
keywords:
  - aperturectl
  - aperturectl_dynamic-config
---

<!-- markdownlint-disable -->

## aperturectl dynamic-config

DynamicConfig of Aperture Policy related commands for the Controller

### Synopsis

Use this command to manage the DynamicConfig of the Aperture Policies to the Controller.

### Options

```
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for dynamic-config
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl dynamic-config apply](/reference/aperturectl/dynamic-config/apply/apply.md) - Apply Aperture DynamicConfig to a Policy
- [aperturectl dynamic-config delete](/reference/aperturectl/dynamic-config/delete/delete.md) - Delete Aperture DynamicConfig of a Policy.
- [aperturectl dynamic-config get](/reference/aperturectl/dynamic-config/get/get.md) - Get Aperture DynamicConfig for a Policy.
