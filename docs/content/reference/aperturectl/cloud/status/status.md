---
sidebar_label: Status
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_status
---

<!-- markdownlint-disable -->

## aperturectl cloud status

Get Jobs status

### Synopsis

Use this command to get the status of internal jobs.

```
aperturectl cloud status [flags]
```

### Examples

```

	aperturectl cloud status

```

### Options

```
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for status
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
