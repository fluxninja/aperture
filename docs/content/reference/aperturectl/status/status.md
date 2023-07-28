---
sidebar_label: Status
hide_title: true
keywords:
  - aperturectl
  - aperturectl_status
---

<!-- markdownlint-disable -->

## aperturectl status

Get Jobs status

### Synopsis

Use this command to get the status of internal jobs.

```
aperturectl status [flags]
```

### Examples

```

	aperturectl status

```

### Options

```
      --api-key string         FluxNinja API Key to be used when using Cloud Controller
      --config string          Path to the Aperture config file
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for status
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
