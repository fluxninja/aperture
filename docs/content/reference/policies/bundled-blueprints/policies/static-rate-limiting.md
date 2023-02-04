---
title: Static Rate Limiting Policy
---

```mdx-code-block
import {apertureVersion} from '../../../../apertureVersion.js';
```

## Blueprint Location

GitHub: <a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/blueprints//policies/static-rate-limiting`}>static-rate-limiting</a>

## Introduction

This blueprint provides a simple static rate limiting policy and a dashboard.

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
    <td><code>{value != '' ? value : "REQUIRED VALUE"}</code></td>
  </tr>
  <tr>
    <td class="blueprints-description"><ParameterHeading>Description</ParameterHeading></td>
    <td class="blueprints-description"><WrappedDescription>{description}</WrappedDescription></td>
  </tr>
</table>
);
```

<h3 class="blueprints-h3">Common</h3>

<ParameterDescription
    name="common.policy_name"
    type="string"
    reference=""
    value=''
    description='Name of the policy.' />

<h3 class="blueprints-h3">Policy</h3>

<ParameterDescription
    name="policy.evaluation_interval"
    type="string"
    reference=""
    value=''
    description='How often should the policy be re-evaluated' />

<ParameterDescription
    name="policy.classifiers"
    type="[]aperture.spec.v1.Classifier"
    reference="../../spec#v1-classifier"
    value=''
    description='List of classification rules.' />

<h4 class="blueprints-h4">Rate Limiter</h4>

<ParameterDescription
    name="policy.rate_limiter.rate_limit"
    type="float64"
    reference=""
    value=''
    description='Number of requests per `policy.rate_limiter.parameters.limit_reset_interval` to accept' />

<ParameterDescription
    name="policy.rate_limiter.flow_selector"
    type="aperture.spec.v1.FlowSelector"
    reference="../../spec#v1-flow-selector"
    value=''
    description='A flow selector to match requests against' />

<ParameterDescription
    name="policy.rate_limiter.parameters"
    type="aperture.spec.v1.RateLimiterParameters"
    reference="../../spec#v1-rate-limiter-parameters"
    value=''
    description='Parameters.' />

<ParameterDescription
    name="policy.rate_limiter.parameters.label_key"
    type="string"
    reference=""
    value=''
    description='Flow label to use for rate limiting.' />

<ParameterDescription
    name="policy.rate_limiter.dynamic_config"
    type="aperture.spec.v1.RateLimiterDefaultConfig"
    reference="../../spec#v1-rate-limiter-default-config"
    value=''
    description='Dynamic configuration for rate limiter that can be applied at the runtime.' />

<h3 class="blueprints-h3">Dashboard</h3>

<ParameterDescription
    name="dashboard.refresh_interval"
    type="string"
    reference=""
    value=''
    description='Refresh interval for dashboard panels.' />

<h4 class="blueprints-h4">Datasource</h4>

<ParameterDescription
    name="dashboard.datasource.name"
    type="string"
    reference=""
    value=''
    description='Datasource name.' />

<ParameterDescription
    name="dashboard.datasource.filter_regex"
    type="string"
    reference=""
    value=''
    description='Datasource filter regex.' />
