---
sidebar_label: Istioconfig
hide_title: true
keywords:
  - aperturectl
  - aperturectl_install_istioconfig
---

## aperturectl install istioconfig

Install example Istio EnvoyFilter

### Synopsis

Use this command to install example Istio EnvoyFilter on your Kubernetes cluster.
Refer https://artifacthub.io/packages/helm/aperture/istioconfig#parameters for list of configurable parameters for preparing values file.

```
aperturectl install istioconfig [flags]
```

### Examples

```
aperturectl install istioconfig --values-file=values.yaml

aperturectl install istioconfig --values-file=values.yaml --namespace=istio-system
```

### Options

```
  -h, --help   help for istioconfig
```

### Options inherited from parent commands

```
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
      --namespace string     Namespace in which the component will be installed. Defaults to 'default' namespace (default "default")
      --values-file string   Values YAML file containing parameters to customize the installation
      --version string       Version of the Aperture (default "latest")
```

### SEE ALSO

- [aperturectl install](/reference/aperturectl/install/install.md) - Install Aperture components
