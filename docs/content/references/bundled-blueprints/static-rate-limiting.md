---
title: Static Rate Limiting Policy
---

```mdx-code-block
import {apertureVersion} from '../../apertureVersion.js';
```

## Blueprint Location

GitHub: <a
href={`https://github.com/fluxninja/aperture/tree/v${apertureVersion}/blueprints/lib/1.0/blueprints/static-rate-limiting`}>static-rate-limiting</a>

## Introduction

This blueprint provides a simple static rate limiting policy and a dashboard.

## Configuration

<!-- Configuration Marker -->

### Common

| Parameter Name      | Parameter Type | Default      | Description         |
| ------------------- | -------------- | ------------ | ------------------- |
| `common.policyName` | `string`       | `(required)` | Name of the policy. |

### Policy

| Parameter Name                   | Parameter Type                  | Default      | Description                                                 |
| -------------------------------- | ------------------------------- | ------------ | ----------------------------------------------------------- |
| `policy.evaluationInterval`      | `string`                        | `"300s"`     | How often should the policy be re-evaluated                 |
| `policy.rateLimit`               | `float64`                       | `(required)` | How many requests per `policy.limitResetInterval` to accept |
| `policy.rateLimiterFlowSelector` | `aperture.spec.v1.FlowSelector` | `(required)` | A flow selector to match requests against                   |
| `policy.limitResetInterval`      | `string`                        | `"1s"`       | The window for `policy.rateLimit`                           |
| `policy.labelKey`                | `string`                        | `(required)` | What flow label to use for rate limiting                    |

#### Overrides

| Parameter Name     | Parameter Type                           | Default | Description                                     |
| ------------------ | ---------------------------------------- | ------- | ----------------------------------------------- |
| `policy.overrides` | `[]aperture.spec.v1.RateLimiterOverride` | `[]`    | A list of limit overrides for the rate limiter. |

#### Lazy Sync

| Parameter Name            | Parameter Type | Default | Description                                                          |
| ------------------------- | -------------- | ------- | -------------------------------------------------------------------- |
| `policy.lazySync.enabled` | `boolean`      | `true`  | Enable lazy syncing.                                                 |
| `policy.lazySync.numSync` | `integer`      | `5`     | Number of times to lazy sync within the `policy.limitResetInterval`. |

### Dashboard

| Parameter Name              | Parameter Type | Default | Description                            |
| --------------------------- | -------------- | ------- | -------------------------------------- |
| `dashboard.refreshInterval` | `string`       | `"10s"` | Refresh interval for dashboard panels. |

#### Datasource

| Parameter Name                     | Parameter Type | Default         | Description              |
| ---------------------------------- | -------------- | --------------- | ------------------------ |
| `dashboard.datasource.name`        | `string`       | `"$datasource"` | Datasource name.         |
| `dashboard.datasource.filterRegex` | `string`       | `""`            | Datasource filter regex. |
