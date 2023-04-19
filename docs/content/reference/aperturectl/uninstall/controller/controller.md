---
sidebar_label: Controller
hide_title: true
keywords:
  - aperturectl
  - aperturectl_uninstall_controller
---

## aperturectl uninstall controller

Uninstall Aperture Controller

### Synopsis

Use this command to uninstall Aperture Controller and its dependencies from your Kubernetes cluster

```
aperturectl uninstall controller [flags]
```

### Examples

```
aperturectl uninstall controller

aperturectl uninstall controller --namespace=aperture
```

### Options

```
  -h, --help   help for controller
```

### Options inherited from parent commands

```
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
      --namespace string     Namespace from which the component will be uninstalled. Defaults to 'default' namespace (default "default")
      --timeout int          Timeout of waiting for uninstallation hooks completion (default 300)
      --values-file string   Values YAML file containing parameters to customize the installation
      --version string       Version of the Aperture (default "latest")
```

### SEE ALSO

- [aperturectl uninstall](/reference/aperturectl/uninstall/uninstall.md) - Uninstall Aperture components
