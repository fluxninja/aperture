---
sidebar_label: Result-Cache
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_result-cache
---

<!-- markdownlint-disable -->

## aperturectl cloud result-cache

Result Cache related commands

### Synopsis

Use this command to interact with Aperture's Result Cache.

### Options

```
      --access-token string   User Access Token to be used while connecting to Aperture Cloud
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
  -h, --help                  help for result-cache
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperture-cli/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
- [aperturectl cloud result-cache delete](/reference/aperture-cli/aperturectl/cloud/result-cache/delete/delete.md) - Delete a result cache entry
- [aperturectl cloud result-cache get](/reference/aperture-cli/aperturectl/cloud/result-cache/get/get.md) - Get a result cache entry
- [aperturectl cloud result-cache set](/reference/aperture-cli/aperturectl/cloud/result-cache/set/set.md) - Set a result cache entry
