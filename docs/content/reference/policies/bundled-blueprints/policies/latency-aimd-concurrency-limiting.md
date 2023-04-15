---
title: Latency AIMD Concurrency Limiting Policy
---

## Introduction

This policy detects traffic overloads and cascading failure build-up by
comparing the real-time latency with its exponential moving average. A gradient
controller calculates a proportional response to limit accepted concurrency,
which is increased additively when the overload is no longer detected.

:::info

AIMD stands for Additive Increase, Multiplicative Decrease. The concurrency is
reduced by a multiplicative factor when the service is overloaded, and increased
by an additive factor while the service is no longer overloaded.

Please see reference for the
[`AIMDConcurrencyController`](/reference/policies/spec.md#a-i-m-d-concurrency-controller)
component that is used within this blueprint.

:::

:::info

See tutorials on
[Basic Concurrency Limiting](/tutorials/flow-control/concurrency-limiting/basic-concurrency-limiting.md)
and
[Workload Prioritization](/tutorials/flow-control/concurrency-limiting/workload-prioritization.md)
to see this blueprint in use.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Code: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/latency-aimd-concurrency-limiting`}>policies/latency-aimd-concurrency-limiting</a>

<!-- vale on -->

### Parameters

#### common {#common}

<a id="common-policy-name"></a> <ParameterDescription
    name="common.policy_name"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

---

#### policy {#policy}

<a id="policy-flux-meter"></a> <ParameterDescription
    name="policy.flux_meter"
    type="
Object (aperture.spec.v1.FluxMeter)"
    reference="../../spec#flux-meter"
    value="{'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}}"
    description='Flux Meter.' />

<a id="policy-classifiers"></a> <ParameterDescription
    name="policy.classifiers"
    type="
Array of
Object (aperture.spec.v1.Classifier)"
    reference="../../spec#classifier"
    value="[]"
    description='List of classification rules.' />

<a id="policy-components"></a> <ParameterDescription
    name="policy.components"
    type="
Array of
Object (aperture.spec.v1.Component)"
    reference="../../spec#component"
    value="[]"
    description='List of additional circuit components.' />

##### policy.latency_baseliner {#policy-latency-baseliner}

<a id="policy-latency-baseliner-ema"></a> <ParameterDescription
    name="policy.latency_baseliner.ema"
    type="
Object (aperture.spec.v1.EMAParameters)"
    reference="../../spec#e-m-a-parameters"
    value="{'correction_factor_on_max_envelope_violation': 0.95, 'ema_window': '1500s', 'warmup_window': '60s'}"
    description='EMA parameters.' />

<a id="policy-latency-baseliner-latency-tolerance-multiplier"></a>
<ParameterDescription
    name="policy.latency_baseliner.latency_tolerance_multiplier"
    type="
Number (double)"
    reference=""
    value="1.1"
    description='Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.' />

<a id="policy-latency-baseliner-latency-ema-limit-multiplier"></a>
<ParameterDescription
    name="policy.latency_baseliner.latency_ema_limit_multiplier"
    type="
Number (double)"
    reference=""
    value="2"
    description='Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.' />

##### policy.concurrency_controller {#policy-concurrency-controller}

<a id="policy-concurrency-controller-flow-selector"></a> <ParameterDescription
    name="policy.concurrency_controller.flow_selector"
    type="
Object (aperture.spec.v1.FlowSelector)"
    reference="../../spec#flow-selector"
    value="{'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}"
    description='Concurrency Limiter flow selector.' />

<a id="policy-concurrency-controller-scheduler"></a> <ParameterDescription
    name="policy.concurrency_controller.scheduler"
    type="
Object (aperture.spec.v1.SchedulerParameters)"
    reference="../../spec#scheduler-parameters"
    value="{'auto_tokens': True}"
    description='Scheduler parameters.' />

<a id="policy-concurrency-controller-gradient"></a> <ParameterDescription
    name="policy.concurrency_controller.gradient"
    type="
Object (aperture.spec.v1.GradientControllerParameters)"
    reference="../../spec#gradient-controller-parameters"
    value="{'max_gradient': 1, 'min_gradient': 0.1, 'slope': -1}"
    description='Gradient Controller parameters.' />

<a id="policy-concurrency-controller-alerter"></a> <ParameterDescription
    name="policy.concurrency_controller.alerter"
    type="
Object (aperture.spec.v1.AlerterParameters)"
    reference="../../spec#alerter-parameters"
    value="{'alert_name': 'Load Shed Event'}"
    description='Whether tokens for workloads are computed dynamically or set statically by the user.' />

<a id="policy-concurrency-controller-max-load-multiplier"></a>
<ParameterDescription
    name="policy.concurrency_controller.max_load_multiplier"
    type="
Number (double)"
    reference=""
    value="2"
    description='Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.' />

<a id="policy-concurrency-controller-load-multiplier-linear-increment"></a>
<ParameterDescription
    name="policy.concurrency_controller.load_multiplier_linear_increment"
    type="
Number (double)"
    reference=""
    value="0.0025"
    description='Linear increment to load multiplier in each execution tick (0.5s) when the system is not in overloaded state.' />

<a id="policy-concurrency-controller-default-config"></a> <ParameterDescription
    name="policy.concurrency_controller.default_config"
    type="
Object (aperture.spec.v1.LoadActuatorDynamicConfig)"
    reference="../../spec#load-actuator-dynamic-config"
    value="{'dry_run': False}"
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.' />

---

#### dashboard {#dashboard}

<a id="dashboard-refresh-interval"></a> <ParameterDescription
    name="dashboard.refresh_interval"
    type="
string"
    reference=""
    value="'5s'"
    description='Refresh interval for dashboard panels.' />

<a id="dashboard-time-from"></a> <ParameterDescription
    name="dashboard.time_from"
    type="
string"
    reference=""
    value="'now-15m'"
    description='From time of dashboard.' />

<a id="dashboard-time-to"></a> <ParameterDescription
    name="dashboard.time_to"
    type="
string"
    reference=""
    value="'now'"
    description='To time of dashboard.' />

##### dashboard.datasource {#dashboard-datasource}

<a id="dashboard-datasource-name"></a> <ParameterDescription
    name="dashboard.datasource.name"
    type="
string"
    reference=""
    value="'$datasource'"
    description='Datasource name.' />

<a id="dashboard-datasource-filter-regex"></a> <ParameterDescription
    name="dashboard.datasource.filter_regex"
    type="
string"
    reference=""
    value="''"
    description='Datasource filter regex.' />

---

## Dynamic Configuration

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

:::

### Parameters

<a id="concurrency-controller"></a> <ParameterDescription
    name="concurrency_controller"
    type="
Object (aperture.spec.v1.LoadActuatorDynamicConfig)"
    reference="../../spec#load-actuator-dynamic-config"
    value="__REQUIRED_FIELD__"
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.' />

---
