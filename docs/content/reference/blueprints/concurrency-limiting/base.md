---
title: Base Concurrency Limiting Policy
keywords:
  - blueprints
sidebar_label: Base Concurrency Limiting Policy
---

## Introduction

This blueprint provides a concurrency limiting policy and a dashboard. This
policy uses the
[`ConcurrencyLimiter`](/reference/configuration/spec.md#concurrency-limiter)
component.

<!-- Configuration Marker -->

```mdx-code-block
import {apertureVersion as aver} from '../../../apertureVersion.js'
import {ParameterDescription} from '../../../parameterComponents.js'
```

## Configuration

<!-- vale off -->

Blueprint name: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/concurrency-limiting/base`}>concurrency-limiting/base</a>

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

##### policy.concurrency_limiter {#policy-concurrency-limiter}

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-limiter-alerter"></a>

<ParameterDescription
    name='policy.concurrency_limiter.alerter'
    description='Alerter.'
    type='Object (aperture.spec.v1.AlerterParameters)'
    reference='../../configuration/spec#alerter-parameters'
    value='{"alert_name": "Too many inflight requests"}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-limiter-max-concurrency"></a>

<ParameterDescription
    name='policy.concurrency_limiter.max_concurrency'
    description='Max concurrency.'
    type='Number (double)'
    reference=''
    value='"__REQUIRED_FIELD__"'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-limiter-parameters"></a>

<ParameterDescription
    name='policy.concurrency_limiter.parameters'
    description='Parameters.'
    type='Object (aperture.spec.v1.ConcurrencyLimiterParameters)'
    reference='../../configuration/spec#concurrency-limiter-parameters'
    value='{"limit_by_label_key": "limit_by_label_key", "max_inflight_duration": "__REQUIRED_FIELD__"}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-limiter-request-parameters"></a>

<ParameterDescription
    name='policy.concurrency_limiter.request_parameters'
    description='Request Parameters.'
    type='Object (aperture.spec.v1.ConcurrencyLimiterRequestParameters)'
    reference='../../configuration/spec#concurrency-limiter-request-parameters'
    value='{}'
/>

<!-- vale on -->

<!-- vale off -->

<a id="policy-concurrency-limiter-selectors"></a>

<ParameterDescription
    name='policy.concurrency_limiter.selectors'
    description='Flow selectors to match requests against'
    type='Array of Object (aperture.spec.v1.Selector)'
    reference='../../configuration/spec#selector'
    value='[{"control_point": "__REQUIRED_FIELD__"}]'
/>

<!-- vale on -->

---
