---
title: Kubernetes Auto Scaler Policy
---

## Introduction

This blueprint provides dashboard and policy which auto scales the targeted
Kubernetes resources based on the results received by executing the provided
PromQL queries for scale-in and scale-out.

This policy uses the
[`PodAutoScaler`](/reference/policies/spec.md#pod-auto-scaler) component.

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

Code: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/kubernetes-auto-scaler`}>policies/kubernetes-auto-scaler</a>

### Parameters

#### common {#common}

<a id="common-policy-name"></a> <ParameterDescription
    name="common.policy_name"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

#### policy {#policy}

<a id="policy-min-replicas"></a> <ParameterDescription
    name="policy.min_replicas"
    type="
string"
    reference=""
    value="1"
    description='Minimum number of replicas.' />

<a id="policy-max-replicas"></a> <ParameterDescription
    name="policy.max_replicas"
    type="
string"
    reference=""
    value="10"
    description='Maximum number of replicas.' />

<a id="policy-scale-in-cooldown"></a> <ParameterDescription
    name="policy.scale_in_cooldown"
    type="
string"
    reference=""
    value="'40s'"
    description='The amount of time to wait after a scale-in operation for another scale-in operation.' />

<a id="policy-scale-out-cooldown"></a> <ParameterDescription
    name="policy.scale_out_cooldown"
    type="
string"
    reference=""
    value="'30s'"
    description='The amount of time to wait after a scale-out operation for another scale-out or scale-in operation.' />

<a id="policy-cooldown-override-percentage"></a> <ParameterDescription
    name="policy.cooldown_override_percentage"
    type="
Number (double)"
    reference=""
    value="50"
    description='Cooldown override percentage defines a threshold change in scale-out beyond which previous cooldown is overridden.' />

<a id="policy-max-scale-in-percentage"></a> <ParameterDescription
    name="policy.max_scale_in_percentage"
    type="
Number (double)"
    reference=""
    value="1"
    description='The maximum decrease of replicas (e.g. pods) at one time.' />

<a id="policy-max-scale-out-percentage"></a> <ParameterDescription
    name="policy.max_scale_out_percentage"
    type="
Number (double)"
    reference=""
    value="10"
    description='The maximum increase of replicas (e.g. pods) at one time.' />

<a id="policy-scale-in-alerter-parameters"></a> <ParameterDescription
    name="policy.scale_in_alerter_parameters"
    type="
Object (aperture.spec.v1.AlerterParameters)"
    reference="../../spec#alerter-parameters"
    value="{'alert_name': 'Kubernetes Auto Scaler Scale In Event'}"
    description='Configuration for scale-in alerter.' />

<a id="policy-scale-out-alerter-parameters"></a> <ParameterDescription
    name="policy.scale_out_alerter_parameters"
    type="
Object (aperture.spec.v1.AlerterParameters)"
    reference="../../spec#alerter-parameters"
    value="{'alert_name': 'Kubernetes Auto Scaler Scale Out Event'}"
    description='Cooldown override percentage.' />

<a id="policy-components"></a> <ParameterDescription
    name="policy.components"
    type="
Array of
Object (aperture.spec.v1.Component)"
    reference="../../spec#component"
    value="[]"
    description='List of additional circuit components.' />

##### policy.kubernetes_object_selector {#policy-kubernetes-object-selector}

<a id="policy-kubernetes-object-selector-namespace"></a> <ParameterDescription
    name="policy.kubernetes_object_selector.namespace"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Namespace.' />

<a id="policy-kubernetes-object-selector-api-version"></a> <ParameterDescription
    name="policy.kubernetes_object_selector.api_version"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='API Version.' />

<a id="policy-kubernetes-object-selector-kind"></a> <ParameterDescription
    name="policy.kubernetes_object_selector.kind"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Kind.' />

<a id="policy-kubernetes-object-selector-name"></a> <ParameterDescription
    name="policy.kubernetes_object_selector.name"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name.' />

<a id="policy-scale-in-criteria"></a> <ParameterDescription
    name="policy.scale_in_criteria"
    type="
Array of
Object (scale_criteria)"
    reference="#scale-criteria"
    value="[{'parameters': {'slope': 1}, 'query': {'promql': {'evaluation_interval': '10s', 'out_ports': {'output': {'signal_name': '__REQUIRED_FIELD__'}}, 'query_string': '__REQUIRED_FIELD__'}}, 'set_point': 0.5}]"
    description='List of scale-in criteria.' />

<a id="policy-scale-out-criteria"></a> <ParameterDescription
    name="policy.scale_out_criteria"
    type="
Array of
Object (scale_criteria)"
    reference="#scale-criteria"
    value="[{'parameters': {'slope': -1}, 'query': {'promql': {'evaluation_interval': '10s', 'out_ports': {'output': {'signal_name': '__REQUIRED_FIELD__'}}, 'query_string': '__REQUIRED_FIELD__'}}, 'set_point': 1}]"
    description='List of scale-out criteria.' />

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

### Schemas

#### scale_criteria {#scale-criteria}

<a id="scale-criteria-query"></a> <ParameterDescription
    name="query"
    type="
Object (aperture.spec.v1.Query)"
    reference="../../spec#query"
    value="{'promql': {'evaluation_interval': '10s', 'out_ports': {'output': {'signal_name': '__REQUIRED_FIELD__'}}, 'query_string': '__REQUIRED_FIELD__'}}"
    description='Query.' />

<a id="scale-criteria-set-point"></a> <ParameterDescription
    name="set_point"
    type="
Number (double)"
    reference=""
    value="1"
    description='Set point.' />

<a id="scale-criteria-parameters"></a> <ParameterDescription
    name="parameters"
    type="
Object (aperture.spec.v1.IncreasingGradientParameters)"
    reference="../../spec#increasing-gradient-parameters"
    value="{'slope': -1}"
    description='Parameters.' />
