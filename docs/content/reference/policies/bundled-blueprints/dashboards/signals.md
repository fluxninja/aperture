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

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/dashboards/signals`}>dashboards/signals</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

#### common {#common}

<!-- vale on -->

<!-- vale off -->

<a id="common-policy-name"></a>

<ParameterDescription
    name="common.policy_name"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.'
/>

<!-- vale on -->

---

<!-- vale off -->

#### dashboard {#dashboard}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-refresh-interval"></a>

<ParameterDescription
    name="dashboard.refresh_interval"
    type="
string"
    reference=""
    value="'5s'"
    description='Refresh interval for dashboard panels.'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-from"></a>

<ParameterDescription
    name="dashboard.time_from"
    type="
string"
    reference=""
    value="'now-15m'"
    description='From time of dashboard.'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-to"></a>

<ParameterDescription
    name="dashboard.time_to"
    type="
string"
    reference=""
    value="'now'"
    description='To time of dashboard.'
/>

<!-- vale on -->

<!-- vale off -->

##### dashboard.datasource {#dashboard-datasource}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-name"></a>

<ParameterDescription
    name="dashboard.datasource.name"
    type="
string"
    reference=""
    value="'$datasource'"
    description='Datasource name.'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-filter-regex"></a>

<ParameterDescription
    name="dashboard.datasource.filter_regex"
    type="
string"
    reference=""
    value="''"
    description='Datasource filter regex.'
/>

<!-- vale on -->

---
