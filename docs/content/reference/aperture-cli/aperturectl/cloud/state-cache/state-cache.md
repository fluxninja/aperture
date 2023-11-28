---
sidebar_label: State-Cache
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_state-cache
---

<!-- markdownlint-disable -->

## aperturectl cloud state-cache

State Cache related commands

### Synopsis

Use this command to interact with Aperture State Cache.

### Options

```
      --access-token string   User Access Token to be used while connecting to Aperture Cloud
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
  -h, --help                  help for state-cache
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperture-cli/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
- [aperturectl cloud state-cache delete](/reference/aperture-cli/aperturectl/cloud/state-cache/delete/delete.md) - Delete a state cache entry
- [aperturectl cloud state-cache get](/reference/aperture-cli/aperturectl/cloud/state-cache/get/get.md) - Get a state cache entry
- [aperturectl cloud state-cache set](/reference/aperture-cli/aperturectl/cloud/state-cache/set/set.md) - Set a state cache entry
