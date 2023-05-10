---
title: Static Rate Limiting Policy
keywords:
  - blueprints
sidebar_position: 4
sidebar_label: Static Rate Limiting Policy
---

## Introduction

This blueprint provides a simple static rate-limiting policy and a dashboard.
This policy uses the [`RateLimiter`](/reference/policies/spec.md#rate-limiter)
component.

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
import {ParameterDescription} from '../../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/static-rate-limiting`}>policies/static-rate-limiting</a>

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

##### policy.rate_limiter {#policy-rate-limiter}

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-rate-limit"></a>

<ParameterDescription
    name='policy.rate_limiter.rate_limit'
    description='Number of requests per `policy.rate_limiter.parameters.limit_reset_interval` to accept'
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
    value='{"label_key": "__REQUIRED_FIELD__", "limit_reset_interval": "__REQUIRED_FIELD__"}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-default-config"></a>

<ParameterDescription
    name='policy.rate_limiter.default_config'
    description='Default configuration for rate limiter that can be updated at the runtime without shutting down the policy.'
    type='Object (aperture.spec.v1.RateLimiterDynamicConfig)'
    reference='../../spec#rate-limiter-dynamic-config'
    value='{"overrides": []}'
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

## Dynamic Configuration

:::note

The following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

:::

### Parameters

<!-- vale off -->

<a id="rate-limiter"></a>

<ParameterDescription
    name='rate_limiter'
    description='Rate limiter dynamic configuration that is updated at runtime.'
    type='Object (aperture.spec.v1.RateLimiterDynamicConfig)'
    reference='../../spec#rate-limiter-dynamic-config'
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

---
