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

Concurrency limiter is responsible for.... TODO

`concurrencyLimiter.workloads` is a list of `aperture.v1.ServiceWorkloadAndLabelMatcher` objects, which can be generated
using aperture libsonnet library (which allows us to do some static checks for validity):

```jsonnet
local aperture = import 'github.com/fluxninja/aperture/libsonnet/1.0/main.libsonnet';

local Workload = aperture.v1.SchedulerWorkload;
local LabelMatcher = aperture.v1.LabelMatcher;
local WorkloadMatcher = aperture.v1.ServiceWorkloadAndLabelMatcher; // Make a local typedef for quicker access

{
  concurrencyLimiter+: {
    workloads: [
      WorkloadMatcher.new(
        workload=Workload.new() + Workload.withPriority(50) + Workload.withTimeout('0.005s')
        label_matcher=LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'guest' })),
      WorkloadWithLabelMatcher.new(
        workload=Workload.new() + Workload.withPriority(200) + Workload.withTimeout('0.005s'),
        label_matcher=LabelMatcher.withMatchLabels({ 'http.request.header.user_type': 'subscriber' }))
    ]
  }
}
```

Or it can be passed as a list of jsonnet objects directly:

```jsonnet
{
  concurrencyLimiter+: {
    workloads: [
      {
        label_matcher: {
          match_labels: {
            "http.request.header.user_type": "guest"
          }
        },
        workload: { priority: 50, timeout: "0.005s" }
      },
      {
        label_matcher: {
          match_labels: {
            "http.request.header.user_type": "subscriber"
          }
        },
        workload: { priority: 200, timeout: "0.005s" }
      }
    ]
  }
}
```

| Parameter Name                                       | Parameter Type                                   | Default | Description                                      |
| ---------------------------------------------------- | ------------------------------------------------ | ------- | ------------------------------------------------ |
| `policy.concurrencyLimiter.defaultWorkload.priority` | `int`                                            | `20`    | TODO                                             |
| `policy.concurrencyLimiter.workloads`                | `[]aperture.v1.SchedulerWorkloadAndLabelMatcher` | `[]`    | A list of additional workloads for the scheduler |

### FluxMeter Dashboard

| Parameter Name         | Parameter Type | Default      | Description                                                               |
| ---------------------- | -------------- | ------------ | ------------------------------------------------------------------------- |
| `dashboard.policyName` | `string`       | `(required)` | A name of the policy used as a promQL query filter for flux meter metrics |
