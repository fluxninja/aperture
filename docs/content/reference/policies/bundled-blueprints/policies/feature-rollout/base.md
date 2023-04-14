---
title: Feature Rollout Policy
---

## Introduction

This policy rolls out new features based on closed loop feedback. The rollout
criteria are defined by drivers that determine conditions for advancing,
reversing, or resetting the rollout to its initial state. The rollout process
consists of a series of steps that progress if the feature is considered
healthy.

:::info

Please see reference for the
[`LoadShaper`](/reference/policies/spec.md#load-shaper) component that is used
within this blueprint.

:::

:::info

See tutorials on TODO tgill to see this blueprint in use.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../../parameterComponents.js'
```

## Configuration

Code: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/feature-rollout/base`}>policies/feature-rollout/base</a>

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

<a id="policy"></a> <ParameterDescription
    name="policy"
    type="
Object (rollout_policy)"
    reference="#rollout-policy"
    value="{'components': [], 'drivers': {}, 'evaluation_interval': '1s', 'load_shaper': {'flow_regulator_parameters': {'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'label_key': ''}, 'steps': [{'duration': '__REQUIRED_FIELD__', 'target_accept_percentage': '__REQUIRED_FIELD__'}]}, 'resources': {'flow_control': {'classifiers': []}}}"
    description='Parameters for the Feature Rollout policy.' />

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

### Schemas

#### promql_driver {#promql-driver}

<a id="promql-driver-query-string"></a> <ParameterDescription
    name="query_string"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The Prometheus query to be run. Must return a scalar or a vector with a single element.' />

##### forward {#promql-driver-forward}

<a id="promql-driver-forward-threshold"></a> <ParameterDescription
    name="threshold"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The threshold for the forward criteria.' />

<a id="promql-driver-forward-operator"></a> <ParameterDescription
    name="operator"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The operator for the forward criteria. oneof: `gt | lt | gte | lte | eq | neq`' />

##### backward {#promql-driver-backward}

<a id="promql-driver-backward-threshold"></a> <ParameterDescription
    name="threshold"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The threshold for the backward criteria.' />

<a id="promql-driver-backward-operator"></a> <ParameterDescription
    name="operator"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The operator for the backward criteria. oneof: `gt | lt | gte | lte | eq | neq`' />

##### reset {#promql-driver-reset}

<a id="promql-driver-reset-threshold"></a> <ParameterDescription
    name="threshold"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The threshold for the reset criteria.' />

<a id="promql-driver-reset-operator"></a> <ParameterDescription
    name="operator"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The operator for the reset criteria. oneof: `gt | lt | gte | lte | eq | neq`' />

---

#### average_latency_driver {#average-latency-driver}

<a id="average-latency-driver-flow-selector"></a> <ParameterDescription
    name="flow_selector"
    type="
Object (aperture.spec.v1.FlowSelector)"
    reference="../../../spec#flow-selector"
    value="{'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}"
    description='Identify the service and flows whose latency needs to be measured.' />

##### forward {#average-latency-driver-forward}

<a id="average-latency-driver-forward-threshold"></a> <ParameterDescription
    name="threshold"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The threshold for the forward criteria.' />

##### backward {#average-latency-driver-backward}

<a id="average-latency-driver-backward-threshold"></a> <ParameterDescription
    name="threshold"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The threshold for the backward criteria.' />

##### reset {#average-latency-driver-reset}

<a id="average-latency-driver-reset-threshold"></a> <ParameterDescription
    name="threshold"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The threshold for the reset criteria.' />

---

#### percentile_latency_driver {#percentile-latency-driver}

<a id="percentile-latency-driver-flux-meter"></a> <ParameterDescription
    name="flux_meter"
    type="
Object (aperture.spec.v1.FluxMeter)"
    reference="../../../spec#flux-meter"
    value="{'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'static_buckets': {'buckets': [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]}}"
    description='FluxMeter specifies the flows whose latency needs to be measured and parameters for the histogram metrics.' />

<a id="percentile-latency-driver-percentile"></a> <ParameterDescription
    name="percentile"
    type="
Number (double)"
    reference=""
    value="95"
    description='The percentile to be used for latency measurement.' />

##### forward {#percentile-latency-driver-forward}

<a id="percentile-latency-driver-forward-threshold"></a> <ParameterDescription
    name="threshold"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The threshold for the forward criteria.' />

##### backward {#percentile-latency-driver-backward}

<a id="percentile-latency-driver-backward-threshold"></a> <ParameterDescription
    name="threshold"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The threshold for the backward criteria.' />

##### reset {#percentile-latency-driver-reset}

<a id="percentile-latency-driver-reset-threshold"></a> <ParameterDescription
    name="threshold"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='The threshold for the reset criteria.' />

---

#### ema_latency_driver {#ema-latency-driver}

<a id="ema-latency-driver-flow-selector"></a> <ParameterDescription
    name="flow_selector"
    type="
Object (aperture.spec.v1.FlowSelector)"
    reference="../../../spec#flow-selector"
    value="{'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}"
    description='Identify the service and flows whose latency needs to be measured.' />

<a id="ema-latency-driver-ema"></a> <ParameterDescription
    name="ema"
    type="
Object (aperture.spec.v1.EMAParameters)"
    reference="../../../spec#e-m-a-parameters"
    value="{'ema_window': '1500s', 'warmup_window': '60s'}"
    description='The parameters for the exponential moving average.' />

