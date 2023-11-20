---
title: Base Rate Limiting Policy
keywords:
  - blueprints
sidebar_label: Base Rate Limiting Policy
---

## Introduction

This blueprint provides a
[token bucket](https://en.wikipedia.org/wiki/Token_bucket) based rate-limiting
policy and a dashboard. This policy uses the
[`RateLimiter`](/reference/configuration/spec.md#rate-limiter) component.

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../apertureVersion.js'
import {ParameterDescription} from '../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/rate-limiting/base`}>rate-limiting/base</a>

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

##### policy.rate_limiter {#policy-rate-limiter}

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-alerter"></a>

<ParameterDescription
    name='policy.rate_limiter.alerter'
    description='Alerter.'
    type='Object (aperture.spec.v1.AlerterParameters)'
    reference='../../configuration/spec#alerter-parameters'
    value='{"alert_name": "More than 90% of requests are being rate limited"}'
/>

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

<a id="policy-rate-limiter-parameters"></a>

<ParameterDescription
    name='policy.rate_limiter.parameters'
    description='Parameters.'
    type='Object (aperture.spec.v1.RateLimiterParameters)'
    reference='../../configuration/spec#rate-limiter-parameters'
    value='{"interval": "__REQUIRED_FIELD__", "limit_by_label_key": ""}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-request-parameters"></a>

<ParameterDescription
    name='policy.rate_limiter.request_parameters'
    description='Request Parameters.'
    type='Object (aperture.spec.v1.RateLimiterRequestParameters)'
    reference='../../configuration/spec#rate-limiter-request-parameters'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-rate-limiter-selectors"></a>

<ParameterDescription
    name='policy.rate_limiter.selectors'
    description='Flow selectors to match requests against'
    type='Array of Object (aperture.spec.v1.Selector)'
    reference='../../configuration/spec#selector'
    value='[{"control_point": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

---
