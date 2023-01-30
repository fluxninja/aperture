---
title: Static Rate Limiting Policy
---

```mdx-code-block
import {apertureVersion} from '../../../apertureVersion.js';
```

## Blueprint Location

GitHub: <a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/blueprints/lib/1.0/policies/static-rate-limiting`}>static-rate-limiting</a>

## Introduction

This blueprint provides a simple static rate limiting policy and a dashboard.

## Configuration

<!-- Configuration Marker -->

### Common

| Parameter Name       | Parameter Type | Default      | Description         |
| -------------------- | -------------- | ------------ | ------------------- |
| `common.policy_name` | `string`       | `(required)` | Name of the policy. |

### Policy

| Parameter Name               | Parameter Type                  | Default  | Description                                 |
| ---------------------------- | ------------------------------- | -------- | ------------------------------------------- |
| `policy.evaluation_interval` | `string`                        | `"300s"` | How often should the policy be re-evaluated |
| `policy.classifiers`         | `[]aperture.spec.v1.Classifier` | `[]`     | List of classification rules.               |

#### Rate Limiter

| Parameter Name                             | Parameter Type                              | Default                                                                                                    | Description                                                                            |
| ------------------------------------------ | ------------------------------------------- | ---------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------- |
| `policy.rate_limiter.rate_limit`           | `float64`                                   | `(required)`                                                                                               | Number of requests per `policy.rate_limiter.parameters.limit_reset_interval` to accept |
| `policy.rate_limiter.flow_selector`        | `aperture.spec.v1.FlowSelector`             | `(required)`                                                                                               | A flow selector to match requests against                                              |
| `policy.rate_limiter.parameters`           | `aperture.spec.v1.RateLimiterParameters`    | `{'label_key': 'FAKE-VALUE', 'lazy_sync': {'enabled': True, 'num_sync': 5}, 'limit_reset_interval': '1s'}` | Parameters.                                                                            |
| `policy.rate_limiter.parameters.label_key` | `string`                                    | `(required)`                                                                                               | Flow label to use for rate limiting.                                                   |
| `policy.rate_limiter.dynamic_config`       | `aperture.spec.v1.RateLimiterDefaultConfig` | `{'overrides': []}`                                                                                        | Dynamic configuration for rate limiter that can be applied at the runtime.             |

### Dashboard

| Parameter Name               | Parameter Type | Default | Description                            |
| ---------------------------- | -------------- | ------- | -------------------------------------- |
| `dashboard.refresh_interval` | `string`       | `"10s"` | Refresh interval for dashboard panels. |

#### Datasource

| Parameter Name                      | Parameter Type | Default         | Description              |
| ----------------------------------- | -------------- | --------------- | ------------------------ |
| `dashboard.datasource.name`         | `string`       | `"$datasource"` | Datasource name.         |
| `dashboard.datasource.filter_regex` | `string`       | `""`            | Datasource filter regex. |
