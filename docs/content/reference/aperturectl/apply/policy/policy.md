---
sidebar_label: Policy
hide_title: true
keywords:
  - aperturectl
  - aperturectl_apply_policy
---

<!-- markdownlint-disable -->

## aperturectl apply policy

Apply Aperture Policy to the cluster

### Synopsis

Use this command to apply the Aperture Policy to the cluster.

```
aperturectl apply policy [flags]
```

### Examples

```
aperturectl apply policy --file=policies/rate-limiting.yaml

aperturectl apply policy --dir=policies
```

### Options

```
      --dir string    Path to directory containing Aperture Policy files
      --file string   Path to Aperture Policy file
  -f, --force         Force apply policy even if it already exists
  -h, --help          help for policy
  -s, --select-all    Apply all policies in the directory
```

### Options inherited from parent commands

```
      --api-key string         FluxNinja ARC API Key to be used when using Cloud Controller
      --config string          Path to the Aperture config file
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl apply](/reference/aperturectl/apply/apply.md) - Apply Aperture Policies
