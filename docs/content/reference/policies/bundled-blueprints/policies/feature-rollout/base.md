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

See the tutorial on
[Feature Rollout with Average Latency Feedback](/tutorials/flow-control/feature-rollout/with-average-latency-feedback.md)
to see this blueprint in use.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/feature-rollout/base`}>policies/feature-rollout/base</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

#### common {#common}

<!-- vale on -->

<!-- vale off -->

<a id="common-policy-name"></a>

<ParameterDescription
    name='common.policy_name'
    description='Name of the policy.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="policy"></a>

<ParameterDescription
    name='policy'
    description='Parameters for the Feature Rollout policy.'
    type='Object (rollout_policy)'
    reference='#rollout-policy'
    value='{"components": [], "drivers": {}, "evaluation_interval": "1s", "load_shaper": {"flow_regulator_parameters": {"flow_selector": {"flow_matcher": {"control_point": "__REQUIRED_FIELD__"}, "service_selector": {"service": "__REQUIRED_FIELD__"}}, "label_key": ""}, "steps": [{"duration": "__REQUIRED_FIELD__", "target_accept_percentage": "__REQUIRED_FIELD__"}]}, "resources": {"flow_control": {"classifiers": []}}}'
/>

<!-- vale on -->

---

<!-- vale off -->

#### dashboard {#dashboard}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-refresh-interval"></a>

<ParameterDescription
    name='dashboard.refresh_interval'
    description='Refresh interval for dashboard panels.'
    type='string'
    reference=''
    value='"5s"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-from"></a>

<ParameterDescription
    name='dashboard.time_from'
    description='From time of dashboard.'
    type='string'
    reference=''
    value='"now-15m"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-to"></a>

<ParameterDescription
    name='dashboard.time_to'
    description='To time of dashboard.'
    type='string'
    reference=''
    value='"now"'
/>

<!-- vale on -->

<!-- vale off -->

##### dashboard.datasource {#dashboard-datasource}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-name"></a>

<ParameterDescription
    name='dashboard.datasource.name'
    description='Datasource name.'
    type='string'
    reference=''
    value='"$datasource"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-filter-regex"></a>

<ParameterDescription
    name='dashboard.datasource.filter_regex'
    description='Datasource filter regex.'
    type='string'
    reference=''
    value='""'
/>

<!-- vale on -->

---

### Schemas

<!-- vale off -->

#### promql_driver {#promql-driver}

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-query-string"></a>

<ParameterDescription
    name='query_string'
    description='The Prometheus query to be run. Must return a scalar or a vector with a single element.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

##### forward {#promql-driver-forward}

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-forward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the forward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-forward-operator"></a>

<ParameterDescription
    name='operator'
    description='The operator for the forward criteria. oneof: `gt | lt | gte | lte | eq | neq`'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

##### backward {#promql-driver-backward}

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-backward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-backward-operator"></a>

<ParameterDescription
    name='operator'
    description='The operator for the backward criteria. oneof: `gt | lt | gte | lte | eq | neq`'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

##### reset {#promql-driver-reset}

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-reset-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the reset criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-reset-operator"></a>

<ParameterDescription
    name='operator'
    description='The operator for the reset criteria. oneof: `gt | lt | gte | lte | eq | neq`'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### average_latency_driver {#average-latency-driver}

<!-- vale on -->

<!-- vale off -->

<a id="average-latency-driver-flow-selector"></a>

<ParameterDescription
    name='flow_selector'
    description='Identify the service and flows whose latency needs to be measured.'
    type='Object (aperture.spec.v1.FlowSelector)'
    reference='../../../spec#flow-selector'
    value='{"flow_matcher": {"control_point": "__REQUIRED_FIELD__"}, "service_selector": {"service": "__REQUIRED_FIELD__"}}'
/>

<!-- vale on -->

<!-- vale off -->

##### forward {#average-latency-driver-forward}

<!-- vale on -->

<!-- vale off -->

<a id="average-latency-driver-forward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the forward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

##### backward {#average-latency-driver-backward}

<!-- vale on -->

<!-- vale off -->

<a id="average-latency-driver-backward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

##### reset {#average-latency-driver-reset}

<!-- vale on -->

<!-- vale off -->

<a id="average-latency-driver-reset-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the reset criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### percentile_latency_driver {#percentile-latency-driver}

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-flux-meter"></a>

<ParameterDescription
    name='flux_meter'
    description='FluxMeter specifies the flows whose latency needs to be measured and parameters for the histogram metrics.'
    type='Object (aperture.spec.v1.FluxMeter)'
    reference='../../../spec#flux-meter'
    value='{"flow_selector": {"flow_matcher": {"control_point": "__REQUIRED_FIELD__"}, "service_selector": {"service": "__REQUIRED_FIELD__"}}, "static_buckets": {"buckets": [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]}}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-percentile"></a>

<ParameterDescription
    name='percentile'
    description='The percentile to be used for latency measurement.'
    type='Number (double)'
    reference=''
    value='95'
/>

<!-- vale on -->

<!-- vale off -->

##### forward {#percentile-latency-driver-forward}

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-forward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the forward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

##### backward {#percentile-latency-driver-backward}

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-backward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

##### reset {#percentile-latency-driver-reset}

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-reset-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the reset criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### ema_latency_driver {#ema-latency-driver}

<!-- vale on -->

<!-- vale off -->

<a id="ema-latency-driver-flow-selector"></a>

<ParameterDescription
    name='flow_selector'
    description='Identify the service and flows whose latency needs to be measured.'
    type='Object (aperture.spec.v1.FlowSelector)'
    reference='../../../spec#flow-selector'
    value='{"flow_matcher": {"control_point": "__REQUIRED_FIELD__"}, "service_selector": {"service": "__REQUIRED_FIELD__"}}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="ema-latency-driver-ema"></a>

<ParameterDescription
    name='ema'
    description='The parameters for the exponential moving average.'
    type='Object (aperture.spec.v1.EMAParameters)'
    reference='../../../spec#e-m-a-parameters'
    value='{"ema_window": "1500s", "warmup_window": "60s"}'
/>

<!-- vale on -->

<!-- vale off -->

##### forward {#ema-latency-driver-forward}

<!-- vale on -->

<!-- vale off -->

<a id="ema-latency-driver-forward-latency-tolerance-multiplier"></a>

<ParameterDescription
    name='latency_tolerance_multiplier'
    description='The threshold for the forward criteria.'
    type='Number (double)'
    reference=''
    value='1.05'
/>

<!-- vale on -->

<!-- vale off -->

##### backward {#ema-latency-driver-backward}

<!-- vale on -->

<!-- vale off -->

<a id="ema-latency-driver-backward-latency-tolerance-multiplier"></a>

<ParameterDescription
    name='latency_tolerance_multiplier'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='1.05'
/>

<!-- vale on -->

<!-- vale off -->

##### reset {#ema-latency-driver-reset}

<!-- vale on -->

<!-- vale off -->

<a id="ema-latency-driver-reset-latency-tolerance-multiplier"></a>

<ParameterDescription
    name='latency_tolerance_multiplier'
    description='The threshold for the reset criteria.'
    type='Number (double)'
    reference=''
    value='1.25'
/>

<!-- vale on -->

---

<!-- vale off -->

#### rollout_policy {#rollout-policy}

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-load-shaper"></a>

<ParameterDescription
    name='load_shaper'
    description='Identify the service and flows of the feature that needs to be rolled out. And specify feature rollout steps.'
    type='Object (aperture.spec.v1.LoadShaperParameters)'
    reference='../../../spec#load-shaper-parameters'
    value='{"flow_regulator_parameters": {"flow_selector": {"flow_matcher": {"control_point": "__REQUIRED_FIELD__"}, "service_selector": {"service": "__REQUIRED_FIELD__"}}, "label_key": ""}, "steps": [{"duration": "__REQUIRED_FIELD__", "target_accept_percentage": "__REQUIRED_FIELD__"}]}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-components"></a>

<ParameterDescription
    name='components'
    description='List of additional circuit components.'
    type='Array of Object (aperture.spec.v1.Component)'
    reference='../../../spec#component'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-resources"></a>

<ParameterDescription
    name='resources'
    description='List of additional resources.'
    type='Object (aperture.spec.v1.Resources)'
    reference='../../../spec#resources'
    value='{"flow_control": {"classifiers": []}}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-evaluation-interval"></a>

<ParameterDescription
    name='evaluation_interval'
    description='The interval between successive evaluations of the Circuit.'
    type='string'
    reference=''
    value='"1s"'
/>

<!-- vale on -->

<!-- vale off -->

##### drivers {#rollout-policy-drivers}

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-drivers-promql-drivers"></a>

<ParameterDescription
    name='promql_drivers'
    description='List of promql drivers that compare results of a Prometheus query against forward, backward and reset thresholds.'
    type='Array of Object (promql_driver)'
    reference='#promql-driver'
    value='[{"backward": {"operator": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}, "forward": {"operator": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}, "query_string": "__REQUIRED_FIELD__", "reset": {"operator": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-drivers-average-latency-drivers"></a>

<ParameterDescription
    name='average_latency_drivers'
    description='List of drivers that compare average latency against forward, backward and reset thresholds.'
    type='Array of Object (average_latency_driver)'
    reference='#average-latency-driver'
    value='[{"backward": {"threshold": "__REQUIRED_FIELD__"}, "flow_selector": {"flow_matcher": {"control_point": "__REQUIRED_FIELD__"}, "service_selector": {"service": "__REQUIRED_FIELD__"}}, "forward": {"threshold": "__REQUIRED_FIELD__"}, "reset": {"threshold": "__REQUIRED_FIELD__"}}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-drivers-percentile-latency-drivers"></a>

<ParameterDescription
    name='percentile_latency_drivers'
    description='List of drivers that compare percentile latency against forward, backward and reset thresholds.'
    type='Array of Object (percentile_latency_driver)'
    reference='#percentile-latency-driver'
    value='[{"backward": {"threshold": "__REQUIRED_FIELD__"}, "flux_meter": {"flow_selector": {"flow_matcher": {"control_point": "__REQUIRED_FIELD__"}, "service_selector": {"service": "__REQUIRED_FIELD__"}}, "static_buckets": {"buckets": [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]}}, "forward": {"threshold": "__REQUIRED_FIELD__"}, "percentile": 95, "reset": {"threshold": "__REQUIRED_FIELD__"}}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-drivers-ema-latency-drivers"></a>

<ParameterDescription
    name='ema_latency_drivers'
    description='List of drivers that compare trend latency against forward, backward and reset thresholds.'
    type='Array of Object (ema_latency_driver)'
    reference='#ema-latency-driver'
    value='[{"backward": {"latency_tolerance_multiplier": 1.05}, "ema": {"ema_window": "1500s", "warmup_window": "60s"}, "flow_selector": {"flow_matcher": {"control_point": "__REQUIRED_FIELD__"}, "service_selector": {"service": "__REQUIRED_FIELD__"}}, "forward": {"latency_tolerance_multiplier": 1.05}, "reset": {"latency_tolerance_multiplier": 1.25}}]'
/>

<!-- vale on -->

---

## Dynamic Configuration

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

:::

### Parameters

<!-- vale off -->

<a id="load-shaper"></a>

<ParameterDescription
    name='load_shaper'
    description='Default configuration for flow regulator that can be updated at the runtime without shutting down the policy.'
    type='Object (aperture.spec.v1.FlowRegulatorDynamicConfig)'
    reference='../../../spec#flow-regulator-dynamic-config'
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---
