# Latency Gradient Policy

## Configuration

[configuration]: # Configuration Marker

### Common

| Parameter Name      | Parameter Type | Default      | Description         |
| ------------------- | -------------- | ------------ | ------------------- |
| `common.policyName` | `string`       | `(required)` | Name of the policy. |

### Policy

| Parameter Name                                                 | Parameter Type                         | Default      | Description                                                                |
| -------------------------------------------------------------- | -------------------------------------- | ------------ | -------------------------------------------------------------------------- |
| `policy.evaluationInterval`                                    | `string`                               | `"0.1s"`     | How often should policy be re-evaluated.                                   |
| `policy.fluxMeter`                                             | `aperture.v1.FluxMeter`                | `(required)` | Flux Meter selector.                                                       |
| `policy.concurrencyLimiterSelector`                            | `aperture.spec.v1.Selector`            | `(required)` | Concurrency Limiter selector.                                              |
| `policy.classifiers`                                           | `[]aperture.spec.v1.Classifier`        | `[]`         | List of classification rules.                                              |
| `policy.components`                                            | `[]aperture.spec.v1.Component`         | `[]`         | List of additional circuit components.                                     |
| `policy.ema.window`                                            | `string`                               | `"1500s"`    | How far back to look when calculating moving average                       |
| `policy.ema.warmUpWindow`                                      | `string`                               | `"10s"`      | How much time to give circuit before we start calculating averages         |
| `policy.concurrencyLimiter.defaultWorkloadParameters.priority` | `int`                                  | `20`         | Workload parameters to use in case none of the configured workloads match. |
| `policy.concurrencyLimiter.workloads`                          | `[]aperture.spec.v1.SchedulerWorkload` | `[]`         | A list of additional workloads for the scheduler                           |

### Dashboard

| Parameter Name                    | Parameter Type | Default         | Description                            |
| --------------------------------- | -------------- | --------------- | -------------------------------------- |
| `dashboard.refreshInterval`       | `string`       | `"10s"`         | Refresh interval for dashboard panels. |
| `dashboard.datasourceName`        | `string`       | `"$datasource"` | Datasource name.                       |
| `dashboard.datasourceFilterRegex` | `string`       | `""`            | Datasource filter regex.               |
