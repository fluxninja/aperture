---
sidebar_label: Agents
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_agents
---

<!-- markdownlint-disable -->

## aperturectl cloud agents

List connected agents

### Synopsis

List connected agents

```
aperturectl cloud agents [flags]
```

### Options

```
      --access-token string   User Access Token to be used while connecting to Aperture Cloud
      --agent-group string    Name of the agent group to list agents for
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
  -h, --help                  help for agents
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperture-cli/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
