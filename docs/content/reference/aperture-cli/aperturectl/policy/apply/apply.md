---
sidebar_label: Apply
hide_title: true
keywords:
  - aperturectl
  - aperturectl_policy_apply
---

<!-- markdownlint-disable -->

## aperturectl policy apply

Apply Aperture Policy to the cluster

### Synopsis

Use this command to apply the Aperture Policy to the cluster.

```
aperturectl policy apply [flags]
```

### Examples

```
aperturectl policy apply --file=policies/rate-limiting.yaml

aperturectl apply policy --dir=policies
```

### Options

```
      --dir string    Path to directory containing Aperture Policy files
      --file string   Path to Aperture Policy file
  -f, --force         Force apply policy even if it already exists
  -h, --help          help for apply
  -s, --select-all    Apply all policies in the directory
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
