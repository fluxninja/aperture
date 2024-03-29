---
sidebar_label: Status
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_status
---

<!-- markdownlint-disable -->

## aperturectl cloud status

Get Jobs status

### Synopsis

Use this command to get the status of internal jobs.

```
aperturectl cloud status [flags]
```

### Examples

```

	aperturectl cloud status

```

### Options

```
      --access-token string   User Access Token to be used while connecting to Aperture Cloud
      --config string         Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string     Address of Aperture Cloud Controller
  -h, --help                  help for status
      --insecure              Allow connection to controller running without TLS
      --project-name string   Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify           Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperture-cli/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
