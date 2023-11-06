---
sidebar_label: Get
hide_title: true
keywords:
  - aperturectl
  - aperturectl_dynamic-config_get
---

<!-- markdownlint-disable -->

## aperturectl dynamic-config get

Get Aperture DynamicConfig for a Policy.

### Synopsis

Use this command to get the Aperture DynamicConfig of a Policy.

```
aperturectl dynamic-config get POLICY_NAME [flags]
```

### Examples

```
aperture dynamic-config get rate-limiting
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

- [aperturectl dynamic-config](/reference/aperture-cli/aperturectl/dynamic-config/dynamic-config.md) - DynamicConfig of Aperture Policy related commands for the Controller
