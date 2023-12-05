---
sidebar_label: Set
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_result-cache_set
---

<!-- markdownlint-disable -->

## aperturectl cloud result-cache set

Set a result cache entry

### Synopsis

Set a result cache entry

```
aperturectl cloud result-cache set [flags]
```

### Options

```
  -a, --agent-group string     Agent group
  -c, --control-point string   Control point
  -h, --help                   help for set
  -k, --key string             Key
  -t, --ttl int                TTL in seconds (default 600)
      --value string           Value
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

- [aperturectl cloud result-cache](/reference/aperture-cli/aperturectl/cloud/result-cache/result-cache.md) - Result Cache related commands
