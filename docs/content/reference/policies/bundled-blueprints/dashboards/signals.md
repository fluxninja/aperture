---
title: Signals
---

## Introduction

This blueprint provides a [policy monitoring](/reference/policies/monitoring.md)
dashboard that visualizes Signals flowing through the
[Circuit](/concepts/policy/circuit.md).

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

## Configuration: <a href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/dashboards/signals`}>dashboards/signals</a> {#configuration}

### Parameters

<a id="common"></a>

### common

<a id="common-policy-name"></a> <ParameterDescription
    name="common.policy_name"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

<a id="dashboard"></a>

### dashboard

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

<a id="dashboard-datasource"></a>

#### dashboard.datasource

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
