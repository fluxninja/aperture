---
sidebar_label: Apply
hide_title: true
keywords:
  - aperturectl
  - aperturectl_apply
---

<!-- markdownlint-disable -->

## aperturectl apply

Apply Aperture Policies

### Synopsis

Use this command to apply the Aperture Policies.

### Options

```
      --api-key string         FluxNinja API Key to be used when using Cloud Controller
      --config string          Path to the Aperture config file
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for apply
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl apply dynamic-config](/reference/aperturectl/apply/dynamic-config/dynamic-config.md) - Apply Aperture DynamicConfig to a Policy
- [aperturectl apply policy](/reference/aperturectl/apply/policy/policy.md) - Apply Aperture Policy to the cluster
