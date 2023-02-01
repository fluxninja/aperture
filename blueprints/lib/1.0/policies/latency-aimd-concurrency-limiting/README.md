# Latency AIMD Concurrency Limiting Policy

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

**`common.policy_name`** (type: _`string`_)

required parameter

Name of the policy.

### Policy

**`policy.flux_meter`** (type: _`aperture.spec.v1.FluxMeter`_)

required parameter

Flux Meter.

**`policy.classifiers`** (type: _`[]aperture.spec.v1.Classifier`_)

default: `[]`

List of classification rules.

**`policy.components`** (type: _`[]aperture.spec.v1.Component`_)

default: `[]`

List of additional circuit components.

#### Latency Baseliner

**`policy.latency_baseliner.ema`** (type: _`aperture.spec.v1.EMAParameters`_)

default:
`{'correction_factor_on_max_envelope_violation': '0.95', 'ema_window': '1500s', 'warmup_window': '60s'}`

EMA parameters.

**`policy.latency_baseliner.latency_tolerance_multiplier`** (type: _`float64`_)

default: `1.1`

Tolerance factor beyond which the service is considered to be in overloaded
state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is
considered to be in overloaded state if current latency is more than 55ms.

**`policy.latency_baseliner.latency_ema_limit_multiplier`** (type: _`float64`_)

default: `2.0`

Current latency value is multiplied with this factor to calculate maximum
envelope of Latency EMA.

#### Concurrency Controller

**`policy.concurrency_controller.flow_selector`** (type:
_`aperture.spec.v1.FlowSelector`_)

required parameter

Concurrency Limiter flow selector.

**`policy.concurrency_controller.scheduler`** (type:
_`aperture.spec.v1.SchedulerParameters`_)

default:
`{'auto_tokens': True, 'default_workload_parameters': {'priority': 20}, 'timeout_factor': '0.5', 'workloads': []}`

Scheduler parameters.

**`policy.concurrency_controller.gradient`** (type:
_`aperture.spec.v1.GradientParameters`_)

default: `{'max_gradient': '1.0', 'min_gradient': '0.1', 'slope': '-1'}`

Gradient parameters.

**`policy.concurrency_controller.alerter`** (type:
_`aperture.spec.v1.AlerterParameters`_)

default:
`{'alert_channels': [], 'alert_name': 'Load Shed Event', 'resolve_timeout': '5s'}`

Whether tokens for workloads are computed dynamically or set statically by the
user.

**`policy.concurrency_controller.concurrency_limit_multiplier`** (type:
_`float64`_)

default: `2.0`

Current accepted concurrency is multiplied with this number to dynamically
calculate the upper concurrency limit of a Service during normal (non-overload)
state. This protects the Service from sudden spikes.

**`policy.concurrency_controller.concurrency_linear_increment`** (type:
_`float64`_)

default: `5.0`

Linear increment to concurrency in each execution tick when the system is not in
overloaded state.

**`policy.concurrency_controller.concurrency_sqrt_increment_multiplier`** (type:
_`float64`_)

default: `1`

Scale factor to multiply square root of current accepted concurrrency. This,
along with concurrency_linear_increment helps calculate overall concurrency
increment in each tick. Concurrency is rapidly ramped up in each execution cycle
during normal (non-overload) state (integral effect).

**`policy.concurrency_controller.dynamic_config`** (type:
_`aperture.v1.LoadActuatorDynamicConfig`_)

default: `{'dry_run': False}`

Dynamic configuration for concurrency controller.

### Dashboard

**`dashboard.refresh_interval`** (type: _`string`_)

default: `"10s"`

Refresh interval for dashboard panels.

#### Datasource

**`dashboard.datasource.name`** (type: _`string`_)

default: `"$datasource"`

Datasource name.

**`dashboard.datasource.filter_regex`** (type: _`string`_)

default: `""`

Datasource filter regex.
