---
title: Signals Dashboard
---

```mdx-code-block
import {apertureVersion} from '../../../../apertureVersion.js';
```

## Blueprint Location

GitHub: <a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/blueprints/dashboards/signals-dashboard`}>signals-dashboard</a>

Sample values file: <a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/blueprints/dashboards/signals-dashboard/values.yaml`}>values.yaml</a>

Sample values file (required fields only): <a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/blueprints/dashboards/signals-dashboard/values_required.yaml`}>values_required.yaml</a>

## Introduction

This blueprint provides a [policy monitoring](/reference/policies/monitoring.md)
dashboard that visualizes Signals flowing through the
[Circuit](/concepts/policy/circuit.md).

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

<h3 class="blueprints-h3">Common</h3>

<ParameterDescription
    name="common.policy_name"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

<h3 class="blueprints-h3">Dashboard</h3>

<ParameterDescription
    name="dashboard.refresh_interval"
    type="string"
    reference=""
    value="'10s'"
    description='Refresh interval for dashboard panels.' />

<ParameterDescription
    name="dashboard.time_from"
    type="string"
    reference=""
    value="'now-30m'"
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
