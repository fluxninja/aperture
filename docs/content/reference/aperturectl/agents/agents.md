---
sidebar_label: Agents
hide_title: true
keywords:
  - aperturectl
  - aperturectl_agents
---

<!-- markdownlint-disable -->

## aperturectl agents

List connected agents

### Synopsis

List connected agents

```
aperturectl agents [flags]
```

### Options

```
      --api-key string         FluxNinja Cloud API Key to be used when using Cloud Controller
      --config string          Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
  -h, --help                   help for agents
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
