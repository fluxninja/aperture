# Latency Gradient Policy

## Configuration

[configuration]: # Configuration Marker

### Latency Gradient Policy

| Parameter Name              | Parameter Type | Default      | Description                                                             |
| --------------------------- | -------------- | ------------ | ----------------------------------------------------------------------- |
| `policy.policyName`         | `string`       | `(required)` | A name of the policy, used within PromQL queries for fluxmeter metrics. |
| `policy.evaluationInterval` | `string`       | `"0.5s"`     | How often should policy be re-evaluated.                                |

#### Flux Meter

| Parameter Name     | Parameter Type          | Default      | Description          |
| ------------------ | ----------------------- | ------------ | -------------------- |
| `policy.fluxMeter` | `aperture.v1.FluxMeter` | `(required)` | Flux Meter selector. |

#### Concurrency Limiter Selector

| Parameter Name                      | Parameter Type              | Default      | Description                   |
| ----------------------------------- | --------------------------- | ------------ | ----------------------------- |
| `policy.concurrencyLimiterSelector` | `aperture.spec.v1.Selector` | `(required)` | Concurrency Limiter selector. |

#### Classification Rules

| Parameter Name       | Parameter Type                  | Default | Description                   |
| -------------------- | ------------------------------- | ------- | ----------------------------- |
| `policy.classifiers` | `[]aperture.spec.v1.Classifier` | `[]`    | List of classification rules. |

#### Additional Circuit Components

| Parameter Name      | Parameter Type                 | Default | Description                            |
| ------------------- | ------------------------------ | ------- | -------------------------------------- |
| `policy.components` | `[]aperture.spec.v1.Component` | `[]`    | List of additional circuit components. |

#### Exponential Moving Average configuration

| Parameter Name            | Parameter Type | Default   | Description                                                        |
| ------------------------- | -------------- | --------- | ------------------------------------------------------------------ |
| `policy.ema.window`       | `string`       | `"1500s"` | How far back to look when calculating moving average               |
| `policy.ema.warmUpWindow` | `string`       | `"10s"`   | How much time to give circuit before we start calculating averages |

#### Concurrency Limiter

| Parameter Name                                                 | Parameter Type                         | Default | Description                                                                |
| -------------------------------------------------------------- | -------------------------------------- | ------- | -------------------------------------------------------------------------- |
| `policy.concurrencyLimiter.defaultWorkloadParameters.priority` | `int`                                  | `20`    | Workload parameters to use in case none of the configured workloads match. |
| `policy.concurrencyLimiter.workloads`                          | `[]aperture.spec.v1.SchedulerWorkload` | `[]`    | A list of additional workloads for the scheduler                           |

### FluxMeter Dashboard

| Parameter Name         | Parameter Type | Default      | Description                                                               |
| ---------------------- | -------------- | ------------ | ------------------------------------------------------------------------- |
| `dashboard.policyName` | `string`       | `(required)` | A name of the policy used as a promQL query filter for flux meter metrics |
