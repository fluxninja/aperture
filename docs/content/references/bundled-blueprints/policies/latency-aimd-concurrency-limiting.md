---
title: Latency Gradient Concurrency Limiting Policy
---

```mdx-code-block
import {apertureVersion} from '../../../apertureVersion.js';
```

## Blueprint Location

GitHub: <a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/blueprints/lib/1.0/policies/latency-aimd-concurrency-limiting`}>latency-aimd-concurrency-limiting</a>

## Introduction

This policy detect overloads/cascading failures by comparing the real-time
latency with it's exponential moving average. Gradient controller is then used
to calculate a proportional response that limits the accepted concurrency.
Concurrency is increased additively when the overload is no longer detected.

:::info

See tutorials on
[Basic Concurrency Limiting](/tutorials/integrations/flow-control/concurrency-limiting/basic-concurrency-limiting.md)
and
[Workload Prioritization](/tutorials/integrations/flow-control/concurrency-limiting/workload-prioritization.md)
to see this blueprint in use.

:::

## Configuration

<!-- Configuration Marker -->

### Common

| Parameter Name       | Parameter Type | Default      | Description         |
| -------------------- | -------------- | ------------ | ------------------- |
| `common.policy_name` | `string`       | `(required)` | Name of the policy. |

### Policy

| Parameter Name       | Parameter Type                  | Default      | Description                            |
| -------------------- | ------------------------------- | ------------ | -------------------------------------- |
| `policy.flux_meter`  | `aperture.spec.v1.FluxMeter`    | `(required)` | Flux Meter.                            |
| `policy.classifiers` | `[]aperture.spec.v1.Classifier` | `[]`         | List of classification rules.          |
| `policy.components`  | `[]aperture.spec.v1.Component`  | `[]`         | List of additional circuit components. |

#### Latency Baseliner

| Parameter Name                                          | Parameter Type                   | Default                                                                                                  | Description                                                                                                                                                                                                                           |
| ------------------------------------------------------- | -------------------------------- | -------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `policy.latency_baseliner.ema`                          | `aperture.spec.v1.EMAParameters` | `{'correction_factor_on_max_envelope_violation': '0.95', 'ema_window': '1500s', 'warmup_window': '60s'}` | EMA parameters.                                                                                                                                                                                                                       |
| `policy.latency_baseliner.latency_tolerance_multiplier` | `float64`                        | `1.1`                                                                                                    | Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms. |
| `policy.latency_baseliner.latency_ema_limit_multiplier` | `float64`                        | `2.0`                                                                                                    | Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.                                                                                                                                    |

#### Concurrency Controller

| Parameter Name                                                        | Parameter Type                          | Default                                                                                                            | Description                                                                                                                                                                                                                                                                                     |
| --------------------------------------------------------------------- | --------------------------------------- | ------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `policy.concurrency_controller.flow_selector`                         | `aperture.spec.v1.FlowSelector`         | `(required)`                                                                                                       | Concurrency Limiter flow selector.                                                                                                                                                                                                                                                              |
| `policy.concurrency_controller.scheduler`                             | `aperture.spec.v1.SchedulerParameters`  | `{'auto_tokens': True, 'default_workload_parameters': {'priority': 20}, 'timeout_factor': '0.5', 'workloads': []}` | Scheduler parameters.                                                                                                                                                                                                                                                                           |
| `policy.concurrency_controller.gradient`                              | `aperture.spec.v1.GradientParameters`   | `{'max_gradient': '1.0', 'min_gradient': '0.1', 'slope': '-1'}`                                                    | Gradient parameters.                                                                                                                                                                                                                                                                            |
| `policy.concurrency_controller.alerter`                               | `aperture.spec.v1.AlerterParameters`    | `{'alert_channels': [], 'alert_name': 'Load Shed Event', 'resolve_timeout': '5s'}`                                 | Whether tokens for workloads are computed dynamically or set statically by the user.                                                                                                                                                                                                            |
| `policy.concurrency_controller.concurrency_limit_multiplier`          | `float64`                               | `2.0`                                                                                                              | Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.                                                                                 |
| `policy.concurrency_controller.concurrency_linear_increment`          | `float64`                               | `5.0`                                                                                                              | Linear increment to concurrency in each execution tick when the system is not in overloaded state.                                                                                                                                                                                              |
| `policy.concurrency_controller.concurrency_sqrt_increment_multiplier` | `float64`                               | `1`                                                                                                                | Scale factor to multiply square root of current accepted concurrrency. This, along with concurrency_linear_increment helps calculate overall concurrency increment in each tick. Concurrency is rapidly ramped up in each execution cycle during normal (non-overload) state (integral effect). |
| `policy.concurrency_controller.dynamic_config`                        | `aperture.v1.LoadActuatorDynamicConfig` | `{'dry_run': False}`                                                                                               | Dynamic configuration for concurrency controller.                                                                                                                                                                                                                                               |

### Dashboard

| Parameter Name               | Parameter Type | Default | Description                            |
| ---------------------------- | -------------- | ------- | -------------------------------------- |
| `dashboard.refresh_interval` | `string`       | `"10s"` | Refresh interval for dashboard panels. |

#### Datasource

| Parameter Name                      | Parameter Type | Default         | Description              |
| ----------------------------------- | -------------- | --------------- | ------------------------ |
| `dashboard.datasource.name`         | `string`       | `"$datasource"` | Datasource name.         |
| `dashboard.datasource.filter_regex` | `string`       | `""`            | Datasource filter regex. |
