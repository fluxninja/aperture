---
title: Feature Rollout with Percentile Latency Feedback
---

## Introduction

This policy rolls out new features based on percentile of latency as the rollout
criteria. The percentile of latency is compared with thresholds to determine
conditions for advancing, reversing, or resetting the rollout to its initial
state. The rollout process consists of a series of steps that progress if the
feature is considered healthy.

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/feature-rollout/percentile-latency`}>policies/feature-rollout/percentile-latency</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

#### policy {#policy}

<!-- vale on -->

<!-- vale off -->

<a id="policy-components"></a>

<ParameterDescription
    name='policy.components'
    description='List of additional circuit components.'
    type='Array of Object (aperture.spec.v1.Component)'
    reference='../../../spec#component'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-evaluation-interval"></a>

<ParameterDescription
    name='policy.evaluation_interval'
    description='The interval between successive evaluations of the Circuit.'
    type='string'
    reference=''
    value='"10s"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-load-ramp"></a>

<ParameterDescription
    name='policy.load_ramp'
    description='Identify the service and flows of the feature that needs to be rolled out. And specify feature rollout steps.'
    type='Object (aperture.spec.v1.LoadRampParameters)'
    reference='../../../spec#load-ramp-parameters'
    value='{"sampler": {"label_key": "", "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}, "steps": [{"duration": "__REQUIRED_FIELD__", "target_accept_percentage": "__REQUIRED_FIELD__"}]}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-policy-name"></a>

<ParameterDescription
    name='policy.policy_name'
    description='Name of the policy.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-resources"></a>

<ParameterDescription
    name='policy.resources'
    description='List of additional resources.'
    type='Object (aperture.spec.v1.Resources)'
    reference='../../../spec#resources'
    value='{"flow_control": {"classifiers": []}}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-rollout"></a>

<ParameterDescription
    name='policy.rollout'
    description='Whether to start the rollout. This setting may be overridden at runtime via dynamic configuration.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.drivers {#policy-drivers}

<!-- vale on -->

<!-- vale off -->

<a id="policy-drivers-percentile-latency-drivers"></a>

<ParameterDescription
    name='policy.drivers.percentile_latency_drivers'
    description='List of drivers that compare percentile latency against forward, backward and reset thresholds.'
    type='Array of Object (percentile_latency_driver)'
    reference='#percentile-latency-driver'
    value='[{"criteria": {"forward": {"threshold": "__REQUIRED_FIELD__"}}, "flux_meter": {"selector": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}], "static_buckets": {"buckets": [5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000]}}, "percentile": 95}]'
/>

<!-- vale on -->

---

<!-- vale off -->

#### dashboard {#dashboard}

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

<a id="dashboard-datasource-filter-regex"></a>

<ParameterDescription
    name='dashboard.datasource.filter_regex'
    description='Datasource filter regex.'
    type='string'
    reference=''
    value='""'
/>

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

---

### Schemas

<!-- vale off -->

#### criteria {#criteria}

<!-- vale on -->

<!-- vale off -->

##### backward {#criteria-backward}

<!-- vale on -->

<!-- vale off -->

<a id="criteria-backward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

##### forward {#criteria-forward}

<!-- vale on -->

<!-- vale off -->

<a id="criteria-forward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the forward criteria.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

##### reset {#criteria-reset}

<!-- vale on -->

<!-- vale off -->

<a id="criteria-reset-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the reset criteria.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

---

<!-- vale off -->

#### percentile_latency_driver {#percentile-latency-driver}

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-criteria"></a>

<ParameterDescription
    name='criteria'
    description='The criteria for percentile latency comparison.'
    type='Object (criteria)'
    reference='#criteria'
    value='{"forward": {"threshold": "__REQUIRED_FIELD__"}}'
/>

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
