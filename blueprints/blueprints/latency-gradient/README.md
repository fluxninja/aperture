# Latency Gradient Policy

## Configuration

[configuration]: # Configuration Marker

### Latency Gradient Policy

| Parameter Name              | Parameter Type | Default      | Description                                                             |
| --------------------------- | -------------- | ------------ | ----------------------------------------------------------------------- |
| `policy.policyName`         | `string`       | `(required)` | A name of the policy, used within PromQL queries for fluxmeter metrics. |
| `policy.evaluationInterval` | `string`       | `"0.5s"`     | How often should policy be re-evaluated.                                |

#### Flux Meter Selector

| Parameter Name             | Parameter Type         | Default      | Description          |
| -------------------------- | ---------------------- | ------------ | -------------------- |
| `policy.fluxMeterSelector` | `aperture.v1.Selector` | `(required)` | Flux Meter selector. |

#### Flux Meters

| Parameter Name                                   | Parameter Type                       | Default      | Description                                 |
| ------------------------------------------------ | ------------------------------------ | ------------ | ------------------------------------------- |
| `policy.fluxMeters`                              | `map[string]aperture.v1.FluxMeter`   | `{}`         | Mappings of fluxMeterName to fluxMeter.     |
| `policy.fluxMeters[policyName].attributeKey`     | `string`                             | `(required)` | Key of the attribute in access log or span. |
| `policy.fluxMeters[policyName].histogramBuckets` | `aperture.v1.FluxMeterStaticBuckets` | `(required)` | Flux Meter static histogram buckets.        |

#### Concurrency Limiter Selector

| Parameter Name                      | Parameter Type         | Default      | Description                   |
| ----------------------------------- | ---------------------- | ------------ | ----------------------------- |
| `policy.concurrencyLimiterSelector` | `aperture.v1.Selector` | `(required)` | Concurrency Limiter selector. |

#### Classification Rules

| Parameter Name       | Parameter Type | Default | Description                   |
| -------------------- | -------------- | ------- | ----------------------------- |
| `policy.classifiers` | `string`       | `[]`    | List of classification rules. |

#### Exponential Moving Average configuration

| Parameter Name            | Parameter Type | Default   | Description                                                        |
| ------------------------- | -------------- | --------- | ------------------------------------------------------------------ |
| `policy.ema.window`       | `string`       | `"1500s"` | How far back to look when calculating moving average               |
| `policy.ema.warmUpWindow` | `string`       | `"10s"`   | How much time to give circuit before we start calculating averages |

#### Concurrency Limiter

| Parameter Name                                                 | Parameter Type                    | Default | Description                                                                |
| -------------------------------------------------------------- | --------------------------------- | ------- | -------------------------------------------------------------------------- |
| `policy.concurrencyLimiter.defaultWorkloadParameters.priority` | `int`                             | `20`    | Workload parameters to use in case none of the configured workloads match. |
| `policy.concurrencyLimiter.workloads`                          | `[]aperture.v1.SchedulerWorkload` | `[]`    | A list of additional workloads for the scheduler                           |

### FluxMeter Dashboard

| Parameter Name         | Parameter Type | Default      | Description                                                               |
| ---------------------- | -------------- | ------------ | ------------------------------------------------------------------------- |
| `dashboard.policyName` | `string`       | `(required)` | A name of the policy used as a promQL query filter for flux meter metrics |
