---
sidebar_label: Discovery
hide_title: true
keywords:
  - aperturectl
  - aperturectl_discovery
---

<!-- markdownlint-disable -->

## aperturectl discovery

Discovery integrations

### Synopsis

Use this command to query information about active Discovery integrations

### Options

```
      --api-key string         FluxNinja Cloud API Key to be used when using Cloud Controller
      --config string          Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for discovery
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl discovery entities](/reference/aperturectl/discovery/entities/entities.md) - List AutoScale control points
