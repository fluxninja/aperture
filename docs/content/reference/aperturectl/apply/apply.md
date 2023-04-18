---
sidebar_label: Apply
hide_title: true
keywords:
  - aperturectl
  - aperturectl_apply
---

## aperturectl apply

Apply Aperture Policies

### Synopsis

Use this command to apply the Aperture Policies.

### Options

```
      --controller string      Address of Aperture controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for apply
      --insecure               Allow insecure connection to controller
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl apply dynamic-config](/reference/aperturectl/apply/dynamic-config/dynamic-config.md) - Apply Aperture DynamicConfig to a Policy
- [aperturectl apply policy](/reference/aperturectl/apply/policy/policy.md) - Apply Aperture Policy to the cluster
