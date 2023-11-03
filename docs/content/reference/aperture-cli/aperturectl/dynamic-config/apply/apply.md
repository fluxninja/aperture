---
sidebar_label: Apply
hide_title: true
keywords:
  - aperturectl
  - aperturectl_dynamic-config_apply
---

<!-- markdownlint-disable -->

## aperturectl dynamic-config apply

Apply Aperture DynamicConfig to a Policy

### Synopsis

Use this command to apply the Aperture DynamicConfig to a Policy.

```
aperturectl dynamic-config apply [flags]
```

### Examples

```
aperturectl dynamic-config apply --policy=rate-limiting --file=dynamic-config.yaml
```

### Options

```
      --file string     Path to the dynamic config file
  -h, --help            help for apply
      --policy string   Name of the Policy to apply the DynamicConfig to
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

- [aperturectl dynamic-config](/reference/aperture-cli/aperturectl/dynamic-config/dynamic-config.md) - DynamicConfig of Aperture Policy related commands for the Controller
