---
sidebar_label: Preview
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_flow-control_preview
---

<!-- markdownlint-disable -->

## aperturectl cloud flow-control preview

Preview control points

### Synopsis

Preview samples of flow labels or HTTP requests on control points

```
aperturectl cloud flow-control preview [--http] CONTROL_POINT [flags]
```

### Options

```
      --agent-group string   Agent group (default "default")
  -h, --help                 help for preview
      --http                 Preview HTTP requests instead of flow labels
      --samples int          Number of samples to collect (default 10)
      --service string       Service FQDN (default "any")
```

### Options inherited from parent commands

```
      --access-token string   User Access Token to be used while connecting to Aperture Cloud
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud flow-control](/reference/aperture-cli/aperturectl/cloud/flow-control/flow-control.md) - Flow Control integrations
