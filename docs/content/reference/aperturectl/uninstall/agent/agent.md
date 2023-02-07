---
sidebar_label: Agent
hide_title: true
keywords:
  - aperturectl
  - aperturectl_uninstall_agent
---

## aperturectl uninstall agent

Uninstall Aperture Agent

### Synopsis

Use this command to uninstall Aperture Agent from your Kubernetes cluster

```
aperturectl uninstall agent [flags]
```

### Examples

```
aperturectl uninstall agent

aperturectl uninstall agent --namespace=aperture
```

### Options

```
  -h, --help   help for agent
```

### Options inherited from parent commands

```
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config'
      --namespace string     Namespace from which the component will be uninstalled. Defaults to component name
      --version string       Version of the Aperture to uninstall. Defaults to latest (default "latest")
```

### SEE ALSO

- [aperturectl uninstall](/reference/aperturectl/uninstall/uninstall.md) - Uninstall Aperture components
