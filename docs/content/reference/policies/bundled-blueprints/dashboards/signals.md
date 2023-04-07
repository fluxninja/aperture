---
title: Signals
---

## Introduction

This blueprint provides a [policy monitoring](/reference/policies/monitoring.md)
dashboard that visualizes Signals flowing through the
[Circuit](/concepts/policy/circuit.md).

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

Code: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/dashboards/signals`}>dashboards/signals</a>

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
