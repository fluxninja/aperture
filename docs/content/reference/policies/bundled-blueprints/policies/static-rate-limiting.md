---
title: Static Rate Limiting Policy
---

```mdx-code-block
import {apertureVersion} from '../../../../apertureVersion.js';
```

## Blueprint Location

GitHub: <a
href={`https://github.com/fluxninja/aperture/tree/${apertureVersion}/blueprints/lib/1.0/policies/static-rate-limiting`}>static-rate-limiting</a>

## Introduction

This blueprint provides a simple static rate limiting policy and a dashboard.

## Configuration

<!-- Configuration Marker -->

### Common

**`common.policy_name`** (type: _`string`_)

required parameter

Name of the policy.

### Policy

**`policy.evaluation_interval`** (type: _`string`_)

default: `"300s"`

How often should the policy be re-evaluated

**`policy.classifiers`** (type: _`[]aperture.spec.v1.Classifier`_)

default: `[]`

List of classification rules.

#### Rate Limiter

**`policy.rate_limiter.rate_limit`** (type: _`float64`_)

required parameter

Number of requests per `policy.rate_limiter.parameters.limit_reset_interval` to
accept

**`policy.rate_limiter.flow_selector`** (type:
_`aperture.spec.v1.FlowSelector`_)

required parameter

A flow selector to match requests against

**`policy.rate_limiter.parameters`** (type:
_`aperture.spec.v1.RateLimiterParameters`_)

default:
`{'label_key': 'FAKE-VALUE', 'lazy_sync': {'enabled': True, 'num_sync': 5}, 'limit_reset_interval': '1s'}`

Parameters.

**`policy.rate_limiter.parameters.label_key`** (type: _`string`_)

required parameter

Flow label to use for rate limiting.

**`policy.rate_limiter.dynamic_config`** (type:
_`aperture.spec.v1.RateLimiterDefaultConfig`_)

default: `{'overrides': []}`

Dynamic configuration for rate limiter that can be applied at the runtime.

### Dashboard

**`dashboard.refresh_interval`** (type: _`string`_)

default: `"10s"`

Refresh interval for dashboard panels.

#### Datasource

**`dashboard.datasource.name`** (type: _`string`_)

default: `"$datasource"`

Datasource name.

**`dashboard.datasource.filter_regex`** (type: _`string`_)

default: `""`

Datasource filter regex.
