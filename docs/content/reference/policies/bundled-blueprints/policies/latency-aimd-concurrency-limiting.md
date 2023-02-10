---
title: Latency AIMD Concurrency Limiting Policy
---

## Introduction

This policy detects overloads/cascading failures by comparing the real-time
latency with it's exponential moving average. Gradient controller is then used
to calculate a proportional response that limits the accepted concurrency.
Concurrency is increased additively when the overload is no longer detected.

:::info

AIMD stands for Additive Increase, Multiplicative Decrease. That is, the
concurrency is reduced by a multiplicative factor when the service is overloaded
and increased by an additive factor when the service is no longer overloaded.

:::

:::info

See tutorials on
[Basic Concurrency Limiting](/tutorials/integrations/flow-control/concurrency-limiting/basic-concurrency-limiting.md)
and
[Workload Prioritization](/tutorials/integrations/flow-control/concurrency-limiting/workload-prioritization.md)
to see this blueprint in use.

:::

## Configuration

<!-- Configuration Marker -->

```mdx-code-block

export const ParameterHeading = ({children}) => (
  <span style={{fontWeight: "bold"}}>{children}</span>
);

export const WrappedDescription = ({children}) => (
  <span style={{wordWrap: "normal"}}>{children}</span>
);

export const RefType = ({type, reference}) => (
  <a href={reference}>{type}</a>
);

export const ParameterDescription = ({name, type, reference, value, description}) => (
  <table class="blueprints-params">
  <tr>
    <td><ParameterHeading>Parameter</ParameterHeading></td>
    <td><code>{name}</code></td>
  </tr>
  <tr>
    <td><ParameterHeading>Type</ParameterHeading></td>
    <td><em>{reference == "" ? type : <RefType type={type} reference={reference} />}</em></td>
  </tr>
  <tr>
    <td class="blueprints-default-heading"><ParameterHeading>Default Value</ParameterHeading></td>
    <td><code>{value}</code></td>
  </tr>
  <tr>
    <td class="blueprints-description"><ParameterHeading>Description</ParameterHeading></td>
    <td class="blueprints-description"><WrappedDescription>{description}</WrappedDescription></td>
  </tr>
</table>
);
```

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
```

Code: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/latency-aimd-concurrency-limiting`}>policies/latency-aimd-concurrency-limiting</a>

<h3 class="blueprints-h3">Common</h3>

<ParameterDescription
    name="common.policy_name"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

<h3 class="blueprints-h3">Policy</h3>

<ParameterDescription
    name="policy.flux_meter"
    type="aperture.spec.v1.FluxMeter"
    reference="../../spec#v1-flux-meter"
    value="{'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'agent_group': 'default', 'service': '__REQUIRED_FIELD__'}}}"
    description='Flux Meter.' />

<ParameterDescription
    name="policy.flux_meter.flow_selector.service_selector.service"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Service Name.' />

<ParameterDescription
    name="policy.flux_meter.flow_selector.flow_matcher.control_point"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Control Point Name.' />

<ParameterDescription
    name="policy.classifiers"
    type="[]aperture.spec.v1.Classifier"
    reference="../../spec#v1-classifier"
    value="[]"
    description='List of classification rules.' />

<ParameterDescription
    name="policy.components"
    type="[]aperture.spec.v1.Component"
    reference="../../spec#v1-component"
    value="[]"
    description='List of additional circuit components.' />

<h4 class="blueprints-h4">Latency Baseliner</h4>

<ParameterDescription
    name="policy.latency_baseliner.ema"
    type="aperture.spec.v1.EMAParameters"
    reference="../../spec#v1-e-m-a-parameters"
    value="{'correction_factor_on_max_envelope_violation': '0.95', 'ema_window': '1500s', 'warmup_window': '60s'}"
    description='EMA parameters.' />

<ParameterDescription
    name="policy.latency_baseliner.latency_tolerance_multiplier"
    type="float64"
    reference=""
    value="1.1"
    description='Tolerance factor beyond which the service is considered to be in overloaded state. E.g. if EMA of latency is 50ms and if Tolerance is 1.1, then service is considered to be in overloaded state if current latency is more than 55ms.' />

