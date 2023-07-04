---
title: Base Load Ramping Policy
keywords:
  - blueprints
sidebar_label: Base Load Ramping Policy
---

## Introduction

This policy rolls out new features based on closed loop feedback. The ramping
criteria are defined by drivers that determine conditions for advancing,
reversing, or resetting the ramping to its initial state. The ramping process
consists of a series of steps that progress if the feature is considered
healthy.

:::info

Please see reference for the
[`LoadRamp`](/reference/configuration/spec.md#load-ramp) component that is used
within this blueprint.

:::

:::info

See the use-case
[Percentage Rollout with Average Latency Feedback](/use-cases/percentage-rollouts/average-latency-feedback.md)
to see this blueprint in use.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../apertureVersion.js'
import {ParameterDescription} from '../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/load-ramping/base`}>load-ramping/base</a>

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
    reference='../../spec#component'
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
    description='Identify the service and flows of the feature that needs to be rolled out. And specify load ramp steps.'
    type='Object (aperture.spec.v1.LoadRampParameters)'
    reference='../../spec#load-ramp-parameters'
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
    reference='../../spec#resources'
    value='{"flow_control": {"classifiers": []}}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-start"></a>

<ParameterDescription
    name='policy.start'
    description='Whether to start the ramp. This setting may be overridden at runtime via dynamic configuration.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.drivers {#policy-drivers}

<!-- vale on -->

<!-- vale off -->

<a id="policy-drivers-average-latency-drivers"></a>

<ParameterDescription
    name='policy.drivers.average_latency_drivers'
    description='List of drivers that compare average latency against forward, backward and reset thresholds.'
    type='Array of Object (average_latency_driver)'
    reference='#average-latency-driver'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-drivers-percentile-latency-drivers"></a>

<ParameterDescription
    name='policy.drivers.percentile_latency_drivers'
    description='List of drivers that compare percentile latency against forward, backward and reset thresholds.'
    type='Array of Object (percentile_latency_driver)'
    reference='#percentile-latency-driver'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-drivers-promql-drivers"></a>

<ParameterDescription
    name='policy.drivers.promql_drivers'
    description='List of promql drivers that compare results of a Prometheus query against forward, backward and reset thresholds.'
    type='Array of Object (promql_driver)'
    reference='#promql-driver'
    value='[]'
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
    value='"Aperture Load Ramp"'
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

#### average_latency_driver {#average-latency-driver}

<!-- vale on -->

<!-- vale off -->

<a id="average-latency-driver-selectors"></a>

<ParameterDescription
    name='selectors'
    description='Identify the service and flows whose latency needs to be measured.'
    type='Array of Object (aperture.spec.v1.Selector)'
    reference='../../spec#selector'
    value='[{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

<!-- vale off -->

##### criteria {#average-latency-driver-criteria}

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
    value='null'
/>

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

###### reset {#average-latency-driver-criteria-reset}

<!-- vale on -->

<!-- vale off -->

<a id="average-latency-driver-criteria-reset-threshold"></a>

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

<a id="percentile-latency-driver-flux-meter"></a>

<ParameterDescription
    name='flux_meter'
    description='FluxMeter specifies the flows whose latency needs to be measured and parameters for the histogram metrics.'
    type='Object (aperture.spec.v1.FluxMeter)'
    reference='../../spec#flux-meter'
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

###### backward {#percentile-latency-driver-criteria-backward}

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-criteria-backward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='null'
/>

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

###### reset {#percentile-latency-driver-criteria-reset}

<!-- vale on -->

<!-- vale off -->

<a id="percentile-latency-driver-criteria-reset-threshold"></a>

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

###### backward {#promql-driver-criteria-backward}

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-criteria-backward-operator"></a>

<ParameterDescription
    name='operator'
    description='The operator for the backward criteria. oneof: `gt | lt | gte | lte | eq | neq`'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-criteria-backward-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the backward criteria.'
    type='Number (double)'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

###### forward {#promql-driver-criteria-forward}

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

###### reset {#promql-driver-criteria-reset}

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-criteria-reset-operator"></a>

<ParameterDescription
    name='operator'
    description='The operator for the reset criteria. oneof: `gt | lt | gte | lte | eq | neq`'
    type='string'
    reference=''
    value='null'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-driver-criteria-reset-threshold"></a>

<ParameterDescription
    name='threshold'
    description='The threshold for the reset criteria.'
    type='Number (double)'
    reference=''
    value='null'
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
    description='Reset load ramp to the first step. This setting can be updated at the runtime without shutting down the policy.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

---

<!-- vale off -->

<a id="start"></a>

<ParameterDescription
    name='start'
    description='Start load ramp. This setting can be updated at runtime without shutting down the policy. The load ramp gets paused if this flag is set to false in the middle of a load ramp.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

---
