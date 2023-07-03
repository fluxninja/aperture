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

	aperturectl decisions --etcd-host="127.0.0.1" --etcd-port="2379" --all
	aperturectl decisions --etcd-host="127.0.0.1" --etcd-port="2379" --decision-type="load_scheduler"
```

### Options

```
      --all                    Get all decisions
      --decision-type string   Type of the decision to get (load_scheduler, rate_limiter, quota_scheduler, pod_scaler, sampler)
      --etcd-host string       Etcd host (default "localhost")
      --etcd-port string       Etcd port (default "2379")
  -h, --help                   help for decisions
```

### SEE ALSO

- [aperturectl](/reference/aperturectl/aperturectl.md) - aperturectl - CLI tool to interact with Aperture
