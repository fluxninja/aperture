---
title: Static Rate Limiting Policy
---

## Introduction

This blueprint provides a simple static rate limiting policy and a dashboard.
This policy uses the [`RateLimiter`](/reference/policies/spec.md#rate-limiter)
component.

## Configuration

<!-- Configuration Marker -->

```mdx-code-block

export const ParameterHeading = ({children}) => (
  <span style={{fontWeight: "bold"}}>{children}</span>
);

export const WrappedDescription = ({children}) => (
  <span style={{wordWrap: "normal"}}>{children}</span>
);

export const RefType = ({type, reference}) => (
  <a href={reference}>{type}</a>
);

export const ParameterDescription = ({name, type, reference, value, description}) => (
  <table class="blueprints-params">
  <tr>
    <td><ParameterHeading>Parameter</ParameterHeading></td>
    <td><code>{name}</code></td>
  </tr>
  <tr>
    <td><ParameterHeading>Type</ParameterHeading></td>
    <td><em>{reference == "" ? type : <RefType type={type} reference={reference} />}</em></td>
  </tr>
  <tr>
    <td class="blueprints-default-heading"><ParameterHeading>Default Value</ParameterHeading></td>
    <td><code>{value}</code></td>
  </tr>
  <tr>
    <td class="blueprints-description"><ParameterHeading>Description</ParameterHeading></td>
    <td class="blueprints-description"><WrappedDescription>{description}</WrappedDescription></td>
  </tr>
</table>
);
```

```mdx-code-block
import {apertureVersion as aver} from '../../../../apertureVersion.js'
```

Code: <a
href={`https://github.com/fluxninja/aperture/tree/${aver}/blueprints/policies/static-rate-limiting`}>policies/static-rate-limiting</a>

<h3 class="blueprints-h3">Common</h3>

<ParameterDescription
    name="common.policy_name"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Name of the policy.' />

<h3 class="blueprints-h3">Policy</h3>

<ParameterDescription
    name="policy.classifiers"
    type="[]aperture.spec.v1.Classifier"
    reference="../../spec#classifier"
    value="[]"
    description='List of classification rules.' />

<h4 class="blueprints-h4">Rate Limiter</h4>

<ParameterDescription
    name="policy.rate_limiter.rate_limit"
    type="float64"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Number of requests per `policy.rate_limiter.parameters.limit_reset_interval` to accept' />

<ParameterDescription
    name="policy.rate_limiter.flow_selector"
    type="aperture.spec.v1.FlowSelector"
    reference="../../spec#flow-selector"
    value="{'flow_matcher': {'control_point': '__REQUIRED_FIELD__'}, 'service_selector': {'service': '__REQUIRED_FIELD__'}}"
    description='A flow selector to match requests against' />

<ParameterDescription
    name="policy.rate_limiter.flow_selector.service_selector.service"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Service Name.' />

<ParameterDescription
    name="policy.rate_limiter.flow_selector.flow_matcher.control_point"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Control Point Name.' />

<ParameterDescription
    name="policy.rate_limiter.parameters"
    type="aperture.spec.v1.RateLimiterParameters"
    reference="../../spec#rate-limiter-parameters"
    value="{'label_key': '__REQUIRED_FIELD__', 'limit_reset_interval': '__REQUIRED_FIELD__'}"
    description='Parameters.' />

<ParameterDescription
    name="policy.rate_limiter.parameters.limit_reset_interval"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Time after which the limit for a given label value will be reset.' />

<ParameterDescription
    name="policy.rate_limiter.parameters.label_key"
    type="string"
    reference=""
    value="__REQUIRED_FIELD__"
    description='Flow label to use for rate limiting.' />

<ParameterDescription
    name="policy.rate_limiter.default_config"
    type="aperture.spec.v1.RateLimiterDynamicConfig"
    reference="../../spec#rate-limiter-dynamic-config"
    value="{'overrides': []}"
    description='Default configuration for rate limiter that can be updated at the runtime without shutting down the policy.' />

<ParameterDescription
    name="policy.rate_limiter.default_config.overrides"
    type="[]aperture.spec.v1.RateLimiterOverride"
    reference="../../spec#rate-limiter-override"
    value="[]"
    description='Allows to specify different limits for particular label values.' />

<h3 class="blueprints-h3">Dashboard</h3>

<ParameterDescription
    name="dashboard.refresh_interval"
    type="string"
    reference=""
    value="'10s'"
    description='Refresh interval for dashboard panels.' />

<h4 class="blueprints-h4">Datasource</h4>

<ParameterDescription
    name="dashboard.datasource.name"
    type="string"
    reference=""
    value="'$datasource'"
    description='Datasource name.' />

<ParameterDescription
    name="dashboard.datasource.filter_regex"
    type="string"
    reference=""
    value="''"
    description='Datasource filter regex.' />## Dynamic Configuration The
following configuration parameters can be
[dynamically configured](/reference/aperturectl/apply/dynamic-config/dynamic-config.md)
at runtime, without reloading the policy.

<h3 class="blueprints-h3">Dynamic Configuration</h3>

<ParameterDescription
    name="rate_limiter"
    type="aperture.spec.v1.RateLimiterDynamicConfig"
    reference="../../spec#rate-limiter-dynamic-config"
    value="__REQUIRED_FIELD__"
    description='Rate limiter dynamic configuration that is updated at runtime.' />
