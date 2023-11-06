---
sidebar_label: Install
hide_title: true
keywords:
  - aperturectl
  - aperturectl_install
---

<!-- markdownlint-disable -->

## aperturectl install

Install Aperture components

### Synopsis

Use this command to install Aperture Controller and Agent on your Kubernetes cluster.

### Options

```
      --dry-run              If set to true, only the manifests will be generated and no installation will be performed
  -h, --help                 help for install
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
      --namespace string     Namespace in which the component will be installed. Defaults to 'default' namespace (default "default")
      --values-file string   Values YAML file containing parameters to customize the installation
      --version string       Version of the Aperture
```

### SEE ALSO

- [aperturectl](/reference/aperture-cli/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
- [aperturectl install agent](/reference/aperture-cli/aperturectl/install/agent/agent.md) - Install Aperture Agent
- [aperturectl install controller](/reference/aperture-cli/aperturectl/install/controller/controller.md) - Install Aperture Controller
- [aperturectl install istioconfig](/reference/aperture-cli/aperturectl/install/istioconfig/istioconfig.md) - Install example Istio EnvoyFilter