##### forward {#ema-latency-driver-forward}

<a id="ema-latency-driver-forward-latency-tolerance-multiplier"></a>
<ParameterDescription
    name="latency_tolerance_multiplier"
    type="
Number (double)"
    reference=""
    value="1.05"
    description='The threshold for the forward criteria.' />

##### backward {#ema-latency-driver-backward}

<a id="ema-latency-driver-backward-latency-tolerance-multiplier"></a>
<ParameterDescription
    name="latency_tolerance_multiplier"
    type="
Number (double)"
    reference=""
    value="1.05"
    description='The threshold for the backward criteria.' />

##### reset {#ema-latency-driver-reset}

<a id="ema-latency-driver-reset-latency-tolerance-multiplier"></a>
<ParameterDescription
    name="latency_tolerance_multiplier"
    type="
Number (double)"
    reference=""
    value="1.25"
    description='The threshold for the reset criteria.' />

---

#### rollout_policy {#rollout-policy}

<a id="rollout-policy-load-shaper"></a> <ParameterDescription
    name="load_shaper"
    type="
Object (aperture.spec.v1.LoadShaperParameters)"
    reference="../../../spec#load-shaper-parameters"
    value="{'flow_regulator_parameters': {'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'label_key': ''}, 'steps': [{'duration': '__REQUIRED_FIELD__', 'target_accept_percentage': '__REQUIRED_FIELD__'}]}"
    description='Identify the service and flows of the feature that needs to be rolled out. And specify feature rollout steps.' />

<a id="rollout-policy-components"></a> <ParameterDescription
    name="components"
    type="
Array of
Object (aperture.spec.v1.Component)"
    reference="../../../spec#component"
    value="[]"
    description='List of additional circuit components.' />

<a id="rollout-policy-resources"></a> <ParameterDescription
    name="resources"
    type="
Object (aperture.spec.v1.Resources)"
    reference="../../../spec#resources"
    value="{'flow_control': {'classifiers': []}}"
    description='List of additional resources.' />

<a id="rollout-policy-evaluation-interval"></a> <ParameterDescription
    name="evaluation_interval"
    type="
string"
    reference=""
    value="'1s'"
    description='The interval between successive evaluations of the Circuit.' />

##### drivers {#rollout-policy-drivers}

<a id="rollout-policy-drivers-promql-drivers"></a> <ParameterDescription
    name="promql_drivers"
    type="
Array of
Object (promql_driver)"
    reference="#promql-driver"
    value="[{'backward': {'operator': '__REQUIRED_FIELD__', 'threshold': '__REQUIRED_FIELD__'}, 'forward': {'operator': '__REQUIRED_FIELD__', 'threshold': '__REQUIRED_FIELD__'}, 'query_string': '__REQUIRED_FIELD__', 'reset': {'operator': '__REQUIRED_FIELD__', 'threshold': '__REQUIRED_FIELD__'}}]"
    description='List of promql drivers that compare results of a Prometheus query against forward, backward and reset thresholds.' />

<a id="rollout-policy-drivers-average-latency-drivers"></a>
<ParameterDescription
    name="average_latency_drivers"
    type="
Array of
Object (average_latency_driver)"
    reference="#average-latency-driver"
    value="[{'backward': {'threshold': '__REQUIRED_FIELD__'}, 'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'forward': {'threshold': '__REQUIRED_FIELD__'}, 'reset': {'threshold': '__REQUIRED_FIELD__'}}]"
    description='List of drivers that compare average latency against forward, backward and reset thresholds.' />

<a id="rollout-policy-drivers-percentile-latency-drivers"></a>
<ParameterDescription
    name="percentile_latency_drivers"
    type="
Array of
Object (percentile_latency_driver)"
    reference="#percentile-latency-driver"
    value="[{'backward': {'threshold': '__REQUIRED_FIELD__'}, 'flux_meter': {'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'static_buckets': {'buckets': [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]}}, 'forward': {'threshold': '__REQUIRED_FIELD__'}, 'percentile': 95, 'reset': {'threshold': '__REQUIRED_FIELD__'}}]"
    description='List of drivers that compare percentile latency against forward, backward and reset thresholds.' />

<a id="rollout-policy-drivers-ema-latency-drivers"></a> <ParameterDescription
    name="ema_latency_drivers"
    type="
Array of
Object (ema_latency_driver)"
    reference="#ema-latency-driver"
    value="[{'backward': {'latency_tolerance_multiplier': 1.05}, 'ema': {'ema_window': '1500s', 'warmup_window': '60s'}, 'flow_selector': {'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}, 'forward': {'latency_tolerance_multiplier': 1.05}, 'reset': {'latency_tolerance_multiplier': 1.25}}]"
    description='List of drivers that compare trend latency against forward, backward and reset thresholds.' />

---

## Dynamic Configuration

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

:::

### Parameters

<a id="load-shaper"></a> <ParameterDescription
    name="load_shaper"
    type="
Object (aperture.spec.v1.FlowRegulatorDynamicConfig)"
    reference="../../../spec#flow-regulator-dynamic-config"
    value="__REQUIRED_FIELD__"
    description='Default configuration for flow regulator that can be updated at the runtime without shutting down the policy.' />

---
