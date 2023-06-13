---
title: Rate Limiting Policy
keywords:
  - blueprints
sidebar_position: 5
sidebar_label: Rate Limiting Policy
---

## Introduction

This blueprint provides a
[token bucket](https://en.wikipedia.org/wiki/Token_bucket) based rate-limiting
policy and a dashboard. This policy uses the
[`RateLimiter`](/reference/policies/spec.md#rate-limiter) component.

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/rate-limiting`}>policies/rate-limiting</a>

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

<a id="policy-components"></a>

<ParameterDescription
    name='policy.components'
    description='List of additional circuit components.'
    type='Array of Object (aperture.spec.v1.Component)'
    reference='../../spec#component'
    value='[]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-resources"></a>

<ParameterDescription
    name='policy.resources'
    description='Additional resources.'
    type='Object (aperture.spec.v1.Resources)'
    reference='../../spec#resources'
    value='{"flow_control": {"classifiers": []}}'
/>

<!-- vale on -->

<!-- vale off -->

##### policy.rate_limiter {#policy-rate-limiter}

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-bucket-capacity"></a>

<ParameterDescription
    name='policy.rate_limiter.bucket_capacity'
    description='Bucket capacity.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-fill-amount"></a>

<ParameterDescription
    name='policy.rate_limiter.fill_amount'
    description='Fill amount.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-selectors"></a>

<ParameterDescription
    name='policy.rate_limiter.selectors'
    description='Flow selectors to match requests against'
    type='Array of Object (aperture.spec.v1.Selector)'
    reference='../../spec#selector'
    value='[{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-parameters"></a>

<ParameterDescription
    name='policy.rate_limiter.parameters'
    description='Parameters.'
    type='Object (aperture.spec.v1.RateLimiterParameters)'
    reference='../../spec#rate-limiter-parameters'
    value='{"interval": "__REQUIRED_FIELD__", "label_key": ""}'
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

<a id="dashboard-title"></a>

<ParameterDescription
    name='dashboard.title'
    description='Name of the main dashboard.'
    type='string'
    reference=''
    value='"Aperture Rate Limiter"'
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