<ParameterDescription
    name="policy.latency_baseliner.latency_ema_limit_multiplier"
    type="float64"
    reference=""
    value="2.0"
    description='Current latency value is multiplied with this factor to calculate maximum envelope of Latency EMA.' />

<h4 class="blueprints-h4">Concurrency Controller</h4>

<ParameterDescription
    name="policy.concurrency_controller.flow_selector"
    type="aperture.spec.v1.FlowSelector"
    reference="../../spec#v1-flow-selector"
    value="{'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'agent_group': 'default', 'service': '__REQUIRED_FIELD__'}}"
    description='Concurrency Limiter flow selector.' />

<ParameterDescription
    name="policy.concurrency_controller.flow_selector.service_selector.service"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Service Name.' />

<ParameterDescription
    name="policy.concurrency_controller.flow_selector.flow_matcher.control_point"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Control Point Name.' />

<ParameterDescription
    name="policy.concurrency_controller.scheduler"
    type="aperture.spec.v1.SchedulerParameters"
    reference="../../spec#v1-scheduler-parameters"
    value="{'auto_tokens': True, 'default_workload_parameters': {'priority': 20}, 'timeout_factor': '0.5', 'workloads': []}"
    description='Scheduler parameters.' />

<ParameterDescription
    name="policy.concurrency_controller.gradient"
    type="aperture.spec.v1.GradientControllerParameters"
    reference="../../spec#v1-gradient-controller-parameters"
    value="{'max_gradient': '1.0', 'min_gradient': '0.1', 'slope': '-1'}"
    description='Gradient Controller parameters.' />

<ParameterDescription
    name="policy.concurrency_controller.alerter"
    type="aperture.spec.v1.AlerterParameters"
    reference="../../spec#v1-alerter-parameters"
    value="{'alert_channels': [], 'alert_name': 'Load Shed Event', 'resolve_timeout': '5s'}"
    description='Whether tokens for workloads are computed dynamically or set statically by the user.' />

<ParameterDescription
    name="policy.concurrency_controller.concurrency_limit_multiplier"
    type="float64"
    reference=""
    value="2.0"
    description='Current accepted concurrency is multiplied with this number to dynamically calculate the upper concurrency limit of a Service during normal (non-overload) state. This protects the Service from sudden spikes.' />

<ParameterDescription
    name="policy.concurrency_controller.concurrency_linear_increment"
    type="float64"
    reference=""
    value="5.0"
    description='Linear increment to concurrency in each execution tick when the system is not in overloaded state.' />

<ParameterDescription
    name="policy.concurrency_controller.concurrency_sqrt_increment_multiplier"
    type="float64"
    reference=""
    value="1"
    description='Scale factor to multiply square root of current accepted concurrrency. This, along with concurrency_linear_increment helps calculate overall concurrency increment in each tick. Concurrency is rapidly ramped up in each execution cycle during normal (non-overload) state (integral effect).' />

<ParameterDescription
    name="policy.concurrency_controller.default_config"
    type="aperture.v1.LoadActuatorDynamicConfig"
    reference=""
    value="{'dry_run': False}"
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.' />

<h3 class="blueprints-h3">Dashboard</h3>

<ParameterDescription
    name="dashboard.refresh_interval"
    type="string"
    reference=""
    value="'10s'"
    description='Refresh interval for dashboard panels.' />

<h4 class="blueprints-h4">Datasource</h4>

<ParameterDescription
    name="dashboard.datasource.name"
    type="string"
    reference=""
    value="'$datasource'"
    description='Datasource name.' />

<ParameterDescription
    name="dashboard.datasource.filter_regex"
    type="string"
    reference=""
    value="''"
    description='Datasource filter regex.' />## Dynamic Configuration The
following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

<h3 class="blueprints-h3">Dynamic Configuration</h3>

<ParameterDescription
    name="concurrency_controller"
    type="aperture.v1.LoadActuatorDynamicConfig"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Default configuration for concurrency controller that can be updated at the runtime without shutting down the policy.' />
