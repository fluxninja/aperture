---
sidebar_label: Controller
hide_title: true
keywords:
  - aperturectl
  - aperturectl_install_controller
---

<!-- markdownlint-disable -->

## aperturectl install controller

Install Aperture Controller

### Synopsis

Use this command to install Aperture Controller and its dependencies on your Kubernetes cluster.
Refer https://artifacthub.io/packages/helm/aperture/aperture-controller#parameters for list of configurable parameters for preparing values file.

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
      --generate-cert   Generate self signed certificates for Aperture Controller
  -h, --help            help for controller
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
