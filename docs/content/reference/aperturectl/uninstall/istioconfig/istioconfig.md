---
sidebar_label: Istioconfig
hide_title: true
keywords:
  - aperturectl
  - aperturectl_uninstall_istioconfig
---

## aperturectl uninstall istioconfig

Install example Istio EnvoyFilter

### Synopsis

Use this command to uninstall example Istio EnvoyFilter from your Kubernetes
cluster.

```
aperturectl uninstall istioconfig [flags]
```

### Examples

```
aperturectl uninstall istioconfig

aperturectl uninstall istioconfig --namespace=istio-system
```

### Options

```
  -h, --help   help for istioconfig
```

### Options inherited from parent commands

```
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
      --namespace string     Namespace from which the component will be uninstalled. Defaults to component name
      --timeout int          Timeout of waiting for uninstallation hooks completion (default 300)
      --version string       Version of the Aperture (default "latest")
```

### SEE ALSO

- [aperturectl uninstall](/reference/aperturectl/uninstall/uninstall.md) -
  Uninstall Aperture components
