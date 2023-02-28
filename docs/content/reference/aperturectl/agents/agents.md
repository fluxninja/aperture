---
sidebar_label: Agents
hide_title: true
keywords:
  - aperturectl
  - aperturectl_agents
---

## aperturectl agents

List connected agents

### Synopsis

List connected agents

```
aperturectl agents {--kube | --controller ADDRESS} [flags]
```

### Options

```
      --controller string    Address of Aperture controller
  -h, --help                 help for agents
      --insecure             Allow insecure connection to controller
      --kube                 Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string   Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
