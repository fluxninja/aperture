---
sidebar_label: Archive
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_policy_archive
---

<!-- markdownlint-disable -->

## aperturectl cloud policy archive

Archive Aperture Policy from the Aperture Cloud Controller

### Synopsis

Use this command to archive the Aperture Policy from the Aperture Cloud Controller.

```
aperturectl cloud policy archive POLICY_NAME [flags]
```

### Examples

```
aperturectl cloud policy archive POLICY_NAME
```

### Options

```
  -h, --help   help for archive
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

- [aperturectl cloud policy](/reference/aperture-cli/aperturectl/cloud/policy/policy.md) - Aperture Policy related commands for the Cloud Controller
