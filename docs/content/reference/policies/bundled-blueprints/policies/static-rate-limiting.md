---
title: Static Rate Limiting Policy
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

#### common {#common}

<!-- vale on -->

<!-- vale off -->

<a id="common-policy-name"></a> <ParameterDescription
    name="common.policy_name"
    type="
string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

<!-- vale on -->

---

<!-- vale off -->

#### policy {#policy}

<!-- vale on -->

<!-- vale off -->

<a id="policy-classifiers"></a> <ParameterDescription
    name="policy.classifiers"
    type="
Array of
Object (aperture.spec.v1.Classifier)"
    reference="../../spec#classifier"
    value="[]"
    description='List of classification rules.' />

<!-- vale on -->

<!-- vale off -->

##### policy.rate_limiter {#policy-rate-limiter}

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-rate-limit"></a> <ParameterDescription
    name="policy.rate_limiter.rate_limit"
    type="
Number (double)"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Number of requests per `policy.rate_limiter.parameters.limit_reset_interval` to accept' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-flow-selector"></a> <ParameterDescription
    name="policy.rate_limiter.flow_selector"
    type="
Object (aperture.spec.v1.FlowSelector)"
    reference="../../spec#flow-selector"
    value="{'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}"
    description='A flow selector to match requests against' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-parameters"></a> <ParameterDescription
    name="policy.rate_limiter.parameters"
    type="
Object (aperture.spec.v1.RateLimiterParameters)"
    reference="../../spec#rate-limiter-parameters"
    value="{'label_key': '__REQUIRED_FIELD__', 'limit_reset_interval': '__REQUIRED_FIELD__'}"
    description='Parameters.' />

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-default-config"></a> <ParameterDescription
    name="policy.rate_limiter.default_config"
    type="
Object (aperture.spec.v1.RateLimiterDynamicConfig)"
    reference="../../spec#rate-limiter-dynamic-config"
    value="{'overrides': []}"
    description='Default configuration for rate limiter that can be updated at the runtime without shutting down the policy.' />

<!-- vale on -->

---

<!-- vale off -->

#### dashboard {#dashboard}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-refresh-interval"></a> <ParameterDescription
    name="dashboard.refresh_interval"
    type="
string"
    reference=""
    value="'10s'"
    description='Refresh interval for dashboard panels.' />

<!-- vale on -->

<!-- vale off -->

##### dashboard.datasource {#dashboard-datasource}

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-name"></a> <ParameterDescription
    name="dashboard.datasource.name"
    type="
string"
    reference=""
    value="'$datasource'"
    description='Datasource name.' />

<!-- vale on -->

<!-- vale off -->

<a id="dashboard-datasource-filter-regex"></a> <ParameterDescription
    name="dashboard.datasource.filter_regex"
    type="
string"
    reference=""
    value="''"
    description='Datasource filter regex.' />

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

<a id="rate-limiter"></a> <ParameterDescription
    name="rate_limiter"
    type="
Object (aperture.spec.v1.RateLimiterDynamicConfig)"
    reference="../../spec#rate-limiter-dynamic-config"
    value="__REQUIRED_FIELD__"
    description='Rate limiter dynamic configuration that is updated at runtime.' />

<!-- vale on -->

---
