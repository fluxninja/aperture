---
title: Static Rate Limiting Policy
keywords:
  - blueprints
sidebar_position: 5
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

<a id="policy-rate-limit"></a>

<ParameterDescription
    name='policy.rate_limit'
    description='Number of requests per `policy.rate_limiter.limit_reset_interval` to accept'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter"></a>

<ParameterDescription
    name='policy.rate_limiter'
    description='Parameters for _Rate Limiter_.'
    type='Object (aperture.spec.v1.RateLimiterParameters)'
    reference='../../spec#rate-limiter-parameters'
    value='{"label_key": "__REQUIRED_FIELD__", "limit_reset_interval": "__REQUIRED_FIELD__", "selectors": [{"control_point": "__REQUIRED_FIELD__", "service": "__REQUIRED_FIELD__"}]}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-custom-limits"></a>

<ParameterDescription
    name='policy.custom_limits'
    description='Allows to specify different limits for particular label values. This setting can be updated at runtime without restarting the policy via dynamic config.'
    type='Array of Object (aperture.spec.v1.RateLimiterCustomLimit)'
    reference='../../spec#rate-limiter-custom-limit'
    value='[]'
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

<a id="custom-limits"></a>

<ParameterDescription
    name='custom_limits'
    description='Allows to specify different limits for particular label values.'
    type='Array of Object (aperture.spec.v1.RateLimiterCustomLimit)'
    reference='../../spec#rate-limiter-custom-limit'
    value='[{"label_value": "__REQUIRED_FIELD__", "limit_scale_factor": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

---
