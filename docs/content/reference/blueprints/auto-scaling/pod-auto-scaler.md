---
title: Pod Auto-Scaler Policy
---

## Introduction

This blueprint provides a dashboard and policy which auto-scales the targeted
Kubernetes resources based on the results received by executing the provided
PromQL queries for scale-in and scale-out.

:::info

See reference for the [`AutoScale`](/reference/configuration/spec.md#auto-scale)
component that is used within this blueprint.

:::

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../apertureVersion.js'
import {ParameterDescription} from '../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/auto-scaling/pod-auto-scaler`}>auto-scaling/pod-auto-scaler</a>

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
    reference='../../configuration/spec#component'
    value='[]'
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
    description='Additional resources.'
    type='Object (aperture.spec.v1.Resources)'
    reference='../../configuration/spec#resources'
    value='{"flow_control": {"classifiers": []}}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-dry-run"></a>

<ParameterDescription
    name='policy.dry_run'
    description='Dry run mode ensures that no scaling is invoked by this auto scaler.'
    type='Boolean'
    reference=''
    value='false'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-promql-scale-in-controllers"></a>

<ParameterDescription
    name='policy.promql_scale_in_controllers'
    description='List of scale in controllers.'
    type='Array of Object (promql_scale_in_controller)'
    reference='#promql-scale-in-controller'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-promql-scale-out-controllers"></a>

<ParameterDescription
    name='policy.promql_scale_out_controllers'
    description='List of scale out controllers.'
    type='Array of Object (promql_scale_out_controller)'
    reference='#promql-scale-out-controller'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-scaling-backend"></a>

<ParameterDescription
    name='policy.scaling_backend'
    description='Scaling backend for the policy.'
    type='Object (aperture.spec.v1.AutoScalerScalingBackend)'
    reference='../../configuration/spec#auto-scaler-scaling-backend'
    value='{"kubernetes_replicas": "__REQUIRED_FIELD__"}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-scaling-parameters"></a>

<ParameterDescription
    name='policy.scaling_parameters'
    description='Parameters that define the scaling behavior.'
    type='Object (aperture.spec.v1.AutoScalerScalingParameters)'
    reference='../../configuration/spec#auto-scaler-scaling-parameters'
    value='{"scale_in_alerter": {"alert_name": "Auto-scaler is scaling in"}, "scale_out_alerter": {"alert_name": "Auto-scaler is scaling out"}}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-pod-cpu"></a>

<ParameterDescription
    name='policy.pod_cpu'
    description='Driver to do scaling of the resource based on the CPU usage.'
    type='Object (scaling_driver)'
    reference='#scaling-driver'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-pod-memory"></a>

<ParameterDescription
    name='policy.pod_memory'
    description='Driver to do scaling of the resource based on the Memory usage.'
    type='Object (scaling_driver)'
    reference='#scaling-driver'
    value='{}'
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
    value='"15s"'
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
    value='"Aperture Auto-scale"'
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

#### promql_scale_in_controller {#promql-scale-in-controller}

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-alerter"></a>

<ParameterDescription
    name='alerter'
    description='Alerter parameters for the controller.'
    type='Object (aperture.spec.v1.AlerterParameters)'
    reference='../../configuration/spec#alerter-parameters'
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-gradient"></a>

<ParameterDescription
    name='gradient'
    description='Gradient parameters for the controller.'
    type='Object (aperture.spec.v1.DecreasingGradientParameters)'
    reference='../../configuration/spec#decreasing-gradient-parameters'
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-query-string"></a>

<ParameterDescription
    name='query_string'
    description='The Prometheus query to be run. Must return a scalar or a vector with a single element.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-in-controller-setpoint"></a>

<ParameterDescription
    name='setpoint'
    description='Setpoint for the controller.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### promql_scale_out_controller {#promql-scale-out-controller}

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-alerter"></a>

<ParameterDescription
    name='alerter'
    description='Alerter parameters for the controller.'
    type='Object (aperture.spec.v1.AlerterParameters)'
    reference='../../configuration/spec#alerter-parameters'
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-gradient"></a>

<ParameterDescription
    name='gradient'
    description='Gradient parameters for the controller.'
    type='Object (aperture.spec.v1.IncreasingGradientParameters)'
    reference='../../configuration/spec#increasing-gradient-parameters'
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-query-string"></a>

<ParameterDescription
    name='query_string'
    description='The Prometheus query to be run. Must return a scalar or a vector with a single element.'
    type='string'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="promql-scale-out-controller-setpoint"></a>

<ParameterDescription
    name='setpoint'
    description='Setpoint for the controller.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### scaling_criteria {#scaling-criteria}

<!-- vale on -->

<!-- vale off -->

<a id="scaling-criteria-enabled"></a>

<ParameterDescription
    name='enabled'
    description='Enables the driver to do scale in or out of the resource.'
    type='Boolean'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="scaling-criteria-threshold"></a>

<ParameterDescription
    name='threshold'
    description='Threshold for the driver.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---

<!-- vale off -->

#### scaling_driver {#scaling-driver}

<!-- vale on -->

<!-- vale off -->

<a id="scaling-driver-scale-in"></a>

<ParameterDescription
    name='scale_in'
    description='The scale in criteria for the driver.'
    type='Object (scaling_criteria)'
    reference='#scaling-criteria'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="scaling-driver-scale-out"></a>

<ParameterDescription
    name='scale_out'
    description='The scale out criteria for the driver.'
    type='Object (scaling_criteria)'
    reference='#scaling-criteria'
    value='{}'
/>

<!-- vale on -->

---

## Dynamic Configuration

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperturectl/dynamic-config/apply/apply.md)
at runtime, without reloading the policy.

:::

### Parameters

<!-- vale off -->

<a id="dry-run"></a>

<ParameterDescription
    name='dry_run'
    description='Dynamic configuration for setting dry run mode at runtime without restarting this policy. Dry run mode ensures that no scaling is invoked by this auto scaler. This is useful for observing the behavior of auto scaler without disrupting any real deployment.'
    type='Boolean'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---
