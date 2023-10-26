---
title: Base Quota Scheduling Policy
keywords:
  - blueprints
sidebar_label: Base Quota Scheduling Policy
---

## Introduction

This blueprint provides a
[token bucket](https://en.wikipedia.org/wiki/Token_bucket) based quota scheduler
policy and a dashboard. This policy uses the
[`QuotaScheduler`](/reference/configuration/spec.md#quota-scheduler) component.

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../apertureVersion.js'
import {ParameterDescription} from '../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/quota-scheduling/base`}>quota-scheduling/base</a>

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

##### policy.quota_scheduler {#policy-quota-scheduler}

<!-- vale on -->

<!-- vale off -->

<a id="policy-quota-scheduler-bucket-capacity"></a>

<ParameterDescription
    name='policy.quota_scheduler.bucket_capacity'
    description='Bucket capacity.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-quota-scheduler-fill-amount"></a>

<ParameterDescription
    name='policy.quota_scheduler.fill_amount'
    description='Fill amount.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-quota-scheduler-rate-limiter"></a>

<ParameterDescription
    name='policy.quota_scheduler.rate_limiter'
    description='Rate Limiter Parameters.'
    type='Object (aperture.spec.v1.RateLimiterParameters)'
    reference='../../configuration/spec#rate-limiter-parameters'
    value='{"interval": "__REQUIRED_FIELD__", "label_key": ""}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-quota-scheduler-scheduler"></a>

<ParameterDescription
    name='policy.quota_scheduler.scheduler'
    description='Scheduler configuration.'
    type='Object (aperture.spec.v1.Scheduler)'
    reference='../../configuration/spec#scheduler'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-quota-scheduler-selectors"></a>

<ParameterDescription
    name='policy.quota_scheduler.selectors'
    description='Flow selectors to match requests against'
    type='Array of Object (aperture.spec.v1.Selector)'
    reference='../../configuration/spec#selector'
    value='[{"control_point": "__REQUIRED_FIELD__"}]'
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
    value='"Aperture Quota Scheduler"'
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
