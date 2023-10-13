---
sidebar_label: Uninstall
hide_title: true
keywords:
  - aperturectl
  - aperturectl_uninstall
---

<!-- markdownlint-disable -->

## aperturectl uninstall

Uninstall Aperture components

### Synopsis

Use this command to uninstall Aperture Controller and Agent from your Kubernetes cluster.

### Options

```
  -h, --help                 help for uninstall
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
      --namespace string     Namespace from which the component will be uninstalled. Defaults to 'default' namespace (default "default")
      --timeout int          Timeout of waiting for uninstallation hooks completion (default 300)
      --values-file string   Values YAML file containing parameters to customize the installation
      --version string       Version of the Aperture
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl uninstall agent](/reference/aperturectl/uninstall/agent/agent.md) - Uninstall Aperture Agent
- [aperturectl uninstall controller](/reference/aperturectl/uninstall/controller/controller.md) - Uninstall Aperture Controller
- [aperturectl uninstall istioconfig](/reference/aperturectl/uninstall/istioconfig/istioconfig.md) - Install example Istio EnvoyFilter
