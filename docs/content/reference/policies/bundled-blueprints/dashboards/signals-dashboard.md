---
title: Signals Dashboard
---

```mdx-code-block
import {apertureVersion} from '../../../../apertureVersion.js';
```

## Blueprint Location

GitHub: <a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/blueprints/lib/1.0/dashboards/signals-dashboard`}>signals-dashboard</a>

## Introduction

This blueprint provides a [policy monitoring](/reference/policies/monitoring.md)
dashboard that visualizes Signals flowing through the
[Circuit](/concepts/policy/circuit.md).

## Configuration

<!-- Configuration Marker -->

### Common

**`common.policy_name`** (type: _`string`_)

required parameter

Name of the policy.

### Dashboard

**`dashboard.refresh_interval`** (type: _`string`_)

default: `"10s"`

Refresh interval for dashboard panels.

**`dashboard.time_from`** (type: _`string`_)

default: `"now-30m"`

From time of dashboard.

**`dashboard.time_to`** (type: _`string`_)

default: `"now"`

To time of dashboard.

#### Datasource

**`dashboard.datasource.name`** (type: _`string`_)

default: `"$datasource"`

Datasource name.

**`dashboard.datasource.filter_regex`** (type: _`string`_)

default: `""`

Datasource filter regex.
