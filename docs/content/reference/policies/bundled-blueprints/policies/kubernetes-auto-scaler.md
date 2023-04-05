---
title: Kubernetes Auto Scaler Policy
---

## Introduction

This blueprint provides dashboard and policy which auto scales the targeted
Kubernetes resources based on the results received by executing the provided
PromQL queries for scale-in and scale-out.

This policy uses the
[`PodAutoScaler`](/reference/policies/spec.md#pod-auto-scaler) component.

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
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/kubernetes-auto-scaler`}>policies/kubernetes-auto-scaler</a>

<h3 class="blueprints-h3">Common</h3>

<ParameterDescription
    name="common.policy_name"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

<h3 class="blueprints-h3">Policy</h3>

<ParameterDescription
    name="policy.min_replicas"
    type="string"
    reference=""
    value="1"
    description='Minimum number of replicas.' />

<ParameterDescription
    name="policy.max_replicas"
    type="string"
    reference=""
    value="10"
    description='Maximum number of replicas.' />

<ParameterDescription
    name="policy.scale_in_cooldown"
    type="string"
    reference=""
    value="'40s'"
    description='The amount of time to wait after a scale-in operation for another scale-in operation.' />

<ParameterDescription
    name="policy.scale_out_cooldown"
    type="string"
    reference=""
    value="'30s'"
    description='The amount of time to wait after a scale-out operation for another scale-out or scale-in operation.' />

<ParameterDescription
    name="policy.cooldown_override_percentage"
    type="number"
    reference=""
    value="50"
    description='Cooldown override percentage defines a threshold change in scale-out beyond which previous cooldown is overridden.' />

<ParameterDescription
    name="policy.max_scale_in_percentage"
    type="number"
    reference=""
    value="1"
    description='The maximum decrease of replicas (e.g. pods) at one time.' />

<ParameterDescription
    name="policy.max_scale_out_percentage"
    type="number"
    reference=""
    value="10"
    description='The maximum increase of replicas (e.g. pods) at one time.' />

<ParameterDescription
    name="policy.scale_in_alerter_parameters"
    type="aperture.spec.v1.AlerterParameters"
    reference="../../spec#alerter-parameters"
    value="{'alert_name': 'Kubernetes Auto Scaler Scale In Event'}"
    description='Configuration for scale-in alerter.' />

<ParameterDescription
    name="policy.scale_in_alerter_parameters.alert_name"
    type="string"
    reference=""
    value="'Kubernetes Auto Scaler Scale In Event'"
    description='Name of the alert.' />

<ParameterDescription
    name="policy.scale_out_alerter_parameters"
    type="aperture.spec.v1.AlerterParameters"
    reference="../../spec#alerter-parameters"
    value="{'alert_name': 'Kubernetes Auto Scaler Scale Out Event'}"
    description='Cooldown override percentage.' />

<ParameterDescription
    name="policy.scale_out_alerter_parameters.alert_name"
    type="string"
    reference=""
    value="'Kubernetes Auto Scaler Scale Out Event'"
    description='Configuration for scale-out alerter.' />

<ParameterDescription
    name="policy.components"
    type="[]aperture.spec.v1.Component"
    reference="../../spec#component"
    value="[]"
    description='List of additional circuit components.' />

<h4 class="blueprints-h4">Kubernetes Object Selector</h4>

<ParameterDescription
    name="policy.kubernetes_object_selector.namespace"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Namespace.' />

<ParameterDescription
    name="policy.kubernetes_object_selector.api_version"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='API Version.' />

<ParameterDescription
    name="policy.kubernetes_object_selector.kind"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Kind.' />

<ParameterDescription
    name="policy.kubernetes_object_selector.name"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name.' />

<h4 class="blueprints-h4">Scale-in criteria</h4>

<ParameterDescription
    name="policy.scale_in_criteria"
    type="[]object"
    reference=""
    value="[{'parameters': {'slope': 1}, 'query': {'promql': {'evaluation_interval': '10s', 'out_ports': {'output': {'signal_name': '__REQUIRED_FIELD__'}}, 'query_string': '__REQUIRED_FIELD__'}}, 'set_point': 0.5}]"
    description='List of scale-in criteria.' />

<ParameterDescription
    name="policy.scale_in_criteria.query"
    type="aperture.spec.v1.Query"
    reference="../../spec#query"
    value="{'promql': {'evaluation_interval': '10s', 'out_ports': {'output': {'signal_name': '__REQUIRED_FIELD__'}}, 'query_string': '__REQUIRED_FIELD__'}}"
    description='Query.' />

<ParameterDescription
    name="policy.scale_in_criteria.query.promql"
    type="aperture.spec.v1.PromQL"
    reference="../../spec#prom-q-l"
    value="{'evaluation_interval': '10s', 'out_ports': {'output': {'signal_name': '__REQUIRED_FIELD__'}}, 'query_string': '__REQUIRED_FIELD__'}"
    description='PromQL query.' />

<ParameterDescription
    name="policy.scale_in_criteria.query.promql.query_string"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='PromQL query string.' />

<ParameterDescription
    name="policy.scale_in_criteria.query.promql.evaluation_interval"
    type="string"
    reference=""
    value="'10s'"
    description='Evaluation interval.' />

<ParameterDescription
    name="policy.scale_in_criteria.query.promql.out_ports"
    type="aperture.spec.v1.PromQLOuts"
    reference="../../spec#prom-q-l-outs"
    value="{'output': {'signal_name': '__REQUIRED_FIELD__'}}"
    description='PromQL query execution output.' />

<ParameterDescription
    name="policy.scale_in_criteria.query.promql.out_ports.output"
    type="aperture.spec.v1.OutPort"
    reference="../../spec#out-port"
    value="{'signal_name': '__REQUIRED_FIELD__'}"
    description='PromQL query execution output port.' />

<ParameterDescription
    name="policy.scale_in_criteria.query.promql.out_ports.output.signal_name"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Output Signal name.' />

<ParameterDescription
    name="policy.scale_in_criteria.set_point"
    type="number"
    reference=""
    value="0.5"
    description='Set point.' />

<ParameterDescription
    name="policy.scale_in_criteria.parameters"
    type="aperture.spec.v1.DecreasingGradientParameters"
    reference="../../spec#decreasing-gradient-parameters"
    value="{'slope': 1}"
    description='Parameters.' />

<ParameterDescription
    name="policy.scale_in_criteria.parameters.slope"
    type="number"
    reference=""
    value="1"
    description='Slope.' />

<h4 class="blueprints-h4">Scale-out criteria</h4>

<ParameterDescription
    name="policy.scale_out_criteria"
    type="[]object"
    reference=""
    value="[{'parameters': {'slope': -1}, 'query': {'promql': {'evaluation_interval': '10s', 'out_ports': {'output': {'signal_name': '__REQUIRED_FIELD__'}}, 'query_string': '__REQUIRED_FIELD__'}}, 'set_point': 1}]"
    description='List of scale-out criteria.' />

<ParameterDescription
    name="policy.scale_out_criteria.query"
    type="aperture.spec.v1.Query"
    reference="../../spec#query"
    value="{'promql': {'evaluation_interval': '10s', 'out_ports': {'output': {'signal_name': '__REQUIRED_FIELD__'}}, 'query_string': '__REQUIRED_FIELD__'}}"
    description='Query.' />

<ParameterDescription
    name="policy.scale_out_criteria.query.promql"
    type="aperture.spec.v1.PromQL"
    reference="../../spec#prom-q-l"
    value="{'evaluation_interval': '10s', 'out_ports': {'output': {'signal_name': '__REQUIRED_FIELD__'}}, 'query_string': '__REQUIRED_FIELD__'}"
    description='PromQL query.' />

<ParameterDescription
    name="policy.scale_out_criteria.query.promql.query_string"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='PromQL query string.' />

<ParameterDescription
    name="policy.scale_out_criteria.query.promql.evaluation_interval"
    type="string"
    reference=""
    value="'10s'"
    description='Evaluation interval.' />

<ParameterDescription
    name="policy.scale_out_criteria.query.promql.out_ports"
    type="aperture.spec.v1.PromQLOuts"
    reference="../../spec#prom-q-l-outs"
    value="{'output': {'signal_name': '__REQUIRED_FIELD__'}}"
    description='PromQL query execution output.' />

<ParameterDescription
    name="policy.scale_out_criteria.query.promql.out_ports.output"
    type="aperture.spec.v1.OutPort"
    reference="../../spec#out-port"
    value="{'signal_name': '__REQUIRED_FIELD__'}"
    description='PromQL query execution output port.' />

<ParameterDescription
    name="policy.scale_out_criteria.query.promql.out_ports.output.signal_name"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Output Signal name.' />

<ParameterDescription
    name="policy.scale_out_criteria.set_point"
    type="number"
    reference=""
    value="1"
    description='Set point.' />

<ParameterDescription
    name="policy.scale_out_criteria.parameters"
    type="aperture.spec.v1.IncreasingGradientParameters"
    reference="../../spec#increasing-gradient-parameters"
    value="{'slope': -1}"
    description='Parameters.' />

<ParameterDescription
    name="policy.scale_out_criteria.parameters.slope"
    type="number"
    reference=""
    value="-1"
    description='Slope.' />

<h3 class="blueprints-h3">Dashboard</h3>

<ParameterDescription
    name="dashboard.refresh_interval"
    type="string"
    reference=""
    value="'5s'"
    description='Refresh interval for dashboard panels.' />

<ParameterDescription
    name="dashboard.time_from"
    type="string"
    reference=""
    value="'now-15m'"
    description='From time of dashboard.' />

<ParameterDescription
    name="dashboard.time_to"
    type="string"
    reference=""
    value="'now'"
    description='To time of dashboard.' />

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
    description='Datasource filter regex.' />
