---
sidebar_label: Dynamic-Config
hide_title: true
keywords:
  - aperturectl
  - aperturectl_apply_dynamic-config
---

## aperturectl apply dynamic-config

Apply Aperture DynamicConfig to a Policy

### Synopsis

Use this command to apply the Aperture DynamicConfig to a Policy.

```
aperturectl apply dynamic-config [flags]
```

### Examples

```
aperturectl apply dynamic-config --policy=static-rate-limiting --file=dynamic-config.yaml
```

### Options

```
      --file string     Path to the dynamic config file
  -h, --help            help for dynamic-config
      --policy string   Name of the Policy to apply the DynamicConfig to
```

### Options inherited from parent commands

```
      --controller string      Address of Aperture controller
      --controller-ns string   Namespace in which the Aperture Controller is running
      --insecure               Allow insecure connection to controller
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
```

### SEE ALSO

- [aperturectl apply](/reference/aperturectl/apply/apply.md) - Apply Aperture Policies
