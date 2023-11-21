---
sidebar_label: Decisions
hide_title: true
keywords:
  - aperturectl
  - aperturectl_cloud_decisions
---

<!-- markdownlint-disable -->

## aperturectl cloud decisions

Get Aperture Decisions

### Synopsis

Use this command to get the Aperture Decisions.

```
aperturectl cloud decisions [flags]
```

### Examples

```

	aperturectl cloud decisions --all
	aperturectl cloud decisions --decision-type="load_scheduler"
```

### Options

```
      --access-token string    User Access Token to be used while connecting to Aperture Cloud
      --all                    Get all decisions
      --config string          Path to the Aperture config file. Defaults to '~/.aperturectl/config' or $APERTURE_CONFIG
      --controller string      Address of Aperture Cloud Controller
      --decision-type string   Type of the decision to get (load_scheduler, rate_limiter, quota_scheduler, pod_scaler, sampler)
  -h, --help                   help for decisions
      --insecure               Allow connection to controller running without TLS
      --project-name string    Aperture Cloud Project Name to be used when using Cloud Controller
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl cloud](/reference/aperture-cli/aperturectl/cloud/cloud.md) - Commands to communicate with the Cloud Controller
