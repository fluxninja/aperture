---
sidebar_label: Entities
hide_title: true
keywords:
  - aperturectl
  - aperturectl_discovery_entities
---

<!-- markdownlint-disable -->

## aperturectl discovery entities

List AutoScale control points

### Synopsis

List AutoScale control points

```
aperturectl discovery entities [flags]
```

### Examples

```
aperturectl discovery entities --kube

aperturectl discovery entities --kube --find-by="name=service1-demo-app-7dfdf9c698-4wmlt"

aperturectl discovery entities --kube --find-by=“ip=10.244.1.24”
```

### Options

```
      --find-by string   Find entity by [name|ip]
  -h, --help             help for entities
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

- [aperturectl discovery](/reference/aperturectl/discovery/discovery.md) - Discovery integrations
