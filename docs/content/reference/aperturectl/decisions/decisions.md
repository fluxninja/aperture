---
sidebar_label: Decisions
hide_title: true
keywords:
  - aperturectl
  - aperturectl_decisions
---

<!-- markdownlint-disable -->

## aperturectl decisions

Get Aperture Decisions

### Synopsis

Use this command to get the Aperture Decisions.

```
aperturectl decisions [flags]
```

### Examples

```

	aperturectl decisions --all
	aperturectl decisions --decision-type="load_scheduler"
```

### Options

```
      --all                    Get all decisions
      --api-key string         FluxNinja API Key to be used when using Cloud Controller
      --config string          Path to the Aperture config file
      --controller string      Address of Aperture Controller
      --controller-ns string   Namespace in which the Aperture Controller is running
      --decision-type string   Type of the decision to get (load_scheduler, rate_limiter, quota_scheduler, pod_scaler, sampler)
  -h, --help                   help for decisions
      --insecure               Allow connection to controller running without TLS
      --kube                   Find controller in Kubernetes cluster, instead of connecting directly
      --kube-config string     Path to the Kubernetes cluster config. Defaults to '~/.kube/config' or $KUBECONFIG
      --skip-verify            Skip TLS certificate verification while connecting to controller
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
