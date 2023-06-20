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
[`LoadRamp`](/reference/configuration/spec.md#load-ramp) component that is used
within this blueprint.

:::

:::info

See the tutorial on
[Feature Rollout with Average Latency Feedback](/use-cases/feature-rollout/with-average-latency-feedback.md)
to see this blueprint in use.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/feature-rollout/base`}>policies/feature-rollout/base</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

<a id="policy"></a>

<ParameterDescription
    name='policy'
    description='Parameters for the Feature Rollout policy.'
    type='Object (rollout_policy)'
    reference='#rollout-policy'
    value='{"components": [], "drivers": {}, "evaluation_interval": "1s", "load_ramp": {"sampler": {"label_key": "", "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "steps": [{"duration": "__REQUIRED_FIELD__", "target_accept_percentage": "__REQUIRED_FIELD__"}]}, "policy_name": "__REQUIRED_FIELD__", "resources": {"flow_control": {"classifiers": []}}, "rollout": false}'
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

<a id="dashboard-extra-filters"></a>

<ParameterDescription
    name='dashboard.extra_filters'
    description='Additional filters to pass to each query to Grafana datasource.'
    type='Object (map[string]string)'
    reference='#map-string-string'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-title"></a>

<ParameterDescription
    name='dashboard.title'
    description='Name of the main dashboard.'
    type='string'
    reference=''
    value='"Aperture Feature Rollout"'
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

##### criteria {#promql-driver-criteria}

<!-- vale on -->

<!-- vale off -->

###### forward {#promql-driver-criteria-forward}

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-criteria-forward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the forward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-criteria-forward-operator"></a>

<ParameterDescription
    name='operator'
    description='The operator for the forward criteria. oneof: `gt | lt | gte | lte | eq | neq`'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

###### backward {#promql-driver-criteria-backward}

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-criteria-backward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-criteria-backward-operator"></a>

<ParameterDescription
    name='operator'
    description='The operator for the backward criteria. oneof: `gt | lt | gte | lte | eq | neq`'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

###### reset {#promql-driver-criteria-reset}

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-criteria-reset-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the reset criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-criteria-reset-operator"></a>

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

<a id="average-latency-driver-selectors"></a>

<ParameterDescription
    name='selectors'
    description='Identify the service and flows whose latency needs to be measured.'
    type='Array of Object (aperture.spec.v1.Selector)'
    reference='../../../spec#selector'
    value='[{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

<!-- vale off -->

##### criteria {#average-latency-driver-criteria}

<!-- vale on -->

<!-- vale off -->

###### forward {#average-latency-driver-criteria-forward}

<!-- vale on -->

<!-- vale off -->

<a id="average-latency-driver-criteria-forward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the forward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

###### backward {#average-latency-driver-criteria-backward}

<!-- vale on -->

<!-- vale off -->

<a id="average-latency-driver-criteria-backward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

###### reset {#average-latency-driver-criteria-reset}

<!-- vale on -->

<!-- vale off -->

<a id="average-latency-driver-criteria-reset-threshold"></a>

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
    value='{"selector": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}], "static_buckets": {"buckets": [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]}}'
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

##### criteria {#percentile-latency-driver-criteria}

<!-- vale on -->

<!-- vale off -->

###### forward {#percentile-latency-driver-criteria-forward}

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-criteria-forward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the forward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

###### backward {#percentile-latency-driver-criteria-backward}

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-criteria-backward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

###### reset {#percentile-latency-driver-criteria-reset}

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-criteria-reset-threshold"></a>

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

<a id="ema-latency-driver-selectors"></a>

<ParameterDescription
    name='selectors'
    description='Identify the service and flows whose latency needs to be measured.'
    type='Array of Object (aperture.spec.v1.Selector)'
    reference='../../../spec#selector'
    value='[{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]'
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

##### criteria {#ema-latency-driver-criteria}

<!-- vale on -->

<!-- vale off -->

###### forward {#ema-latency-driver-criteria-forward}

<!-- vale on -->

<!-- vale off -->

<a id="ema-latency-driver-criteria-forward-latency-tolerance-multiplier"></a>

<ParameterDescription
    name='latency_tolerance_multiplier'
    description='The threshold for the forward criteria.'
    type='Number (double)'
    reference=''
    value='1.05'
/>

<!-- vale on -->

<!-- vale off -->

###### backward {#ema-latency-driver-criteria-backward}

<!-- vale on -->

<!-- vale off -->

<a id="ema-latency-driver-criteria-backward-latency-tolerance-multiplier"></a>

<ParameterDescription
    name='latency_tolerance_multiplier'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='1.05'
/>

<!-- vale on -->

<!-- vale off -->

###### reset {#ema-latency-driver-criteria-reset}

<!-- vale on -->

<!-- vale off -->

<a id="ema-latency-driver-criteria-reset-latency-tolerance-multiplier"></a>

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

<a id="rollout-policy-policy-name"></a>

<ParameterDescription
    name='policy_name'
    description='Name of the policy.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-load-ramp"></a>

<ParameterDescription
    name='load_ramp'
    description='Identify the service and flows of the feature that needs to be rolled out. And specify feature rollout steps.'
    type='Object (aperture.spec.v1.LoadRampParameters)'
    reference='../../../spec#load-ramp-parameters'
    value='{"sampler": {"label_key": "", "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "steps": [{"duration": "__REQUIRED_FIELD__", "target_accept_percentage": "__REQUIRED_FIELD__"}]}'
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

<a id="rollout-policy-rollout"></a>

<ParameterDescription
    name='rollout'
    description='Whether to start the rollout. This setting may be overridden at runtime via dynamic configuration.'
    type='Boolean'
    reference=''
    value='false'
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
    value='[{"criteria": {"backward": {"operator": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}, "forward": {"operator": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}, "reset": {"operator": "__REQUIRED_FIELD__", "threshold": "__REQUIRED_FIELD__"}}, "query_string": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-drivers-average-latency-drivers"></a>

<ParameterDescription
    name='average_latency_drivers'
    description='List of drivers that compare average latency against forward, backward and reset thresholds.'
    type='Array of Object (average_latency_driver)'
    reference='#average-latency-driver'
    value='[{"criteria": {"backward": {"threshold": "__REQUIRED_FIELD__"}, "forward": {"threshold": "__REQUIRED_FIELD__"}, "reset": {"threshold": "__REQUIRED_FIELD__"}}, "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-drivers-percentile-latency-drivers"></a>

<ParameterDescription
    name='percentile_latency_drivers'
    description='List of drivers that compare percentile latency against forward, backward and reset thresholds.'
    type='Array of Object (percentile_latency_driver)'
    reference='#percentile-latency-driver'
    value='[{"criteria": {"backward": {"threshold": "__REQUIRED_FIELD__"}, "forward": {"threshold": "__REQUIRED_FIELD__"}, "reset": {"threshold": "__REQUIRED_FIELD__"}}, "flux_meter": {"selector": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}], "static_buckets": {"buckets": [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]}}, "percentile": 95}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="rollout-policy-drivers-ema-latency-drivers"></a>

<ParameterDescription
    name='ema_latency_drivers'
    description='List of drivers that compare trend latency against forward, backward and reset thresholds.'
    type='Array of Object (ema_latency_driver)'
    reference='#ema-latency-driver'
    value='[{"criteria": {"backward": {"latency_tolerance_multiplier": 1.05}, "forward": {"latency_tolerance_multiplier": 1.05}, "reset": {"latency_tolerance_multiplier": 1.25}}, "ema": {"ema_window": "1500s", "warmup_window": "60s"}, "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}]'
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

<a id="pass-through-label-values"></a>

<ParameterDescription
    name='pass_through_label_values'
    description='Specify certain label values to be always accepted by the _Sampler_ regardless of accept percentage. This configuration can be updated at the runtime without shutting down the policy.'
    type='Array of string'
    reference=''
    value='["__REQUIRED_FIELD__"]'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="rollout"></a>

<ParameterDescription
    name='rollout'
    description='Start feature rollout. This setting can be updated at runtime without shutting down the policy. The feature rollout gets paused if this flag is set to false in the middle of a feature rollout.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="reset"></a>

<ParameterDescription
    name='reset'
    description='Reset feature rollout to the first step. This setting can be updated at the runtime without shutting down the policy.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

---
