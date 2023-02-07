---
sidebar_label: Controller
hide_title: true
keywords:
  - aperturectl
  - aperturectl_install_controller
---

## aperturectl install controller

Install Aperture Controller

### Synopsis

Use this command to install Aperture Controller and its dependencies on your Kubernetes cluster.
Refer https://github.com/fluxninja/aperture/blob/v0.0.1/manifests/charts/aperture-controller/README.md for list of configurable parameters for preparing values file.

```
aperturectl install controller [flags]
```

### Examples

```
aperturectl install controller --values-file=values.yaml

aperturectl install controller --values-file=values.yaml --namespace=aperture
```

### Options

```
  -h, --help   help for controller
```

### Options inherited from parent commands

```
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
      --namespace string     Namespace in which the component will be installed. Defaults to component name
      --values-file string   Values YAML file containing parameters to customize the installation
      --version string       Version of the Aperture to uninstall. Defaults to latest (default "latest")
```

### SEE ALSO

- [aperturectl install](/reference/aperturectl/install/install.md) - Install Aperture components
