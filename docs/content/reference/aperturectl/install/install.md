---
sidebar_label: Install
hide_title: true
keywords:
  - aperturectl
  - aperturectl_install
---

## aperturectl install

Install Aperture components

### Synopsis

Use this command to install Aperture Controller and Agent on your Kubernetes
cluster.

### Options

```
  -h, --help                 help for install
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
      --namespace string     Namespace in which the component will be installed. Defaults to component name
      --values-file string   Values YAML file containing parameters to customize the installation
      --version string       Version of the Aperture (default "latest")
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool
  to interact with Aperture
- [aperturectl install agent](/reference/aperturectl/install/agent/agent.md) -
  Install Aperture Agent
- [aperturectl install controller](/reference/aperturectl/install/controller/controller.md) -
  Install Aperture Controller
- [aperturectl install istioconfig](/reference/aperturectl/install/istioconfig/istioconfig.md) -
  Install example Istio EnvoyFilter
