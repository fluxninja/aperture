# Latency Gradient Policy

## Configuration

[configuration]: # Configuration Marker

### Common

| Parameter Name      | Parameter Type | Default      | Description         |
| ------------------- | -------------- | ------------ | ------------------- |
| `common.policyName` | `string`       | `(required)` | Name of the policy. |

### Policy

| Parameter Name                          | Parameter Type                  | Default      | Description                            |
| --------------------------------------- | ------------------------------- | ------------ | -------------------------------------- |
| `policy.fluxMeter`                      | `aperture.spec.v1.FluxMeter`    | `(required)` | Flux Meter.                            |
| `policy.concurrencyLimiterFlowSelector` | `aperture.spec.v1.FlowSelector` | `(required)` | Concurrency Limiter flow selector.     |
| `policy.classifiers`                    | `[]aperture.spec.v1.Classifier` | `[]`         | List of classification rules.          |
| `policy.components`                     | `[]aperture.spec.v1.Component`  | `[]`         | List of additional circuit components. |

#### Concurrency Limiter

| Parameter Name                                                 | Parameter Type                         | Default | Description                                                                                    |
| -------------------------------------------------------------- | -------------------------------------- | ------- | ---------------------------------------------------------------------------------------------- |
| `policy.concurrencyLimiter.autoTokens`                         | `bool`                                 | `true`  | Whether tokens for workloads are computed dynamically or set statically by the user.           |
| `policy.concurrencyLimiter.timeoutFactor`                      | `float64`                              | `0.5`   | The maximum time a request can wait for tokens as a factor of tokens for a flow in a workload. |
| `policy.concurrencyLimiter.defaultWorkloadParameters.priority` | `int`                                  | `20`    | Workload parameters to use in case none of the configured workloads match.                     |
| `policy.concurrencyLimiter.workloads`                          | `[]aperture.spec.v1.SchedulerWorkload` | `[]`    | A list of additional workloads for the scheduler.                                              |

#### Dynamic Config

| Parameter Name                | Parameter Type | Default | Description                                                                                                                                                                                                                           |
| ----------------------------- | -------------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `policy.dynamicConfig.dryRun` | `bool`         | `false` | Decides whether this Policy runs in dry-run mode I.E. no traffic dropped. The signals would show how this Policy behaves and when it decides to drop any traffic. Useful for evaluating a Policy without disrupting any real traffic. |

#### Constants

| Parameter Name                                        | Parameter Type | Default | Description                                                                                                                                                                                                                                                                                   |
| ----------------------------------------------------- | -------------- | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `policy.constants.latencyToleranceMultiplier`         | `float64`      | `1.1`   | Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.                                                         |
| `policy.constants.latencyEMALimitMultiplier`          | `float64`      | `2.0`   | Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.                                                                                                                                                                                            |
| `policy.constants.concurrencyLimitMultiplier`         | `float64`      | `2.0`   | Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.                                                                               |
| `policy.constants.concurrencyLinearIncrement`         | `float64`      | `5.0`   | Linear increment to concurrency in each execution tick when the system is not in overloaded state.                                                                                                                                                                                            |
| `policy.constants.concurrencySQRTIncrementMultiplier` | `float64`      | `1`     | Scale factor to multiply square root of current accepted concurrrency. This, along with concurrencyLinearIncrement helps calculate overall concurrency increment in each tick. Concurrency is rapidly ramped up in each execution cycle during normal (non-overload) state (integral effect). |

#### EMA

| Parameter Name                | Parameter Type | Default   | Description                                                                                   |
| ----------------------------- | -------------- | --------- | --------------------------------------------------------------------------------------------- |
| `policy.ema.window`           | `string`       | `"1500s"` | How far back to look when calculating moving average.                                         |
| `policy.ema.warmUpWindow`     | `string`       | `"60s"`   | How much time to give circuit to learn the average value before we start emitting EMA values. |
| `policy.ema.correctionFactor` | `string`       | `0.95`    | Factor that is applied to the EMA value when it's above the maximum envelope.                 |

#### Gradient Controller

| Parameter Name                | Parameter Type | Default | Description                                                                                         |
| ----------------------------- | -------------- | ------- | --------------------------------------------------------------------------------------------------- |
| `policy.gradient.slope`       | `float64`      | `-1`    | Gradient that adjusts the response of the controller based on current latency and setpoint latency. |
| `policy.gradient.minGradient` | `float64`      | `0.1`   | Minimum gradient cap.                                                                               |
| `policy.gradient.maxGradient` | `float64`      | `1.0`   | Maximum gradient cap.                                                                               |

### Dashboard

| Parameter Name              | Parameter Type | Default | Description                            |
| --------------------------- | -------------- | ------- | -------------------------------------- |
| `dashboard.refreshInterval` | `string`       | `"10s"` | Refresh interval for dashboard panels. |

#### Datasource

| Parameter Name                     | Parameter Type | Default         | Description              |
| ---------------------------------- | -------------- | --------------- | ------------------------ |
| `dashboard.datasource.name`        | `string`       | `"$datasource"` | Datasource name.         |
| `dashboard.datasource.filterRegex` | `string`       | `""`            | Datasource filter regex. |
