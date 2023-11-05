---
sidebar_label: Entities
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_discovery_entities
---

<!-- markdownlint-disable -->

## aperturectl cloud discovery entities

List AutoScale control points

### Synopsis

List AutoScale control points

```
aperturectl cloud discovery entities [flags]
```

### Examples

```
aperturectl cloud discovery entities

aperturectl cloud discovery entities --find-by="name=service1-demo-app-7dfdf9c698-4wmlt"

aperturectl cloud discovery entities --find-by=“ip=10.244.1.24”
```

### Options

```
      --agent-group string   Name of the agent group to list agents for
      --find-by string       Find entity by [name|ip]
  -h, --help                 help for entities
```

### Options inherited from parent commands

```
      --api-key string        Aperture Cloud User API Key to be used when using Cloud Controller
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud discovery](/reference/aperture-cli/aperturectl/cloud/discovery/discovery.md) - Discovery integrations
