---
title: Quota Scheduler Policy
keywords:
  - blueprints
sidebar_position: 6
sidebar_label: Quota Scheduler Policy
---

## Introduction

This blueprint provides a
[token bucket](https://en.wikipedia.org/wiki/Token_bucket) based quota scheduler
policy and a dashboard. This policy uses the
[`QuotaScheduler`](/reference/policies/spec.md#quota-scheduler) component.

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/quota-scheduler`}>policies/quota-scheduler</a>

<!-- vale on -->

### Parameters

<!-- vale off -->

#### policy {#policy}

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

<a id="policy-classifiers"></a>

<ParameterDescription
    name='policy.classifiers'
    description='List of classification rules.'
    type='Array of Object (aperture.spec.v1.Classifier)'
    reference='../../spec#classifier'
    value='[]'
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

<a id="policy-quota-scheduler-selectors"></a>

<ParameterDescription
    name='policy.quota_scheduler.selectors'
    description='Flow selectors to match requests against'
    type='Array of Object (aperture.spec.v1.Selector)'
    reference='../../spec#selector'
    value='[{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-quota-scheduler-parameters"></a>

<ParameterDescription
    name='policy.quota_scheduler.parameters'
    description='Parameters.'
    type='Object (aperture.spec.v1.RateLimiterParameters)'
    reference='../../spec#rate-limiter-parameters'
    value='{"interval": "__REQUIRED_FIELD__", "label_key": ""}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-quota-scheduler-scheduler"></a>

<ParameterDescription
    name='policy.quota_scheduler.scheduler'
    description='Scheduler configuration.'
    type='Object (aperture.spec.v1.Scheduler)'
    reference='../../spec#scheduler'
    value='{}'
/>

<!-- vale on -->

---

<!-- vale off -->

#### dashboard {#dashboard}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-refresh-interval"></a>

<ParameterDescription
    name='dashboard.refresh_interval'
    description='Refresh interval for dashboard panels.'
    type='string'
    reference=''
    value='"10s"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-from"></a>

<ParameterDescription
    name='dashboard.time_from'
    description='Time from of dashboard.'
    type='string'
    reference=''
    value='"now-15m"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-time-to"></a>

<ParameterDescription
    name='dashboard.time_to'
    description='Time to of dashboard.'
    type='string'
    reference=''
    value='"now"'
/>

<!-- vale on -->

<!-- vale off -->

##### dashboard.datasource {#dashboard-datasource}

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

---
