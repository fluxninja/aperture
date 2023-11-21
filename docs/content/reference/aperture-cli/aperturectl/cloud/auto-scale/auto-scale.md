---
sidebar_label: Auto-Scale
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_auto-scale
---

<!-- markdownlint-disable -->

## aperturectl cloud auto-scale

AutoScale integrations

### Synopsis

Use this command to query information about active AutoScale integrations

### Options

```
      --access-token string   User Access Token to be used while connecting to Aperture Cloud
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
  -h, --help                  help for auto-scale
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperture-cli/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
- [aperturectl cloud auto-scale control-points](/reference/aperture-cli/aperturectl/cloud/auto-scale/control-points/control-points.md) - List AutoScale control points
