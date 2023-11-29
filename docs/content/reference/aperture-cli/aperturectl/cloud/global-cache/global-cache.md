---
sidebar_label: Global-Cache
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_global-cache
---

<!-- markdownlint-disable -->

## aperturectl cloud global-cache

Global Cache related commands

### Synopsis

Use this command to interact with Aperture's Global Cache.

### Options

```
      --access-token string   User Access Token to be used while connecting to Aperture Cloud
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
  -h, --help                  help for global-cache
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperture-cli/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
- [aperturectl cloud global-cache delete](/reference/aperture-cli/aperturectl/cloud/global-cache/delete/delete.md) - Delete a global cache entry
- [aperturectl cloud global-cache get](/reference/aperture-cli/aperturectl/cloud/global-cache/get/get.md) - Get a global cache entry
- [aperturectl cloud global-cache set](/reference/aperture-cli/aperturectl/cloud/global-cache/set/set.md) - Set a global cache entry
